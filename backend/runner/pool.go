package runner

import (
	"context"
	"os"
	"strings"
	"sync"

	"github.com/juju/errors"

	"github.com/localhots/cmdui/backend/config"
	"github.com/localhots/cmdui/backend/db"
	"github.com/localhots/cmdui/backend/log"
)

var pool = &processPool{}

type processPool struct {
	lock  sync.RWMutex
	wg    sync.WaitGroup
	procs map[string]*Process
}

func FindProcess(id string) *Process {
	pool.lock.RLock()
	defer pool.lock.RUnlock()
	return pool.procs[id]
}

func Shutdown() {
	pool.close()
}

func PrepareLogsDir() error {
	return os.MkdirAll(config.Get().LogDir, 0755)
}

func (pp *processPool) add(ctx context.Context, p *Process) error {
	// Validate process
	if err := pp.validate(p); err != nil {
		return errors.Annotate(err, "Process validation failed")
	}

	// Register process
	pp.lock.Lock()
	if pp.procs == nil {
		pp.procs = make(map[string]*Process)
	}
	pp.procs[p.ID] = p
	pp.lock.Unlock()

	// Start process

	if err := p.exec.Start(); err != nil {
		log.WithFields(log.F{
			"id":    p.ID,
			"error": err,
		}).Error("Failed to start a command")
		return errors.Annotate(err, "Failed to start a process")
	}
	p.PID = p.exec.Process.Pid
	log.WithFields(log.F{
		"id":      p.ID,
		"pid":     p.PID,
		"command": strings.Join(p.exec.Args, " "),
	}).Info("Command started")
	tryUpdateState(p, db.JobStateStarted)

	pp.wg.Add(1)
	go pp.handleProcessExit(p)
	return nil
}

func (pp *processPool) handleProcessExit(p *Process) {
	defer pp.wg.Done()
	err := p.exec.Wait()

	pp.lock.Lock()
	delete(pp.procs, p.ID)
	pp.lock.Unlock()

	if err != nil {
		log.WithFields(log.F{
			"id":    p.ID,
			"pid":   p.PID,
			"error": err,
		}).Error("Command failed")
		tryUpdateState(p, db.JobStateFailed)
	} else {
		log.WithFields(log.F{
			"id":  p.ID,
			"pid": p.PID,
		}).Info("Command finished")
		tryUpdateState(p, db.JobStateFinished)
	}

	if err := p.close(); err != nil {
		log.WithFields(log.F{
			"id":    p.ID,
			"error": err,
		}).Error("Failed to close a job")
	}

	for _, onExit := range p.exitCallbacks {
		onExit(p)
	}
}

func (pp *processPool) procsList() []*Process {
	pp.lock.RLock()
	defer pp.lock.RUnlock()

	list := make([]*Process, len(pp.procs))
	i := 0
	for _, p := range pp.procs {
		list[i] = p
		i++
	}

	return list
}

func (pp *processPool) validate(p *Process) error {
	switch {
	case p == nil:
		return errors.New("Can't add an empty process")
	case p.ID == "":
		return errors.New("Can't add a process without an ID")
	case p.exec == nil:
		return errors.New("Process executable can't be empty")
	default:
		return nil
	}
}

func (pp *processPool) close() {
	pp.wg.Wait()
}

func tryUpdateState(p *Process, s db.JobState) {
	log.WithFields(log.F{
		"id":    p.ID,
		"state": s,
	}).Debug("Job state changed")
	if p.Job == nil {
		return
	}
	if err := p.Job.UpdateState(s); err != nil {
		log.WithFields(log.F{
			"id":    p.ID,
			"state": s,
			"error": err,
		}).Error("Job state changed")
	}
}

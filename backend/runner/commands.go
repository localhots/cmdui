package runner

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"time"

	"github.com/juju/errors"

	"github.com/localhots/cmdui/backend/commands"
	"github.com/localhots/cmdui/backend/config"
	"github.com/localhots/cmdui/backend/db"
	"github.com/localhots/cmdui/backend/log"
)

func Start(cmd commands.Command, user db.User) (*Process, error) {
	job := db.NewJob(cmd, user)
	if err := job.Create(); err != nil {
		return nil, errors.Annotate(err, "Failed to create a job")
	}

	p := &Process{
		ID:  job.ID,
		Job: &job,
	}

	basePath := config.Get().Commands.BasePath
	p.exec = exec.Command(basePath, cmd.CombinedArgs()...)
	fd, err := p.useLogfile(p.logfile())
	if err != nil {
		return nil, err
	}
	p.exec.Stdout = fd
	p.exec.Stderr = fd

	if err := pool.add(context.Background(), p); err != nil {
		return p, errors.Annotate(err, "Failed to start a process")
	}

	return p, nil
}

func ReadLogUpdates(ctx context.Context, p *Process, w io.Writer) (done <-chan struct{}, err error) {
	cmde := exec.CommandContext(ctx, "/usr/bin/tail", "-n", "100", "-f", p.logfile())
	return readLog(ctx, cmde, p, w)
}

func ReadFullLog(ctx context.Context, p *Process, w io.Writer) (done <-chan struct{}, err error) {
	cmde := exec.CommandContext(ctx, "/bin/cat", p.logfile())
	return readLog(ctx, cmde, p, w)
}

func readLog(ctx context.Context, cmde *exec.Cmd, p *Process, w io.Writer) (done <-chan struct{}, err error) {
	cmde.Stdout = w
	cmde.Stderr = w

	tp := &Process{
		ID:   sysID("log-access"),
		exec: cmde,
	}

	exited := make(chan struct{}, 2)
	p.onExit(func(p *Process) {
		exited <- struct{}{}
	})
	tp.onExit(func(p *Process) {
		exited <- struct{}{}
	})

	if err := pool.add(ctx, tp); err != nil {
		return nil, errors.Annotate(err, "Failed to start a tail process")
	}

	return exited, nil
}

func CommandsList() ([]commands.Command, error) {
	var buf bytes.Buffer
	cmde := exec.Command("/bin/bash", "-c", config.Get().Commands.ConfigCommand)
	cmde.Stdout = &buf
	cmde.Stderr = &buf

	p := &Process{
		ID:   sysID("commands-list"),
		exec: cmde,
	}

	exited := make(chan struct{})
	p.onExit(func(p *Process) {
		close(exited)
	})

	if err := pool.add(context.Background(), p); err != nil {
		return nil, errors.Annotate(err, "Failed to import commands")
	}
	<-exited

	body := buf.Bytes()
	var list []commands.Command
	if err := json.Unmarshal(body, &list); err != nil {
		log.WithFields(log.F{
			"error": err,
			"json":  string(body),
		}).Error("Invalid commands JSON")
		return nil, errors.Annotate(err, "Failed to decode commands JSON")
	}

	return list, nil
}

func sysID(name string) string {
	return fmt.Sprintf("%s-%d", name, time.Now().UnixNano())
}

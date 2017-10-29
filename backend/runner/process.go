package runner

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"

	"github.com/juju/errors"

	"github.com/localhots/cmdui/backend/commands"
	"github.com/localhots/cmdui/backend/config"
	"github.com/localhots/cmdui/backend/db"
)

type Process struct {
	ID      string           `json:"id"`
	PID     int              `json:"pid"`
	Job     *db.Job          `json:"job"`
	Command commands.Command `json:"command"`
	Out     io.Writer        `json:"-"`

	exec *exec.Cmd
	log  *os.File

	exitCallbacks []func(p *Process) `json:"-"`
}

func (p *Process) Signal(s syscall.Signal) error {
	return p.exec.Process.Signal(s)
}

func (p *Process) logfile() string {
	return fmt.Sprintf("%s/%s.log", config.Get().LogDir, p.ID)
}

func (p *Process) useLogfile(path string) (io.Writer, error) {
	fd, err := os.OpenFile(p.logfile(), os.O_CREATE|os.O_WRONLY, 0744)
	if err != nil {
		return nil, errors.Annotate(err, "Failed to create log file")
	}
	p.log = fd

	return fd, nil
}

func (p *Process) onExit(fn func(p *Process)) {
	p.exitCallbacks = append(p.exitCallbacks, fn)
}

func (p *Process) close() error {
	if p.log != nil {
		err := p.log.Close()
		if err != nil {
			return errors.Annotate(err, "Failed to close log file")
		}
	}
	return nil
}

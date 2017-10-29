package api

import (
	"context"
	"fmt"
	"net/http"
	"syscall"

	"github.com/juju/errors"

	"github.com/localhots/cmdui/backend/commands"
	"github.com/localhots/cmdui/backend/db"
	"github.com/localhots/cmdui/backend/runner"
)

func init() {
	router.GET("/api/jobs", protectedEndpoint(jobsIndexHandler))
	router.GET("/api/jobs/:job_id", protectedEndpoint(jobShowHandler))
	router.PUT("/api/jobs/:job_id", protectedEndpoint(jobActionHandler))
	router.GET("/api/jobs/:job_id/log", protectedEndpoint(jobLogHandler))
	router.GET("/api/commands/:cmd/jobs", protectedEndpoint(jobsIndexHandler))
	router.GET("/api/users/:user_id/jobs", protectedEndpoint(jobsIndexHandler))
}

// jobLogHandler returns job's log. If the command is still running, than the
// client would be attached to a log file and receive updates as well as few
// previous lines of it. If the command is completed, the entire log would be
// returned. If a `full` parameter is provided and the command is still running,
// the client would receive all existing log contents and future updates.
// GET /api/jobs/:job_id/log
// FIXME: What the fuck is this function
func jobLogHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	id := param(ctx, "job_id")
	proc := runner.FindProcess(id)
	if proc != nil {
		var done <-chan struct{}
		var err error
		if r.FormValue("full") != "" {
			done, err = runner.ReadFullLog(ctx, proc, unbufferedWriter{w})
		} else {
			done, err = runner.ReadLogUpdates(ctx, proc, unbufferedWriter{w})
		}
		if err != nil {
			renderError(w, err, http.StatusInternalServerError, "Failed to tail a log")
		}
		<-done
		return
	}

	proc = &runner.Process{ID: id}
	done, err := runner.ReadFullLog(ctx, proc, unbufferedWriter{w})
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	<-done
}

// jobsIndexHandler returns a list of jobs for a given criteria.
// GET /api/jobs
// GET /api/commands/:cmd/jobs
// GET /api/users/:user_id/jobs
func jobsIndexHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	renderJobs := func(jobs []db.Job, err error) {
		if err != nil {
			renderError(w, err, http.StatusInternalServerError, "Failed to find jobs")
		} else {
			renderJSON(w, jobs)
		}
	}

	switch {
	case param(ctx, "user_id") != "":
		id := param(ctx, "user_id")
		u, err := db.FindUser(id)
		if err != nil {
			renderError(w, err, http.StatusInternalServerError, "Failed to find a user")
		}
		if u == nil {
			err := fmt.Errorf("User not found: %s", id)
			renderError(w, err, http.StatusNotFound, "User not found")
			return
		}

		renderJobs(db.FindUserJobs(id, requestedPage(r)))
	case param(ctx, "cmd") != "":
		cmdName := param(ctx, "cmd")
		if _, ok := commands.Map()[cmdName]; ok {
			renderJobs(db.FindCommandJobs(cmdName, requestedPage(r)))
		} else {
			err := fmt.Errorf("Command not found: %s", cmdName)
			renderError(w, err, http.StatusNotFound, "Command not found")
		}
	default:
		renderJobs(db.FindAllJobs(requestedPage(r)))
	}
}

// jobShowHandler returns job details.
// GET /api/jobs/:job_id
func jobShowHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	id := param(ctx, "job_id")
	job, err := db.FindJob(id)
	if err != nil {
		renderError(w, err, http.StatusInternalServerError, "Failed to find job")
		return
	}
	if job != nil {
		renderJSON(w, job)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

// jobActionHandler performs certain actions on a job.
// PUT /api/jobs/:job_id
func jobActionHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		renderError(w, err, http.StatusBadRequest, "Failed to parse form")
		return
	}

	id := r.PostForm.Get("id")
	if id == "" {
		err := errors.New("Job ID is required")
		renderError(w, err, http.StatusBadRequest, err.Error())
		return
	}

	proc := runner.FindProcess(id)
	if proc == nil {
		err := fmt.Errorf("Job %q was not found", id)
		renderError(w, err, http.StatusNotFound, err.Error())
		return
	}

	sigName := r.PostForm.Get("signal")
	if sigName != "" {
		sig, err := signalFromName(sigName)
		if err != nil {
			renderError(w, err, http.StatusBadRequest, err.Error())
			return
		}

		err = proc.Signal(sig)
		if err != nil {
			renderError(w, err, http.StatusInternalServerError, "Failed to send signal to a process")
			return
		}
	}
}

func signalFromName(name string) (syscall.Signal, error) {
	switch name {
	case "SIGHUP":
		return syscall.SIGHUP, nil
	case "SIGINT":
		return syscall.SIGINT, nil
	case "SIGKILL":
		return syscall.SIGKILL, nil
	case "SIGQUIT":
		return syscall.SIGQUIT, nil
	case "SIGTERM":
		return syscall.SIGTERM, nil
	case "SIGTTIN":
		return syscall.SIGTTIN, nil
	case "SIGTTOU":
		return syscall.SIGTTOU, nil
	case "SIGUSR1":
		return syscall.SIGUSR1, nil
	case "SIGUSR2":
		return syscall.SIGUSR2, nil
	default:
		return 0, fmt.Errorf("Signal not supported: %s", name)
	}
}

package api

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/localhots/cmdui/backend/api/auth"
	"github.com/localhots/cmdui/backend/commands"
	"github.com/localhots/cmdui/backend/runner"
)

func init() {
	router.GET("/api/commands", protectedEndpoint(commandsHandler))
	router.POST("/api/exec", protectedEndpoint(execHandler))
}

// commandsHandler returns a list of commands that are registered on the server.
// GET /commands
func commandsHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	renderJSON(w, commands.List())
}

// execHandler accepts a form from the UI, decodes it into a command and
// attempts to execute it. Raw job ID is returned in the response.
// POST /exec
func execHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	cmd, err := commandFromRequest(r)
	if err != nil {
		renderError(w, err, http.StatusBadRequest, "Failed to build a command from request form")
		return
	}

	proc, err := runner.Start(cmd, auth.User(ctx))
	if err != nil {
		renderError(w, err, http.StatusInternalServerError, "Failed to execute a command")
		return
	}

	w.WriteHeader(http.StatusCreated)
	renderJSON(w, proc.Job)
}

// commandFromRequest parses a form submitted from UI and converts it into a
// command.
func commandFromRequest(r *http.Request) (commands.Command, error) {
	var cmd commands.Command
	err := r.ParseForm()
	if err != nil {
		return cmd, err
	}

	for _, c := range commands.List() {
		if c.Name == r.PostForm.Get("command") {
			cmd = c
			break
		}
	}
	if cmd.Name == "" {
		return cmd, fmt.Errorf("Unknown command: %q", r.PostForm.Get("command"))
	}

	return buildCommand(cmd, r.PostForm), nil
}

func buildCommand(cmd commands.Command, p url.Values) commands.Command {
	cmd.Args = p.Get("args")
	for i, f := range cmd.Flags {
		cmd.Flags[i].Value = p.Get(fmt.Sprintf("flags[%s]", f.Name))
	}
	return cmd
}

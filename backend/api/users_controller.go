package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/localhots/cmdui/backend/db"
)

func init() {
	router.GET("/api/users/:user_id", protectedEndpoint(userDetailsHandler))
}

// userDetailsHandler returns user details.
// GET /api/users/:user_id
func userDetailsHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
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

	renderJSON(w, u)
}

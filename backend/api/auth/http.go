package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/juju/errors"
)

const (
	sessionCookieName = "cmdui_session_id"
)

func AuthenticateRequest(w http.ResponseWriter, r *http.Request) (context.Context, error) {
	cook, err := r.Cookie(sessionCookieName)
	if err != nil {
		return r.Context(), errors.Annotate(err, "Failed to get cookie value")
	}
	sess, err := FindSession(cook.Value)
	if err != nil {
		return r.Context(), errors.Annotate(err, "Failed to authenticate request")
	}
	ctx := ContextWithSession(r.Context(), sess)
	if sess.ExpiresAt.Before(time.Now()) {
		ClearCookie(ctx, w)
		return ctx, errors.New("Session expired")
	}

	u, err := sess.User()
	if err != nil {
		return ctx, errors.Annotatef(err, "Failed to find user %d", sess.UserID)
	}
	if u == nil {
		return ctx, errors.UserNotFoundf("User %s was not found", sess.UserID)
	}
	u.Authorized = true

	ctx = ContextWithUser(ctx, *u)
	return ctx, nil
}

func AuthorizeResponse(ctx context.Context, w http.ResponseWriter) {
	if sess, ok := SessionFromContext(ctx); ok {
		http.SetCookie(w, &http.Cookie{
			Name:     sessionCookieName,
			Value:    sess.ID,
			Path:     "/",
			Expires:  sess.ExpiresAt,
			HttpOnly: true,
		})
	}
}

func ClearCookie(ctx context.Context, w http.ResponseWriter) {
	sess := Session(ctx)
	sess.ExpiresAt = time.Time{}
	ctx = ContextWithSession(ctx, sess)
	AuthorizeResponse(ctx, w)
}

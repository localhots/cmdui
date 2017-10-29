package auth

import (
	"context"

	"github.com/localhots/cmdui/backend/db"
)

type ctxKey string

const (
	ctxSessionKey ctxKey = "session"
	ctxUserKey    ctxKey = "user"
)

func ContextWithSession(ctx context.Context, sess db.Session) context.Context {
	return context.WithValue(ctx, ctxSessionKey, sess)
}

func SessionFromContext(ctx context.Context) (sess db.Session, ok bool) {
	v := ctx.Value(ctxSessionKey)
	if v == nil {
		return db.Session{}, false
	}

	return v.(db.Session), true
}

func Session(ctx context.Context) db.Session {
	sess, _ := SessionFromContext(ctx)
	return sess
}

func ContextWithUser(ctx context.Context, u db.User) context.Context {
	return context.WithValue(ctx, ctxUserKey, u)
}

func UserFromContext(ctx context.Context) (u db.User, ok bool) {
	v := ctx.Value(ctxUserKey)
	if v == nil {
		return db.User{}, false
	}

	return v.(db.User), true
}

func User(ctx context.Context) db.User {
	u, _ := UserFromContext(ctx)
	return u
}

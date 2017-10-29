package auth

import (
	"sync"

	"github.com/juju/errors"

	"github.com/localhots/cmdui/backend/db"
)

var (
	sessionCacheMux    sync.Mutex
	sessionCache       = map[string]db.Session{}
	errSessionNotFound = errors.New("Session not found")
)

func FindSession(id string) (db.Session, error) {
	if id == "" {
		return db.Session{}, errSessionNotFound
	}

	sessionCacheMux.Lock()
	sessc, ok := sessionCache[id]
	sessionCacheMux.Unlock()
	if ok {
		return sessc, nil
	}

	sess, err := db.FindSession(id)
	if err != nil {
		return db.Session{}, errors.Annotate(err, "Session lookup failed")
	}
	if sess == nil {
		return db.Session{}, errSessionNotFound
	}

	sessionCacheMux.Lock()
	sessionCache[sess.ID] = *sess
	sessionCacheMux.Unlock()

	return *sess, nil
}

func CacheSession(sess db.Session) {
	if sess.ID == "" {
		return
	}

	sessionCacheMux.Lock()
	sessionCache[sess.ID] = sess
	sessionCacheMux.Unlock()
}

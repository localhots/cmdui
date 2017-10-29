package db

import (
	"database/sql"
	"time"

	"github.com/juju/errors"
)

type Session struct {
	ID        string    `db:"id"`
	UserID    string    `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
	ExpiresAt time.Time `db:"expires_at"`
}

func NewSession(userID string) Session {
	const ttl = 6 * 30 * 24 * time.Hour // 6 months
	now := time.Now().UTC()
	exp := now.Add(ttl)

	return Session{
		ID:        newID(),
		UserID:    userID,
		CreatedAt: now,
		ExpiresAt: exp,
	}
}

func FindSession(id string) (*Session, error) {
	var s Session
	err := db.Get(&s, "SELECT * FROM sessions WHERE id = ?", id)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, errors.Annotate(err, "Failed to load session details")
		}
		return nil, nil
	}

	return &s, nil
}

func (s Session) Create() error {
	_, err := db.NamedExec(`
        INSERT INTO sessions
        SET
            id = :id,
            user_id = :user_id,
            created_at = :created_at,
            expires_at = :expires_at
    `, s)
	if err != nil {
		return errors.Annotate(err, "Failed to create a session")
	}
	return nil
}

func (s Session) User() (*User, error) {
	return FindUser(s.UserID)
}

package db

import (
	"database/sql"
	"fmt"

	"github.com/juju/errors"
)

type User struct {
	ID          string `json:"id" db:"id"`
	Name        string `json:"name" db:"github_name"`
	Picture     string `json:"picture" db:"github_picture"`
	GithubID    uint   `json:"-" db:"github_id"`
	GithubLogin string `json:"-" db:"github_login"`

	Authorized bool `json:"authorized" db:"-"`
}

func NewUser() User {
	return User{ID: newID()}
}

func FindAllUsers() (map[string]User, error) {
	return findUsers("SELECT * FROM users ORDER BY id ASC")
}

func FindUsers(ids ...string) (map[string]User, error) {
	if len(ids) == 0 {
		return map[string]User{}, nil
	}

	query := fmt.Sprintf("SELECT * FROM users WHERE id IN (%s)", placeholders(ids))
	return findUsers(query, iargs(ids)...)
}

func findUsers(query string, args ...interface{}) (map[string]User, error) {
	var users []User
	err := db.Select(&users, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return nil, errors.Annotate(err, "Failed to load users list")
	}

	return usersSliceToMap(users), nil
}

func FindUser(id string) (*User, error) {
	return findUser("SELECT * FROM users WHERE id = ?", id)
}

func FindUserByGithubID(id uint) (*User, error) {
	return findUser("SELECT * FROM users WHERE github_id = ?", id)
}

func FindUserByLogin(login string) (*User, error) {
	return findUser("SELECT * FROM users WHERE github_login = ?", login)
}

func findUser(query string, args ...interface{}) (*User, error) {
	var u User
	err := db.Get(&u, query, args...)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, errors.Annotate(err, "Failed to load user details")
		}
		return nil, nil
	}

	return &u, nil
}

func (u User) Create() error {
	_, err := db.NamedExec(`
        INSERT INTO users (id, github_id, github_login, github_name, github_picture)
        VALUES (:id, :github_id, :github_login, :github_name, :github_picture)
    `, u)
	if err != nil {
		return errors.Annotate(err, "Failed to create a user")
	}
	return nil
}

func (u User) Update() error {
	_, err := db.NamedExec(`
        UPDATE users
        SET
            github_login = :github_login,
            github_name = :github_name,
            github_picture = :github_picture
        WHERE
            github_id = :github_id
    `, u)
	if err != nil {
		return errors.Annotate(err, "Failed to update a user")
	}
	return nil
}

func usersSliceToMap(s []User) map[string]User {
	m := make(map[string]User, len(s))
	for _, u := range s {
		m[u.ID] = u
	}
	return m
}

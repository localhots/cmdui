package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/juju/errors"

	"github.com/localhots/cmdui/backend/commands"
)

type JobState string

const (
	JobStateNew      JobState = "new"
	JobStateCreated  JobState = "created"
	JobStateStarted  JobState = "started"
	JobStateAborted  JobState = "aborted"
	JobStateFailed   JobState = "failed"
	JobStateFinished JobState = "finished"
)

type Job struct {
	ID         string     `json:"id" db:"id"`
	Command    string     `json:"command" db:"command"`
	Args       string     `json:"args" db:"args"`
	Flags      string     `json:"flags" db:"flags"`
	UserID     string     `json:"user_id" db:"user_id"`
	User       *User      `json:"user" db:"-"`
	State      string     `json:"state" db:"state"`
	CreatedAt  *time.Time `json:"created_at" db:"created_at"`
	StartedAt  *time.Time `json:"started_at" db:"started_at"`
	FinishedAt *time.Time `json:"finished_at" db:"finished_at"`
}

func NewJob(c commands.Command, u User) Job {
	return Job{
		ID:      newID(),
		Command: c.Name,
		Args:    c.Args,
		Flags:   c.FlagsString(),
		UserID:  u.ID,
		User:    &u,
		State:   string(JobStateNew),
	}
}

func FindAllJobs(p Page) ([]Job, error) {
	return findJobsWhere(jobsIndexQuery("", p))
}

func FindCommandJobs(name string, p Page) ([]Job, error) {
	return findJobsWhere(jobsIndexQuery("WHERE command = ?", p), name)
}

func FindUserJobs(id string, p Page) ([]Job, error) {
	return findJobsWhere(jobsIndexQuery("WHERE user_id = ?", p), id)
}

func jobsIndexQuery(where string, p Page) string {
	p = p.normalize()
	return fmt.Sprintf("SELECT * FROM jobs %s ORDER BY created_at DESC LIMIT %d, %d",
		where, p.Offset, p.Limit)
}

func findJobsWhere(query string, args ...interface{}) ([]Job, error) {
	var jobs []Job
	err := db.Select(&jobs, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return nil, errors.Annotate(err, "Failed to load Jobs list")
	}

	userIDs := stringSet{}
	for _, r := range jobs {
		userIDs.add(r.UserID)
	}
	users, err := FindUsers(userIDs.items()...)
	if err != nil {
		return nil, errors.Annotate(err, "Failed to find users to embed into jobs")
	}
	for i, r := range jobs {
		if u, ok := users[r.UserID]; ok {
			jobs[i].User = &u
		}
	}

	return jobs, nil
}

func FindJob(id string) (*Job, error) {
	var r Job
	err := db.Get(&r, "SELECT * FROM jobs WHERE id = ?", id)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, errors.Annotate(err, "Failed to load Job details")
		}
		return nil, nil
	}

	if r.UserID != "" {
		r.User, err = FindUser(r.UserID)
		if err != nil {
			return nil, errors.Annotate(err, "Failed to find a user to embed into a job")
		}
	}

	return &r, nil
}

func (r *Job) UpdateState(s JobState) error {
	r.State = string(s)
	ts := time.Now().UTC()
	switch s {
	case JobStateStarted:
		r.StartedAt = &ts
	case JobStateFinished, JobStateAborted, JobStateFailed:
		r.FinishedAt = &ts
	}

	return r.Update()
}

func (r *Job) Create() error {
	ts := time.Now().UTC()
	r.CreatedAt = &ts
	r.State = string(JobStateCreated)

	_, err := db.NamedExec(`
        INSERT INTO jobs
        SET
            id = :id,
            command = :command,
            args = :args,
            flags = :flags,
            user_id = :user_id,
            state = :state,
            created_at = :created_at
    `, r)
	if err != nil {
		return errors.Annotate(err, "Failed to create a job")
	}
	return nil
}

func (r Job) Update() error {
	_, err := db.NamedExec(`
        UPDATE jobs
        SET
            state = :state,
            started_at = :started_at,
            finished_at = :finished_at
        WHERE
            id = :id
    `, r)
	if err != nil {
		return errors.Annotate(err, "Failed to update a job")
	}
	return nil
}

func jobsSliceToMap(s []Job) map[string]Job {
	m := make(map[string]Job, len(s))
	for _, r := range s {
		m[r.ID] = r
	}
	return m
}

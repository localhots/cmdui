package db

import (
	"reflect"
	"strings"

	_ "github.com/go-sql-driver/mysql" // Import MySQL driver
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // Import SQLite driver
	uuid "github.com/satori/go.uuid"

	"github.com/localhots/cmdui/backend/config"
)

var (
	db *sqlx.DB
)

func Connect() error {
	var err error
	cfg := config.Get().Database
	db, err = sqlx.Connect(cfg.Driver, cfg.Spec)
	return err
}

type Page struct {
	Offset uint
	Limit  uint
}

func (p Page) normalize() Page {
	const defaultPerPage = 50
	if p.Limit == 0 {
		p.Limit = defaultPerPage
	}
	return p
}

//
// Helpers
//

func newID() string {
	return uuid.NewV4().String()
}

func placeholders(val interface{}) string {
	v := reflect.ValueOf(val)
	if v.Kind() == reflect.Slice {
		s := strings.Repeat("?, ", v.Len())
		return s[0 : len(s)-2]
	}

	return "?"
}

func iargs(args []string) []interface{} {
	iargs := make([]interface{}, len(args))
	for i, arg := range args {
		iargs[i] = arg
	}
	return iargs
}

type stringSet map[string]struct{}

func (s stringSet) add(items ...string) {
	for _, item := range items {
		s[item] = struct{}{}
	}
}

func (s stringSet) items() []string {
	l := make([]string, len(s))
	i := 0
	for item := range s {
		l[i] = item
		i++
	}
	return l
}

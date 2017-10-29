package log

import (
	"github.com/Sirupsen/logrus"
)

// F is short for "fields".
type F map[string]interface{}

// WithFields is an entry function for logging.
func WithFields(f F) *logrus.Entry {
	return logrus.WithFields(logrus.Fields(f))
}

// Logger returns a logger instance.
func Logger() *logrus.Logger {
	return logrus.StandardLogger()
}

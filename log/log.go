// Package log provides customized Logrus loggers.
package log

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/ahiho/gocandy/env"
)

// Create a logger configured via environment variables.
//
// The following is a list of environment variables used:
//	APP_VERSION (string) sets app.version which is part of all log messages.
//	HOST (string) sets app.host which is part of all log messages.
//	LOG_JSON (bool) controls whether to output logs in JSON format with timestamps in time.RFC3339Nano format. Defaults to false.
//	LOG_CALLERS (bool) controls whether to include the calling method as a field in logs. Defaults to false. https://godoc.org/github.com/sirupsen/logrus#SetReportCaller
//	LOG_LEVEL (string) sets the log level. Defaults to "trace". https://godoc.org/github.com/sirupsen/logrus#Level
//
// If LOG_JSON is enabled, timestamps are switched
func New() *logrus.Entry {
	l := logrus.New()
	if env.GetBool("LOG_JSON") {
		l.Formatter = &logrus.JSONFormatter{
			DataKey:         "data",
			TimestampFormat: time.RFC3339Nano,
		}
	}

	l.SetReportCaller(env.GetBool("LOG_CALLERS"))

	if lvl, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL")); err == nil {
		l.Level = lvl
	} else {
		l.Level = logrus.TraceLevel
	}

	log := l.WithFields(logrus.Fields{
		"app": map[string]string{
			"host":    os.Getenv("HOST"),
			"version": os.Getenv("APP_VERSION"),
		},
	})

	return log
}

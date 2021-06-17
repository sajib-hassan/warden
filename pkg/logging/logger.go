// Package logging provides structured logging with logrus.
package logging

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	// Logger is a configured logrus.Logger.
	Logger *logrus.Logger
)

// NewLogger creates and configures a new logrus Logger.
func NewLogger() *logrus.Logger {
	Logger = logrus.New()
	if viper.GetBool("log_textlogging") {
		Logger.Formatter = &logrus.TextFormatter{
			DisableTimestamp: true,
		}
	} else {
		Logger.Formatter = &logrus.JSONFormatter{
			DisableTimestamp: true,
		}
	}

	level := viper.GetString("log_level")
	if level == "" {
		level = "error"
	}
	l, err := logrus.ParseLevel(level)
	if err != nil {
		log.Fatal(err)
	}
	Logger.Level = l
	return Logger
}

// NewStructuredLogger implements a custom structured logrus Logger.
func NewStructuredLogger(logger *logrus.Logger) func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: logger})
}

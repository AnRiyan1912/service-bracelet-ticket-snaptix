package xlogger

import (
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"

	"bracelet-ticket-system-be/internal/config"
)

var (
	Logger *zerolog.Logger
)

func Setup(cfg config.Config) {
	if cfg.IsDevelopment {

		l := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).
			With().
			Timestamp().
			Logger()
		l.Level(zerolog.DebugLevel)
		Logger = &l
		return
	}

	logDir := "logs"
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.Mkdir(logDir, os.ModePerm)
		if err != nil {
			panic("Failed to create logs directory: " + err.Error())
		}
	}

	logFile := filepath.Join(logDir, "app.log")
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("Failed to open log file: " + err.Error())
	}

	l := zerolog.New(file).
		With().
		Timestamp().
		Logger()
	Logger = &l
}

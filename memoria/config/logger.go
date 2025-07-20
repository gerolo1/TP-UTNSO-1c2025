package config

import (
	"io"
	"log/slog"
	"os"
)

func InitLogger(logName string, logLevel string) *slog.Logger {
	logger := &slog.Logger{}

	logFile, e := os.OpenFile(logName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if e != nil {
		panic(e)
	}

	mw := io.MultiWriter(os.Stdout, logFile)

	logger = slog.New(slog.NewTextHandler(mw, &slog.HandlerOptions{
		Level: GetLogLevel(logLevel),
	}))

	return logger
}

func GetLogLevel(logLevel string) slog.Level {
	switch logLevel {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "ERROR":
		return slog.LevelError
	case "WARN":
		return slog.LevelWarn
	}
	return slog.LevelInfo
}

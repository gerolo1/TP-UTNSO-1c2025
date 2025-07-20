package config

import (
	"io"
	"log/slog"
	"os"
)

var Logger *slog.Logger

func InitLogger(logName string) {
	logFile, e := os.OpenFile(logName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if e != nil {
		panic(e)
	}

	mw := io.MultiWriter(os.Stdout, logFile)

	Logger = slog.New(slog.NewTextHandler(mw, &slog.HandlerOptions{
		Level: GetLogLevel(),
	}))
}

func GetLogLevel() slog.Level {
	switch Config.LogLevel {
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

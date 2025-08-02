package config

import (
	"context"
	"log"
	"log/slog"
	"mimic/lib/utils"
)

type AppConfig struct {
	GoMimicPort  uint16
	AdminPort    uint16
	AdminToken   string
	LogFilter    slog.Level
	MongodbUrl   string
	DatabaseName string
	Ctx          context.Context
}

func DefaultLogLevel() slog.Level {
	logLevelStr := utils.EnvOrDefault("LOG_LEVEL", "info")

	logLevel := slog.LevelInfo
	switch logLevelStr {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		log.Printf("invalid value set for LOG_LEVEL, default to 'info'")
	}

	return logLevel
}

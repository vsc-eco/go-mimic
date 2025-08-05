package config

import (
	"log"
	"log/slog"
	"mimic/lib/utils"
)

type AppConfig struct {
	GoMimic GoMimicConfig
	Admin   AdminConfig

	LogFilter    slog.Level
	MongodbUrl   string
	DatabaseName string
}

type GoMimicConfig struct {
	Port uint16
}

type AdminConfig struct {
	Port  uint16
	Token string
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

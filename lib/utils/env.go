package utils

import (
	"fmt"
	"log/slog"
	"os"
)

func EnvOrPanic(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		panic(fmt.Sprintf("environment variable not set: `%s`", key))
	}
	return value
}

func EnvOrDefault(key, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		slog.Warn("Environment variable not set, using default.", "key", key, "default", defaultValue)
		return defaultValue
	}
	return value
}

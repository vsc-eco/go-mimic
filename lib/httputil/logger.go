package httputil

import (
	"log/slog"

	"github.com/go-chi/httplog/v3"
)

var defaultLoggerOptions = httplog.Options{
	// Level defines the verbosity of the request logs:
	// slog.LevelDebug - log all responses (incl. OPTIONS)
	// slog.LevelInfo  - log responses (excl. OPTIONS)
	// slog.LevelWarn  - log 4xx and 5xx responses only (except for 429)
	// slog.LevelError - log 5xx responses only
	Level: slog.LevelInfo,

	// Set log output to Elastic Common Schema (ECS) format.
	Schema: httplog.SchemaECS,

	// RecoverPanics recovers from panics occurring in the underlying HTTP handlers
	// and middlewares. It returns HTTP 500 unless response status was already set.
	//
	// NOTE: Panics are logged as errors automatically, regardless of this setting.
	RecoverPanics: true,
}

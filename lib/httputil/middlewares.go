package httputil

import (
	"crypto/subtle"
	"encoding/hex"
	"log/slog"
	"net/http"

	"github.com/go-chi/httplog/v3"
)

func RequestTrace(logger *slog.Logger) func(http.Handler) http.Handler {
	return httplog.RequestLogger(
		logger,
		&defaultLoggerOptions,
	)
}

func AuthMiddleware(
	adminToken []byte,
	logger *slog.Logger,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr := r.Header.Get("X-ADMIN-TOKEN")
			if len(tokenStr) == 0 {
				w.WriteHeader(http.StatusForbidden)
				return

			}

			token, err := hex.DecodeString(tokenStr)
			if err != nil || len(token) != len(adminToken) {
				logger.Warn("invalid token.", "token", tokenStr)
				w.WriteHeader(http.StatusForbidden)
				return
			}

			tokenMatched := subtle.ConstantTimeCompare(adminToken, token) == 1
			if !tokenMatched {
				logger.Warn(
					"failed to authenticate request.",
					"token",
					tokenStr,
				)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

package httputil

import (
	"crypto/subtle"
	"encoding/hex"
	"log/slog"
	"net/http"
)

// verifies that the header `X-ADMIN-TOKEN` is present and matches the exported admin token.
// token verification is done with `subtle.ConstantTimeCompare` to prevent timing attacks.
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

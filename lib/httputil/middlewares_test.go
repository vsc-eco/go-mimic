package httputil

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	var validToken [64]byte
	if _, err := io.ReadFull(rand.Reader, validToken[:]); err != nil {
		t.Fatal(err)
	}

	validTokenHex := hex.EncodeToString(validToken[:])
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	// Mock handler to verify successful authentication
	mockHandler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusTeapot)
		},
	)

	tests := []struct {
		name           string
		token          string
		expectedStatus int
	}{
		{
			name:           "valid token",
			token:          validTokenHex,
			expectedStatus: http.StatusTeapot,
		},
		{
			name:           "missing token",
			token:          "",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "invalid hex token",
			token:          "invalid-hex",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "wrong length token",
			token:          hex.EncodeToString([]byte("short")),
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "wrong token value",
			token:          hex.EncodeToString(make([]byte, len(validToken))),
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			middleware := AuthMiddleware(validToken[:], logger)
			handler := middleware(mockHandler)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.token != "" {
				req.Header.Set("X-ADMIN-TOKEN", tt.token)
			}

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

package admin

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandlerUserCreate(t *testing.T) {
	var token [64]byte

	_, err := io.ReadFull(rand.Reader, token[:])
	if err != nil {
		t.Fatal(err)
	}

	srv := serverHandler{slog.Default()}

	// test requests
	requestJson, err := json.Marshal(map[string]string{
		"account":  "test@example.com",
		"password": "password",
	})
	if err != nil {
		t.Fatal(err)
	}

	requestBody := io.NopCloser(bytes.NewReader(requestJson))
	req := &http.Request{
		Method: http.MethodPost,
		Header: http.Header{
			"Content-Type": []string{"application/json"},
			// "X-ADMIN-TOKEN": []string{tokenStr},
		},
		Body: requestBody,
	}

	// Create a ResponseWriter for testing
	w := httptest.NewRecorder()
	srv.newUser(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

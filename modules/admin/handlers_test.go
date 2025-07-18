package admin

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"mimic/modules/db/mimic/accountdb"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	mongoDuplicateError = mongo.WriteError{
		Index:   0,
		Code:    11000, // MongoDB duplicate key error code
		Message: "E11000 duplicate key error collection: test.accounts index: email_1 dup key: { email: \"test@example.com\" }",
	}

	errServerError = errors.New("stub server error")
)

// mockAccountDB implements accountdb.AccountQuery
type mockAccountDB struct {
	insertAccountErr error
}

func (mockaccountdb *mockAccountDB) InsertAccount(
	_ context.Context,
	_ *accountdb.Account,
) error {
	return mockaccountdb.insertAccountErr
}

func (mockaccountdb *mockAccountDB) QueryAccountByNames(
	_ context.Context,
	_ *[]accountdb.Account,
	_ []string,
) error {
	panic("not implemented")
}

func TestHandlerUserCreate(t *testing.T) {
	var token [64]byte

	_, err := io.ReadFull(rand.Reader, token[:])
	if err != nil {
		t.Fatal(err)
	}

	mockDb := &mockAccountDB{}

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	srv := serverHandler{logger, mockDb}

	testCases := []struct {
		insertError    error
		account        string
		password       string
		expectedStatus int
	}{
		{nil, "", "", http.StatusBadRequest},
		{nil, "foo", "", http.StatusBadRequest},
		{nil, "", "bar", http.StatusBadRequest},
		{nil, "foo", "bar", http.StatusCreated},
		{errServerError, "foo", "bar", http.StatusInternalServerError},
		{errServerError, "", "", http.StatusBadRequest},
		{errServerError, "foo", "", http.StatusBadRequest},
		{errServerError, "", "bar", http.StatusBadRequest},
		{mongoDuplicateError, "foo", "bar", http.StatusConflict},
		{mongoDuplicateError, "", "", http.StatusBadRequest},
		{mongoDuplicateError, "foo", "", http.StatusBadRequest},
		{mongoDuplicateError, "", "bar", http.StatusBadRequest},
	}

	for _, tt := range testCases {
		w := httptest.NewRecorder()
		mockDb.insertAccountErr = tt.insertError

		requestJson, err := json.Marshal(map[string]string{
			"account":  tt.account,
			"password": tt.password,
		})

		if err != nil {
			t.Fatal(err)
		}

		req := &http.Request{
			Method: http.MethodPost,
			Header: http.Header{
				"Content-Type": []string{"application/json"},
				// "X-ADMIN-TOKEN": []string{tokenStr},
			},
			Body: io.NopCloser(bytes.NewReader(requestJson)),
		}

		srv.newUser(w, req)
		assert.Equal(t, tt.expectedStatus, w.Code)
	}
}

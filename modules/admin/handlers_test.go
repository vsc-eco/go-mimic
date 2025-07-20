package admin

import (
	"bytes"
	"context"
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
	err error
}

func (mockaccountdb *mockAccountDB) UpdateAccountKeySet(
	_ context.Context,
	_ string,
	_ *accountdb.UserKeySet,
) error {
	return mockaccountdb.err
}

func (mockaccountdb *mockAccountDB) InsertAccount(
	_ context.Context,
	_ *accountdb.Account,
) error {
	return mockaccountdb.err
}

func (mockaccountdb *mockAccountDB) QueryAccountByNames(
	_ context.Context,
	_ *[]accountdb.Account,
	_ []string,
) error {
	panic("not implemented")
}

func TestHandlerUserCreate(t *testing.T) {
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
		mockDb.err = tt.insertError

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
			},
			Body: io.NopCloser(bytes.NewReader(requestJson)),
		}

		srv.newUser(w, req)
		assert.Equal(t, tt.expectedStatus, w.Code)
	}
}

func TestHandlerUserUpdateKey(t *testing.T) {
	type TestCase struct {
		dbErr          error
		account        string
		password       string
		expectedStatus int
	}

	testTable := map[string][]TestCase{
		"no database error": {
			{nil, "foo", "bar", http.StatusNoContent},
			{nil, "", "", http.StatusBadRequest},
			{nil, "foo", "", http.StatusBadRequest},
			{nil, "", "bar", http.StatusBadRequest},
		},
		"database internal error": {
			{errServerError, "foo", "bar", http.StatusInternalServerError},
			{errServerError, "", "", http.StatusBadRequest},
			{errServerError, "foo", "", http.StatusBadRequest},
			{errServerError, "", "bar", http.StatusBadRequest},
		},
		"database ErrDocumentNotFound": {
			{accountdb.ErrAccountNotFound, "foo", "bar", http.StatusNotFound},
			{accountdb.ErrAccountNotFound, "", "", http.StatusBadRequest},
			{accountdb.ErrAccountNotFound, "foo", "", http.StatusBadRequest},
			{accountdb.ErrAccountNotFound, "", "bar", http.StatusBadRequest},
		},
	}

	mockDb := &mockAccountDB{}

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	srv := serverHandler{logger, mockDb}

	for testName, testCases := range testTable {
		t.Run(testName, func(t *testing.T) {
			for _, tt := range testCases {
				w := httptest.NewRecorder()
				mockDb.err = tt.dbErr

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
					},
					Body: io.NopCloser(bytes.NewReader(requestJson)),
				}

				srv.updateUser(w, req)
				assert.Equal(t, tt.expectedStatus, w.Code)

			}
		})
	}
}

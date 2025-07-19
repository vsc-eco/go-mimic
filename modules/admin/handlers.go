package admin

import (
	"context"
	"encoding/json"
	"errors"
	"mimic/lib/hivekey"
	"mimic/modules/admin/services"
	"mimic/modules/db/mimic/accountdb"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type userCredentials struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

// Request:
//   - request json body: { account: string, password: string }.
//
// Response:
//   - returns 400 Bad Request if the request body is invalid
//   - returns 409 Conflict if the account already exists
//   - returns 201 Created if the user is created successfully.
func (h *serverHandler) newUser(w http.ResponseWriter, r *http.Request) {
	var credentials userCredentials

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		h.logger.Error("failed to decode request body.", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	account, err := services.MakeAccount(
		credentials.Account,
		credentials.Password,
	)
	if err != nil {
		h.logger.Error("failed to create user account.", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := h.db.InsertAccount(ctx, account); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			w.WriteHeader(http.StatusConflict)
		} else {
			h.logger.Error("failed to insert user account to database.", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *serverHandler) updateUser(w http.ResponseWriter, r *http.Request) {
	var credentials userCredentials

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		h.logger.Error("failed to decode request body.", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(credentials.Account) == 0 || len(credentials.Password) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newKeySet := hivekey.MakeHiveKeySet(
		credentials.Account,
		credentials.Password,
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := h.db.UpdateAccountKeySet(ctx, credentials.Account, &newKeySet)
	if err != nil {
		if errors.Is(err, accountdb.ErrDocumentNotFound) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			h.logger.Error("failed to update user key.", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

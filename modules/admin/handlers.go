package admin

import (
	"context"
	"encoding/json"
	"mimic/modules/admin/services"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// Request:
//   - request json body: { account: string, password: string }.
//
// Response:
//   - returns 400 Bad Request if the request body is invalid
//   - returns 409 Conflict if the account already exists
//   - returns 201 Created if the user is created successfully.
func (h *serverHandler) newUser(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Account  string `json:"account"`
		Password string `json:"password"`
	}

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

package admin

import (
	"context"
	"encoding/json"
	"errors"
	"mimic/lib/utils"
	"mimic/lib/validator"
	"mimic/modules/db/mimic/accountdb"
	"net/http"
	"time"

	"github.com/vsc-eco/hivego"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type userCreateParam struct {
	Account             string        `json:"account"               validate:"required"`
	Owner               *hivego.Auths `json:"owner"                 validate:"required"`
	Active              *hivego.Auths `json:"active"                validate:"required"`
	Posting             *hivego.Auths `json:"posting"               validate:"required"`
	JsonMetadata        string        `json:"json_metadata"         validate:"json,omitempty"`
	PostingJsonMetadata string        `json:"posting_json_metadata" validate:"json,omitempty"`
}

// Request: POST http://0.0.0.0:3001/user
//   - request json body: userCreateParam
//
// Response:
//   - returns 201 Created if the user is created successfully.
//   - returns 400 Bad Request otherwise
func (h *serverHandler) newUser(w http.ResponseWriter, r *http.Request) {
	var c userCreateParam

	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		h.logger.Error("failed to decode request json.", "err", err)
		http.Error(w, "failed to decode request", http.StatusBadRequest)
		return
	}

	if err := validator.New().Struct(&c); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	ts := time.Now().Format(utils.TimeFormat)

	account := &accountdb.Account{
		ObjectId: primitive.NilObjectID,
		Name:     c.Account,
		UserKeySet: accountdb.UserKeySet{
			Owner:   c.Owner,
			Active:  c.Active,
			Posting: c.Posting,
		},
		LastOwnerUpdate:   ts,
		LastAccountUpdate: ts,
		Created:           ts,
		JsonMeta:          c.JsonMetadata,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := h.db.InsertAccount(ctx, account); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			http.Error(w, "user exists", http.StatusBadRequest)
		} else {
			h.logger.Error("failed to insert user account to database.", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
}

type userUpdateParam struct {
	Account             string        `json:"account"               validate:"required"`
	Owner               *hivego.Auths `json:"owner"`
	Active              *hivego.Auths `json:"active"`
	Posting             *hivego.Auths `json:"posting"`
	JsonMetadata        string        `json:"json_metadata"         validate:"omitempty,json"`
	PostingJsonMetadata string        `json:"posting_json_metadata" validate:"omitempty,json"`
}

// Request: POST http://0.0.0.0:3001/user
//   - request json body: userUpdateParam
//
// Response:
//   - returns 204 No Content if the user is created successfully.
//   - returns 400 Bad Request otherwise
func (h *serverHandler) updateUser(w http.ResponseWriter, r *http.Request) {
	var c userUpdateParam

	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		h.logger.Error("failed to decode request body.", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := validator.New().Struct(&c); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	ts := time.Now().Format(utils.TimeFormat)

	updateParams := accountdb.Account{
		Name:                c.Account,
		JsonMeta:            c.JsonMetadata,
		JsonPostingMetadata: c.PostingJsonMetadata,
		LastAccountUpdate:   ts,
		UserKeySet: accountdb.UserKeySet{
			Owner:   c.Owner,
			Active:  c.Active,
			Posting: c.Posting,
		},
	}

	err := h.db.UpdateAccount(ctx, &updateParams)
	if err != nil {
		if errors.Is(err, accountdb.ErrAccountNotFound) {
			http.Error(w, "account not found.", http.StatusBadRequest)
		} else {
			h.logger.Error("failed to update user key.", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

package admin

import (
	"encoding/json"
	"mimic/modules/admin/services"
	"net/http"
)

// Request:
//   - request json body: { account: string, password: string }.
//
// Response:
//   - returns 400 Bad Request if the request body is invalid
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

	if err := services.CreateUser(credentials.Account, credentials.Password); err != nil {
		h.logger.Error("failed to create user.", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

package api

import (
	"log/slog"
	"net/http"
)

const RootMsg = "go-mimic v1.0.0; Hive blockchain end to end simulation. To learn more, visit https://github.com/vsc-eco/go-mimic"

type httpHandler struct {
	logger *slog.Logger
}

func (h *httpHandler) root(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(RootMsg)); err != nil {
		h.logger.Error("failed to write response", "err", err)
	}
}

func (h *httpHandler) health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"reflect"
	"strings"
)

const rootMsg = "go-mimic v1.0.0; Hive blockchain end to end simulation. To learn more, visit https://github.com/vsc-eco/go-mimic"

type requestHandler struct {
	logger    *slog.Logger
	rpcRoutes map[string]*ServiceMethod
	services  map[string]reflect.Value
}

func (h *requestHandler) root(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(rootMsg))
}

func (h *requestHandler) health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *requestHandler) jsonrpc(w http.ResponseWriter, r *http.Request) {
	var req map[string]any
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode incoming requests.", "err", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	method, valid := req["method"].(string)

	if !valid {
		http.Error(w, "invalid method", http.StatusBadRequest)
		return
	}

	if h.rpcRoutes[method] == nil {
		http.Error(w, "method not found", http.StatusNotFound)
		return
	}

	methodSpec := h.rpcRoutes[method]

	args := reflect.New(methodSpec.argType)
	paramsJSON, err := json.Marshal(req["params"])
	if err != nil {
		http.Error(w, "invalid params", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(paramsJSON, args.Interface()); err != nil {
		slog.Error("Failed to decode params",
			"raw", paramsJSON, "err", err)
		http.Error(w, "failed to decode params", http.StatusBadRequest)
		return
	}
	reply := reflect.New(h.rpcRoutes[method].replyType)

	strs := strings.Split(method, ".")
	methodSpec.method.Func.Call([]reflect.Value{
		h.services[strs[0]],
		args,
		reply,
	})

	res := map[string]any{
		"jsonrpc": "2.0",
		"id":      req["id"],
		"result":  reply.Interface(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

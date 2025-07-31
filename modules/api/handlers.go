package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"reflect"

	"github.com/sourcegraph/jsonrpc2"
)

const rootMsg = "go-mimic v1.0.0; Hive blockchain end to end simulation. To learn more, visit https://github.com/vsc-eco/go-mimic"

type connBuffer struct {
	*bytes.Buffer
}

func (c *connBuffer) Read(p []byte) (n int, err error) {
	return c.Buffer.Read(p)
}

func (c *connBuffer) Write(p []byte) (n int, err error) {
	return c.Buffer.Write(p)
}

func (c *connBuffer) Close() error { return nil } // nop

type rpcHandler struct {
	*requestHandler
}

// Handle is called to handle a request. No other requests are handled
// until it returns. If you do not require strict ordering behavior
// of received RPCs, it is suggested to wrap your handler in
// AsyncHandler.
func (rH *rpcHandler) Handle(
	ctx context.Context,
	conn *jsonrpc2.Conn,
	req *jsonrpc2.Request,
) {
	select {
	case <-ctx.Done():
		return
	default:
		rH.requestHandler.Handle(conn, req)
		fmt.Println("Request received", req.Method, conn)
	}
}

type requestHandler struct {
	logger    *slog.Logger
	rpcRoutes map[string]*ServiceMethod
	services  map[string]reflect.Value
}

func (rH *requestHandler) root(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(rootMsg)); err != nil {
		rH.logger.Error("failed to write response", "err", err)
	}
}

func (rH *requestHandler) health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (rH *requestHandler) jsonrpc(w http.ResponseWriter, r *http.Request) {
	req := new(jsonrpc2.Request)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		rH.logger.Error("failed to decode incoming requests.", "err", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	buf := &connBuffer{new(bytes.Buffer)}
	stream := jsonrpc2.NewPlainObjectStream(buf)
	conn := jsonrpc2.NewConn(context.Background(), stream, nil)

	handler := rpcHandler{rH}
	handler.Handle(context.Background(), conn, req)
	/*

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
	*/
}

func (rH *requestHandler) Handle(
	conn *jsonrpc2.Conn,
	req *jsonrpc2.Request,
) {
}

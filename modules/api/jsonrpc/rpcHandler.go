package apijsonrpc

import (
	"encoding/json"
	"fmt"
	"log/slog"
	jsonrpcutils "mimic/lib/utils/jsonrpc"
	"net/http"
	"reflect"
	"strings"

	"github.com/sourcegraph/jsonrpc2"
)

type Handler struct {
	Logger   *slog.Logger
	Routes   map[string]*ServiceMethod
	Services map[string]reflect.Value
}

func (rpc *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	req := &jsonrpc2.Request{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		rpc.Logger.Error("failed to decode incoming requests.", "err", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	response := jsonrpc2.Response{ID: req.ID}

	result, err := rpc.rpcHandle(req)
	if err != nil {
		response.Error = err
	} else {
		err := response.SetResult(result)
		if err != nil {
			rpc.Logger.Error("failed to sererialize response", "err", err)
			response.Error = jsonrpcutils.ErrInternalServer
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&response); err != nil {
		rpc.Logger.Error("failed to write response", "err", err)
	}
}

func (rpc *Handler) rpcHandle(req *jsonrpc2.Request) (any, *jsonrpc2.Error) {
	methodSpec, ok := rpc.Routes[req.Method]
	if !ok {
		return nil, jsonrpcutils.ErrUnsupportedMethod
	}

	args := reflect.New(methodSpec.ArgType)

	if err := json.Unmarshal(*req.Params, args.Interface()); err != nil {
		rpc.Logger.Error("Failed to decode params",
			"raw", string(*req.Params),
			"err", err,
		)

		return nil, jsonrpcutils.NewInvalidRequestErr("invalid params")
	}

	strs := strings.Split(req.Method, ".")
	jsonrpcResponse := methodSpec.Method.Func.Call([]reflect.Value{
		rpc.Services[strs[0]],
		args,
	})

	if len(jsonrpcResponse) != 2 {
		msg := fmt.Sprintf(
			`invalid return type for method handler %s
	expected return type: (any, *jsonrpc2.Error)
	got: %v`,
			req.Method, jsonrpcResponse,
		)
		panic(msg)
	}

	var (
		out = jsonrpcResponse[0].Interface()
		err = jsonrpcResponse[1].Interface().(*jsonrpc2.Error)
	)

	return out, err
}

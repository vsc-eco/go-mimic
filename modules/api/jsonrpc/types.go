package apijsonrpc

import (
	"encoding/json"
	"reflect"

	"github.com/sourcegraph/jsonrpc2"
)

type Request struct {
	JsonRpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      int         `json:"id"`
}

// Response envelope from our RPC.
type Response struct {
	Result json.RawMessage `json:"result"`
	Error  string          `json:"error"`
}

type ServiceMethod struct {
	Method  reflect.Method
	ArgType reflect.Type
}

type JsonRpcResult struct {
	Result json.RawMessage
	Error  *jsonrpc2.Error
}

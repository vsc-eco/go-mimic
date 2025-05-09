package api

import (
	"encoding/json"
	"reflect"
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
	method    reflect.Method
	argType   reflect.Type
	replyType reflect.Type
}

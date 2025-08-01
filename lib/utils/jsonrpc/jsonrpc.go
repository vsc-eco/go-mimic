package jsonrpcutils

import "github.com/sourcegraph/jsonrpc2"

type Error *jsonrpc2.Error

var (
	ErrInternalServer = &jsonrpc2.Error{
		Code:    jsonrpc2.CodeInternalError,
		Message: "internal server error",
	}

	ErrUnsupportedMethod = &jsonrpc2.Error{
		Code:    jsonrpc2.CodeMethodNotFound,
		Message: "unsupported method",
	}
)

func NewInvalidRequestErr(msg string) Error {
	return &jsonrpc2.Error{
		Code:    jsonrpc2.CodeInvalidRequest,
		Message: msg,
	}
}

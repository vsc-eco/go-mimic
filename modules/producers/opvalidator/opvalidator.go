package opvalidator

import (
	"errors"
	"mimic/lib/validator"

	"github.com/vsc-eco/hivego"
)

var (
	ErrValidatorNotImplemented = errors.New("validator not implemented")
	ErrInvalidOperation        = errors.New("invalid operation")

	v = validator.New()
)

type OperationValidator interface {
	ValidateOperation(hivego.HiveOperation) error
}

func NewValidator(opName string) (OperationValidator, error) {
	switch opName {
	case "custom_json":
		return &customJsonValidator{v}, nil
	case "account_update":
		return &accountUpdateValidator{v}, nil
	default:
		return nil, ErrInvalidOperation
	}
}

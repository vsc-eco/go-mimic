package opvalidator

import (
	"errors"
	"mimic/lib/validator"

	"github.com/vsc-eco/hivego"
)

var (
	ErrValidatorNotImplemented = errors.New("validator not implemented")
	ErrInvalidOperation        = errors.New("invalid operation")
)

type OperationValidator interface {
	ValidateOperation(hivego.HiveOperation) error
}

var validatorMap map[string]OperationValidator

func init() {
	v := validator.New()

	validatorMap = map[string]OperationValidator{
		"custom_json": &customJsonValidator{v},
	}
}

func NewValidator(opName string) (OperationValidator, error) {
	v, ok := validatorMap[opName]
	if !ok {
		return nil, ErrValidatorNotImplemented
	}
	return v, nil
}

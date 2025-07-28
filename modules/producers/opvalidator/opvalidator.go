package opvalidator

import (
	"errors"
	"mimic/lib/validator"

	"github.com/vsc-eco/hivego"
)

var ErrUnimplementedValidator = errors.New("invalid operation type")

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
		return nil, errors.New("validator not implemented")
	}
	return v, nil
}

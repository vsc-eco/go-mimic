package opvalidator

import (
	"errors"

	"github.com/vsc-eco/hivego"
)

type OperationValidator interface {
	Validate(hivego.HiveOperation) error
}

var validator = map[string]OperationValidator{
	"custom_json": &customJsonValidator{},
}

func NewValidator(opName string) (OperationValidator, error) {
	v, ok := validator[opName]
	if !ok {
		return nil, errors.New("validator not implemented")
	}
	return v, nil
}

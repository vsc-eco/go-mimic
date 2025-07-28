package opvalidator

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/vsc-eco/hivego"
)

type customJsonValidator struct {
	*validator.Validate
}

func (c *customJsonValidator) ValidateOperation(o hivego.HiveOperation) error {
	op, ok := o.(*hivego.CustomJsonOperation)
	if !ok {
		return errInvalidOperationType
	}

	if len(op.RequiredAuths) == 0 && len(op.RequiredPostingAuths) == 0 {
		return errors.New("missing required auths")
	}

	return validateFields(
		c.Validate,
		fieldV{op.Json, "required,json"},
		fieldV{op.Id, "required"},
	)
}

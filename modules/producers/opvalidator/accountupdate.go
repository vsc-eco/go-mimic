package opvalidator

import (
	"github.com/go-playground/validator/v10"
	"github.com/vsc-eco/hivego"
)

type accountUpdateValidator struct {
	*validator.Validate
}

func (c *accountUpdateValidator) ValidateOperation(
	o hivego.HiveOperation,
) error {
	op, ok := o.(*hivego.AccountUpdateOperation)
	if !ok {
		return ErrInvalidOperation
	}

	return validateFields(
		c.Validate,
		fieldV{op.Account, "required"},
	)
}

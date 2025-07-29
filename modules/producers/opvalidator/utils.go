package opvalidator

import "github.com/go-playground/validator/v10"

type fieldV struct {
	field any
	tags  string
}

func validateFields(
	v *validator.Validate,
	fields ...fieldV,
) error {
	for _, f := range fields {
		if err := v.Var(f.field, f.tags); err != nil {
			return err
		}
	}

	return nil
}

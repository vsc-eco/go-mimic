package producers

import (
	"fmt"
	"mimic/modules/producers/opvalidator"

	"github.com/vsc-eco/hivego"
)

func ValidateTransaction(transaction *hivego.HiveTransaction) error {
	for _, op := range transaction.Operations {
		v, err := opvalidator.NewValidator(op.OpName())
		if err != nil {
			panic(
				fmt.Sprintf(
					"%v\nvalidator for this type [%s] isn't implemented",
					err,
					op.OpName(),
				),
			)
		}

		if err := v.Validate(op); err != nil {
			return err
		}
	}
	return nil
}

package opvalidator_test

import (
	"mimic/modules/producers/opvalidator"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vsc-eco/hivego"
)

func TestValidateAccountUpdate(t *testing.T) {
	opDefault := hivego.AccountUpdateOperation{
		Account: "foo",
	}

	v, err := opvalidator.NewValidator(opDefault.OpName())
	assert.NoError(t, err)

	t.Run("catches no account error", func(t *testing.T) {
		op := opDefault
		op.Account = ""
		assert.Error(t, v.ValidateOperation(op))
	})
}

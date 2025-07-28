package opvalidator_test

import (
	"mimic/modules/producers/opvalidator"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vsc-eco/hivego"
)

func TestValidateCustomJson(t *testing.T) {
	opDefault := hivego.CustomJsonOperation{
		RequiredAuths:        []string{"foo"},
		RequiredPostingAuths: []string{"foo"},
		Id:                   "1",
		Json:                 `{"foo": "bar"}`,
	}

	v, err := opvalidator.NewValidator("custom_json")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("catches no auth error", func(t *testing.T) {
		op := opDefault
		op.RequiredAuths = make([]string, 0)
		op.RequiredPostingAuths = make([]string, 0)
		assert.NotNil(t, v.ValidateOperation(op))
	})

	t.Run("catches invalid ID", func(t *testing.T) {
		op := opDefault
		op.Id = ""
		assert.NotNil(t, v.ValidateOperation(op))
	})

	t.Run("catches invalid json", func(t *testing.T) {
		op := opDefault
		op.Json = ""

		assert.NotNil(t, v.ValidateOperation(op))

		op.Json = "not json string"
		assert.NotNil(t, v.ValidateOperation(op))
	})
}

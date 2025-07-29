package producers

import (
	"errors"
	"mimic/lib/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vsc-eco/hivego"
)

func TestValidateTransaction(t *testing.T) {
	customJson := hivego.CustomJsonOperation{
		RequiredAuths:        []string{},
		RequiredPostingAuths: []string{utils.EnvOrPanic("TEST_USERNAME")},
		Id:                   "1",
		Json:                 `{"foo":"bar"}`,
	}

	trx := &hivego.HiveTransaction{
		RefBlockNum:    0,
		RefBlockPrefix: 0,
		Expiration:     time.Now().Format(utils.TimeFormat),
		Operations:     []hivego.HiveOperation{&customJson},
		OperationsJs:   [][2]any{{"custom_json", &customJson}},
		Extensions:     []string{},
		Signatures:     []string{},
	}

	t.Run("transaction with no signatures", func(t *testing.T) {
		err := ValidateTransaction(trx)
		assert.True(t, errors.Is(err, errMissingSignature))
	})
}

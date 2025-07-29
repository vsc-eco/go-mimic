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

	t.Run("transaction with valid signature", func(t *testing.T) {
		keyPair, err := hivego.KeyPairFromWif(
			utils.EnvOrPanic("TEST_POSTING_KEY_PRIVATE"),
		)
		if err != nil {
			t.Fatal(err)
		}

		sig, err := trx.Sign(*keyPair)
		if err != nil {
			t.Fatal(err)
		}

		trx.AddSig(sig)
		assert.NoError(t, ValidateTransaction(trx))
	})
}

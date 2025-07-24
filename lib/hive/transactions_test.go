package hive

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vsc-eco/hivego"
)

func TestTransactionJsonEncoding(t *testing.T) {
	trx := Transaction{
		HiveTransaction: &hivego.HiveTransaction{
			RefBlockNum:    12,
			RefBlockPrefix: 12,
			Expiration:     "2025-07-24T00:00:00",
			Operations: []hivego.HiveOperation{
				&hivego.CustomJsonOperation{
					RequiredAuths:        []string{"hive-io-1"},
					RequiredPostingAuths: []string{"hive-io-2"},
					Id:                   "123",
					Json:                 "{\"key\":\"value\"}",
				},
			},
			Signatures: []string{"sig1", "sig2"},
		},
	}

	jsonBytes, err := json.Marshal(&trx)
	assert.Nil(t, err)

	t.Run("tests json.Marshaler interface", func(t *testing.T) {
		var b struct {
			RefBlockNum    uint16   `json:"ref_block_num"`
			RefBlockPrefix uint32   `json:"ref_block_prefix"`
			Expiration     string   `json:"expiration"`
			Extensions     []any    `json:"extensions"`
			Operations     [][2]any `json:"operations"`
			Signatures     []string `json:"signatures"`
		}

		err = json.Unmarshal(jsonBytes, &b)
		assert.Nil(t, err)

		assert.Equal(t, trx.RefBlockNum, b.RefBlockNum)
		assert.Equal(t, trx.RefBlockPrefix, b.RefBlockPrefix)
		assert.Equal(t, trx.Expiration, b.Expiration)
		assert.Equal(t, len(trx.Operations), len(b.Operations))
		assert.Equal(t, len(trx.OperationsJs), len(b.Operations))
		assert.Equal(t, len(trx.Signatures), len(b.Signatures))
		assert.Equal(t, "custom_json", b.Operations[0][0])
	})

	t.Run("tests json.Unmarshaler interface", func(t *testing.T) {
		trx2 := &Transaction{}

		err := json.Unmarshal(jsonBytes, trx2)
		assert.Nil(t, err)

		assert.Equal(t, trx.RefBlockNum, trx2.RefBlockNum)
		assert.Equal(t, trx.RefBlockPrefix, trx2.RefBlockPrefix)
		assert.Equal(t, trx.Expiration, trx2.Expiration)
		assert.Equal(t, len(trx.Operations), len(trx2.Operations))
		assert.Equal(t, len(trx.OperationsJs), len(trx2.OperationsJs))
		assert.Equal(t, len(trx.Signatures), len(trx2.Signatures))
		assert.Equal(t, "custom_json", trx2.OperationsJs[0][0])
	})
}

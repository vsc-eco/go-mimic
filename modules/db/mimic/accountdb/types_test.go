package accountdb_test

import (
	"encoding/json"
	"mimic/modules/db/mimic/accountdb"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccountAuthJSONMarshalling(t *testing.T) {
	raw := []byte(`{
		"weight_threshold": 1000,
		"account_auths": [["hive-io", 1000]],
		"key_auths": [["hive-io-pubkey", 1000]]
	}`)

	stub := accountdb.AccountAuthority{}

	err := json.Unmarshal(raw, &stub)
	assert.Nil(t, err)
	assert.Equal(t, "hive-io", stub.AccountAuths[0].Account)
	assert.Equal(t, 1000, stub.AccountAuths[0].Weight)
	assert.Equal(t, "hive-io-pubkey", stub.KeyAuths[0].PublicKey)
	assert.Equal(t, 1000, stub.KeyAuths[0].Weight)
}

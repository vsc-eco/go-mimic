package accountdb_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vsc-eco/hivego"
)

func TestAccountAuthJSONMarshalling(t *testing.T) {
	raw := []byte(`{
		"weight_threshold": 1000,
		"account_auths": [["hive-io", 1000]],
		"key_auths": [["hive-io-pubkey", 1000]]
	}`)

	stub := hivego.Auths{}

	err := json.Unmarshal(raw, &stub)
	assert.Nil(t, err)
	assert.Equal(t, "hive-io", stub.AccountAuths[0][0])
	assert.Equal(t, 1000, stub.AccountAuths[0][1])
	assert.Equal(t, "hive-io-pubkey", stub.KeyAuths[0][0])
	assert.Equal(t, 1000, stub.KeyAuths[0][1])
}

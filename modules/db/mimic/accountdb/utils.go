package accountdb

import (
	"fmt"
	"mimic/mock"
	"mimic/modules/crypto"

	"github.com/vsc-eco/hivego"
)

type privateKeyMapType = map[string]crypto.HiveKeySet

var privateKeyMap privateKeyMapType = make(privateKeyMapType)

type accountPrivateKeys struct {
	OwnerKey   string `json:"owner_key,omitempty"`
	PostingKey string `json:"posting_key,omitempty"`
	ActiveKey  string `json:"active_key,omitempty"`
}

func GetPrivateKey(username string) (*crypto.HiveKeySet, error) {
	k, ok := privateKeyMap[username]
	if !ok {
		return nil, fmt.Errorf("Private key not loaded for %s.", username)
	}
	return &k, nil
}

// generate the keys for the seed account, the private keys are stored in
// memory, writing to the global variable `privateKeyMap`
func GetSeedAccounts() ([]Account, error) {
	accounts, err := mock.LoadSeedUserCredentials()
	if err != nil {
		return nil, err
	}

	accountBuf := make([]Account, len(accounts))

	for i, account := range accounts {
		username, password := account.Username, account.Password

		keySet, err := crypto.MakeHiveKeySet(username, password)
		if err != nil {
			return nil, err
		}

		privateKeyMap[username] = *keySet

		accountBuf[i] = Account{
			Name: username,
			Active: hivego.Auths{
				WeightThreshold: 1,
				AccountAuths:    [][2]any{{username, 1}},
				KeyAuths: [][2]any{
					{*keySet.ActiveKey().GetPublicKeyString(), 1},
				},
			},

			Owner: hivego.Auths{
				WeightThreshold: 1,
				AccountAuths:    [][2]any{{username, 1}},
				KeyAuths: [][2]any{
					{*keySet.OwnerKey().GetPublicKeyString(), 1},
				},
			},

			Posting: hivego.Auths{
				WeightThreshold: 1,
				AccountAuths:    [][2]any{{username, 1}},
				KeyAuths: [][2]any{
					{*keySet.PostingKey().GetPublicKeyString(), 1},
				},
			},

			MemoKey: keySet.MemoKey(),
		}
	}

	return accountBuf, nil
}

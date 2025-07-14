package accountdb

import (
	"fmt"
	"mimic/mock"
	"mimic/modules/crypto"
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

		keySet := crypto.MakeHiveKeySet(username, password)
		privateKeyMap[username] = keySet

		accountBuf[i] = Account{
			Name: username,
			Active: AccountAuthority{
				WeightThreshold: 1,
				AccountAuths: []AccountAuth{{
					Account: username,
					Weight:  1,
				}},
				KeyAuths: []KeyAuth{{
					PublicKey: *keySet.ActiveKey().GetPublicKeyString(),
					Weight:    1,
				}},
			},

			Owner: AccountAuthority{
				WeightThreshold: 1,
				AccountAuths: []AccountAuth{{
					Account: username,
					Weight:  1,
				}},
				KeyAuths: []KeyAuth{{
					PublicKey: *keySet.OwnerKey().GetPublicKeyString(),
					Weight:    1,
				}},
			},

			Posting: AccountAuthority{
				WeightThreshold: 1,
				AccountAuths: []AccountAuth{{
					Account: username,
					Weight:  1,
				}},
				KeyAuths: []KeyAuth{{
					PublicKey: *keySet.PostingKey().GetPublicKeyString(),
					Weight:    1,
				}},
			},

			MemoKey: keySet.MemoKey(),
		}
	}

	return accountBuf, nil
}

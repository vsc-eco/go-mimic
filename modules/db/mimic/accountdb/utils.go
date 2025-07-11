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

// the call to `GetSeedAccounts` stores the users' private keys in a
// global map for this function to read.
func GetPrivateKey(username string) (*crypto.HiveKeySet, error) {
	k, ok := privateKeyMap[username]
	if !ok {
		return nil, fmt.Errorf("user not found %s.", username)
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
				AccountAuths:    []AccountAuth{},
				KeyAuths: []KeyAuth{{
					PublicKey: keySet.ActiveKey().PublicKeyWIF(),
					Weight:    1,
				}},
			},

			Owner: AccountAuthority{
				WeightThreshold: 1,
				AccountAuths:    []AccountAuth{},
				KeyAuths: []KeyAuth{{
					PublicKey: keySet.OwnerKey().PublicKeyWIF(),
					Weight:    1,
				}},
			},

			Posting: AccountAuthority{
				WeightThreshold: 1,
				AccountAuths:    []AccountAuth{},
				KeyAuths: []KeyAuth{{
					PublicKey: keySet.PostingKey().PublicKeyWIF(),
					Weight:    1,
				}},
			},

			MemoKey: keySet.MemoKey().PublicKeyWIF(),
		}
	}

	return accountBuf, nil
}

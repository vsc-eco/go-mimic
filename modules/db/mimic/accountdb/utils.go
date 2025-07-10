package accountdb

import (
	"fmt"
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

func makeAccount(username, password string) (Account, crypto.HiveKeySet) {
	keySet := crypto.MakeHiveKeySet(username, password)
	account := Account{
		Name: username,
		Active: AccountAuthority{
			WeightThreshold: 1,
			AccountAuths: []AccountAuth{{
				Account: username,
				Weight:  1,
			}},
			KeyAuths: []KeyAuth{{
				PublicKey: keySet.ActiveKey().PublicKeyHex(),
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
				PublicKey: keySet.OwnerKey().PublicKeyHex(),
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
				PublicKey: keySet.PostingKey().PublicKeyHex(),
				Weight:    1,
			}},
		},

		MemoKey: keySet.MemoKey(),
	}

	return account, keySet
}

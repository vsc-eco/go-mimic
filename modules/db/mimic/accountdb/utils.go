package accountdb

import (
	"fmt"
	"mimic/lib/hive"
	"mimic/lib/utils"
	"time"

	"github.com/vsc-eco/hivego"
)

type privateKeyMapType = map[string]hive.HiveKeySet

var privateKeyMap privateKeyMapType = make(privateKeyMapType)

func GetPrivateKey(username string) (*hive.HiveKeySet, error) {
	k, ok := privateKeyMap[username]
	if !ok {
		return nil, fmt.Errorf("Private key not loaded for %s", username)
	}
	return &k, nil
}

func GetSeedAccounts() ([]Account, error) {
	var (
		username         = utils.EnvOrPanic("TEST_USERNAME")
		activePubKeyWif  = utils.EnvOrPanic("TEST_ACTIVE_KEY_PUBLIC")
		ownerPubKeyWif   = utils.EnvOrPanic("TEST_OWNER_KEY_PUBLIC")
		postingPubKeyWif = utils.EnvOrPanic("TEST_POSTING_KEY_PUBLIC")
	)

	ts := time.Now().Format(utils.TimeFormat)
	account := Account{
		Name: username,
		UserKeySet: UserKeySet{
			Owner: &hivego.Auths{
				WeightThreshold: 1,
				AccountAuths:    [][2]any{{username, 1}},
				KeyAuths:        [][2]any{{ownerPubKeyWif, 1}},
			},
			Active: &hivego.Auths{
				WeightThreshold: 1,
				AccountAuths:    [][2]any{{username, 1}},
				KeyAuths:        [][2]any{{activePubKeyWif, 1}},
			},
			Posting: &hivego.Auths{
				WeightThreshold: 1,
				AccountAuths:    [][2]any{{username, 1}},
				KeyAuths:        [][2]any{{postingPubKeyWif, 1}},
			},
		},
		Created:           ts,
		LastOwnerUpdate:   ts,
		LastAccountUpdate: ts,
	}

	return []Account{account}, nil
}

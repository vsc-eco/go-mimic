package services

import (
	"errors"
	"mimic/lib/hivekey"
	"mimic/lib/utils"
	"mimic/modules/db/mimic/accountdb"
	"time"

	"github.com/vsc-eco/hivego"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func MakeAccount(username, password string) (*accountdb.Account, error) {
	if len(username) == 0 || len(password) == 0 {
		return nil, errors.New("invalid account or password")
	}

	keySet := hivekey.MakeHiveKeySet(username, password)

	account := &accountdb.Account{
		ObjectId: primitive.NilObjectID,
		Id:       0,
		Name:     username,
		Owner: hivego.Auths{
			WeightThreshold: 10000,
			AccountAuths:    [][2]any{},
			KeyAuths: [][2]any{
				{keySet.OwnerKey().GetPublicKeyString(), 1},
			},
		},
		Active: hivego.Auths{
			WeightThreshold: 10000,
			AccountAuths:    [][2]any{},
			KeyAuths: [][2]any{
				{keySet.ActiveKey().GetPublicKeyString(), 1},
			},
		},
		Posting: hivego.Auths{
			WeightThreshold: 10000,
			AccountAuths:    [][2]any{},
			KeyAuths: [][2]any{
				{keySet.PostingKey().GetPublicKeyString(), 1},
			},
		},
		MemoKey:             "",
		JsonMeta:            "",
		JsonPostingMetadata: "",
		LastOwnerUpdate:     time.Now().Format(utils.TimeFormat),
		LastAccountUpdate:   time.Now().Format(utils.TimeFormat),
		Created:             time.Now().Format(utils.TimeFormat),
		Balance:             "",
		HbdBalance:          "",
		SavingsHbdBalance:   "",
		VestingShares:       "",
		Reputation:          0,

		PrivateKeys: accountdb.PrivateKeys{
			OwnerKey:   keySet.OwnerKey().PrivateKeyWif(),
			ActiveKey:  keySet.ActiveKey().PrivateKeyWif(),
			PostingKey: keySet.PostingKey().PrivateKeyWif(),
		},
	}

	return account, nil
}

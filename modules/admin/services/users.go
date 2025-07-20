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

	keySet := MakeNewUserKey(username, password)

	account := &accountdb.Account{
		ObjectId: primitive.NilObjectID,
		Name:     username,
		KeySet: accountdb.UserKeySet{
			Owner:         keySet.Owner,
			Active:        keySet.Active,
			Posting:       keySet.Posting,
			PrivateKeySet: keySet.PrivateKeySet,
		},
		LastOwnerUpdate:   time.Now().Format(utils.TimeFormat),
		LastAccountUpdate: time.Now().Format(utils.TimeFormat),
		Created:           time.Now().Format(utils.TimeFormat),
	}

	return account, nil
}

func MakeNewUserKey(account, password string) accountdb.UserKeySet {
	keySet := hivekey.MakeHiveKeySet(account, password)

	userKeySet := accountdb.UserKeySet{
		Owner: hivego.Auths{
			WeightThreshold: 0,
			AccountAuths:    [][2]any{},
			KeyAuths: [][2]any{
				{keySet.OwnerKey().GetPublicKeyString(), 1},
			},
		},
		Active: hivego.Auths{
			WeightThreshold: 0,
			AccountAuths:    [][2]any{},
			KeyAuths: [][2]any{
				{keySet.ActiveKey().GetPublicKeyString(), 1},
			},
		},
		Posting: hivego.Auths{
			WeightThreshold: 0,
			AccountAuths:    [][2]any{},
			KeyAuths: [][2]any{
				{keySet.PostingKey().GetPublicKeyString(), 1},
			},
		},
		PrivateKeySet: accountdb.PrivateKeys{
			OwnerKey:   keySet.OwnerKey().PrivateKeyWif(),
			ActiveKey:  keySet.ActiveKey().PrivateKeyWif(),
			PostingKey: keySet.PostingKey().PrivateKeyWif(),
		},
	}

	return userKeySet
}

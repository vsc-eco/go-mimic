package accountdb

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"mimic/lib/hive"

	"github.com/decred/dcrd/dcrec/secp256k1/v2"
	"github.com/vsc-eco/hivego"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (db *AccountDB) InsertAccount(
	ctx context.Context,
	account *Account,
) error {
	if _, err := db.collection.InsertOne(ctx, account); err != nil {
		return err
	}
	slog.Debug("accounts inserted", "object-id", account.ObjectId)
	return nil
}

func (a *AccountDB) QueryAccountByNames(
	ctx context.Context,
	buf *[]Account,
	names []string,
) error {
	filter := bson.M{"name": bson.M{"$in": names}}

	cursor, err := a.collection.Find(ctx, filter)
	if err != nil {
		return err
	}

	defer cursor.Close(ctx)

	return cursor.All(ctx, buf)
}

func (a *AccountDB) UpdateAccount(
	ctx context.Context,
	account *Account,
) error {
	if len(account.Name) == 0 {
		return errors.New("required field `account.Name` not set")
	}

	updateDoc := bson.M{"last_account_update": account.LastAccountUpdate}
	if account.Owner != nil {
		updateDoc["owner"] = account.Owner
	}
	if account.Active != nil {
		updateDoc["active"] = account.Active
	}
	if account.Posting != nil {
		updateDoc["posting"] = account.Posting
	}
	if len(account.JsonMeta) != 0 {
		updateDoc["json_metadata"] = account.JsonMeta
	}
	if len(account.JsonPostingMetadata) != 0 {
		updateDoc["posting_json_metadata"] = account.JsonPostingMetadata
	}

	filter := bson.M{"name": account.Name}
	update := bson.M{"$set": updateDoc}

	result, err := a.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return ErrAccountNotFound
	}

	return nil
}

func (a *AccountDB) QueryAccountByPubKeyWIF(
	ctx context.Context,
	account *Account,
	keyString string,
) error {
	filter := bson.M{
		"$or": []bson.M{
			{
				"posting.key_auths": bson.M{
					"$elemMatch": bson.M{"$eq": keyString},
				},
			},
			{
				"active.key_auths": bson.M{
					"$elemMatch": bson.M{"$eq": keyString},
				},
			},
			{
				"owner.key_auths": bson.M{
					"$elemMatch": bson.M{"$eq": keyString},
				},
			},
		},
	}

	result := a.collection.FindOne(ctx, filter)

	return result.Decode(account)
}

func (a *AccountDB) QueryPubKeysByAccount(
	ctx context.Context,
	keyBuf map[string]map[hive.KeyRole]*secp256k1.PublicKey,
	account []string,
) error {
	opts := options.Find()
	opts.SetProjection(bson.M{
		"name":              1,
		"active.key_auths":  1,
		"owner.key_auths":   1,
		"posting.key_auths": 1,
	})

	filter := bson.M{
		"name": bson.M{
			"$in": account,
		},
	}

	cur, err := a.collection.Find(ctx, filter, opts)
	if err != nil {
		return err
	}

	var buf []Account
	if err := cur.All(ctx, &buf); err != nil {
		return err
	}

	for _, account := range buf {
		var (
			_, hasOwner   = keyBuf[account.Name][hive.OwnerKeyRole]
			_, hasActive  = keyBuf[account.Name][hive.ActiveKeyRole]
			_, hasPosting = keyBuf[account.Name][hive.PostingKeyRole]

			err error
		)

		if hasOwner {
			keyBuf[account.Name][hive.OwnerKeyRole], err = hivego.DecodePublicKey(
				account.Owner.KeyAuths[0][0].(string),
			)
			if err != nil {
				return fmt.Errorf(
					"failed to decode owner public key for account %s: %w",
					account.Name,
					err,
				)
			}
		}

		if hasActive {
			keyBuf[account.Name][hive.ActiveKeyRole], err = hivego.DecodePublicKey(
				account.Active.KeyAuths[0][0].(string),
			)
			if err != nil {
				return fmt.Errorf(
					"failed to decode active public key for account %s: %w",
					account.Name,
					err,
				)
			}
		}

		if hasPosting {
			keyBuf[account.Name][hive.PostingKeyRole], err = hivego.DecodePublicKey(
				account.Posting.KeyAuths[0][0].(string),
			)
			if err != nil {
				return fmt.Errorf(
					"failed to decode posting public key for account %s: %w",
					account.Name,
					err,
				)
			}
		}
	}

	return nil
}

package accountdb

import (
	"context"
	"errors"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
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
	if account.KeySet.Owner != nil {
		updateDoc["owner"] = account.KeySet.Owner
	}
	if account.KeySet.Active != nil {
		updateDoc["active"] = account.KeySet.Active
	}
	if account.KeySet.Posting != nil {
		updateDoc["posting"] = account.KeySet.Posting
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

package accountdb

import (
	"context"
	"log/slog"
	"mimic/lib/utils"
	"time"

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

func (a *AccountDB) UpdateAccountKeySet(
	ctx context.Context,
	accountName string,
	newKeySet *UserKeySet,
) error {
	ts := time.Now().Format(utils.TimeFormat)

	doc := struct {
		NewKeySet         *UserKeySet `bson:",inline"`
		LastOwnerUpdate   string      `bson:"last_owner_update"`
		LastAccountUpdate string      `bson:"last_account_update"`
	}{
		NewKeySet:         newKeySet,
		LastOwnerUpdate:   ts,
		LastAccountUpdate: ts,
	}

	filter := bson.M{"name": accountName}
	update := bson.M{"$set": doc}

	result, err := a.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return ErrAccountNotFound
	}

	return nil
}

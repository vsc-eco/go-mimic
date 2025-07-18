package accountdb

import (
	"context"
	"log/slog"
	"mimic/lib/hivekey"

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
	newKeySet *hivekey.HiveKeySet,
) error {
	return nil
}

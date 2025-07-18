package accountdb

import (
	"context"
	"log/slog"
	"mimic/lib/hivekey"
	"mimic/lib/utils"
	"mimic/modules/db"
	"time"

	"github.com/chebyrash/promise"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AccountQuery interface {
	InsertAccount(context.Context, *Account) error
	UpdateAccountKeySet(context.Context, string, *hivekey.HiveKeySet) error
	QueryAccountByNames(context.Context, *[]Account, []string) error
}

type AccountDB struct {
	collection *mongo.Collection
}

var collection AccountQuery = &AccountDB{nil}

func Collection() *AccountDB {
	return collection.(*AccountDB)
}

// Runs initialization in order of how they are passed in to `Aggregate`
func (accountdb *AccountDB) Init() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	db.CreateIndex(ctx, accountdb.collection, mongo.IndexModel{
		Keys:    bson.D{{Key: "name", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("name_unique"),
	})

	// load seed user
	accounts, err := GetSeedAccounts()
	if err != nil {
		return err
	}

	documents := make([]any, len(accounts))
	for i, a := range accounts {
		documents[i] = a
	}

	result, err := accountdb.collection.InsertMany(ctx, documents)
	if err != nil && !mongo.IsDuplicateKeyError(err) {
		return err
	}

	slog.Debug(
		"Seeded collection.",
		"collection",
		accountdb.collection.Name(),
		"documents",
		len(result.InsertedIDs),
	)

	slog.Debug("Mock private keys loaded.", "key-num", len(privateKeyMap))

	return nil
}

// Runs startup and should be non blocking
func (accountdb *AccountDB) Start() *promise.Promise[any] {
	return utils.PromiseResolve[any](nil)
}

// Runs cleanup once the `Aggregate` is finished
func (accountdb *AccountDB) Stop() error {
	return nil
}

func New(d *mongo.Database) *AccountDB {
	collection.(*AccountDB).collection = db.NewCollection(d, "accounts")
	return collection.(*AccountDB)
}

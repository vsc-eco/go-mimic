package accountdb

import (
	"context"
	"log/slog"
	"mimic/lib/utils"
	"mimic/modules/db"
	"mimic/modules/db/mimic"
	"time"

	"github.com/chebyrash/promise"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AccountDBQueries interface{}

type AccountDB struct {
	*mongo.Collection
}

var collection AccountDBQueries = &AccountDB{nil}

func Collection() *AccountDB {
	return collection.(*AccountDB)
}

// Runs initialization in order of how they are passed in to `Aggregate`
func (accountdb *AccountDB) Init() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	db.CreateIndex(ctx, accountdb.Collection, mongo.IndexModel{
		Keys:    bson.D{{Key: "name", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("name_unique"),
	})

	accounts := [...][2]string{
		{"go-mimic-root-username", "go-mimic-root-password"},
		{"alice", "alice-password"},
		{"bob", "bob-password"},
	}

	documents := make([]any, len(accounts))
	for i, a := range accounts {
		documents[i], privateKeyMap[a[0]] = makeAccount(a[0], a[1])
	}

	result, err := accountdb.Collection.InsertMany(ctx, documents)
	if err != nil && !mongo.IsDuplicateKeyError(err) {
		return err
	}

	slog.Debug(
		"Seeded collection.",
		"collection",
		accountdb.Collection.Name(),
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

func New(d *mimic.MimicDb) *AccountDB {
	collection.(*AccountDB).Collection = db.NewCollection(
		d.DbInstance,
		"accounts",
	)
	return collection.(*AccountDB)
}

func (a *AccountDB) QueryAccountByNames(
	ctx context.Context,
	buf *[]Account,
	names []string,
) error {
	filter := bson.M{"name": bson.M{"$in": names}}

	cursor, err := a.Find(ctx, filter)
	if err != nil {
		return err
	}

	defer cursor.Close(ctx)

	return cursor.All(ctx, buf)
}

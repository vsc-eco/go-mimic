package condenserdb

import (
	"context"
	"log/slog"
	"mimic/lib/utils"
	"mimic/mock"
	"mimic/modules/db"
	"mimic/modules/db/mimic"
	"time"

	"github.com/chebyrash/promise"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Condenser struct {
	*mongo.Collection
}

var condenserDb = &Condenser{nil}

func Collection() *Condenser {
	return condenserDb
}

func New(d *mimic.MimicDb) *Condenser {
	condenserDb.Collection = db.NewCollection(d.DbInstance, "condenser")
	return condenserDb
}

// Condenser implements `aggregate.Plugin`
func (c *Condenser) Init() error {
	indexName, err := c.Collection.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    bson.D{{Key: "name", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("name_unique"),
	})

	if err != nil {
		slog.Info("Failed to create index.", "collection", c.Name(), "err", err)
		return err
	}

	slog.Info("Index created.", "collection", c.Name(), "index", indexName)

	return nil
}

func (c *Condenser) Start() *promise.Promise[any] {
	data, err := mock.GetMockData[Account]("mock/condenser_api_get_accounts.mock.json")
	if err != nil {
		panic(err)
	}

	entries := make([]any, 0, len(data))
	for _, v := range data {
		entries = append(entries, v)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	result, err := c.InsertMany(ctx, entries)
	if err != nil && !mongo.IsDuplicateKeyError(err) {
		slog.Error("Failed to seed collection.", "err", err)
	} else {
		slog.Info("Seed collection.", "collection", c.Name(), "new-record", len(result.InsertedIDs))
	}

	return utils.PromiseResolve[any](nil)
}

func (c *Condenser) Stop() error {
	return nil
}

// Queries

func (c *Condenser) QueryGetAccounts(namedQueries []string) ([]Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	filter := bson.M{"name": bson.M{"$in": namedQueries}}

	cursor, err := c.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var accounts []Account
	if err := cursor.All(ctx, &accounts); err != nil {
		return nil, err
	}

	return accounts, nil
}

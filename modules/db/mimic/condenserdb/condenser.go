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

func New(d *mimic.MimicDb) *Condenser {
	condenserDb.Collection = db.NewCollection(d.DbInstance, "condenser")
	return condenserDb
}

// Init implements aggregate.Plugin.
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

// Start implements aggregate.Plugin.
func (c *Condenser) Start() *promise.Promise[any] {
	data, err := mock.GetMockData[Account]("mock/condenser_api_get_accounts.mock.json")
	if err != nil {
		panic(err)
	}

	entries := make([]any, 0, len(data))
	for k, v := range data {
		slog.Info("Seeding account.", "name", k)
		entries = append(entries, v)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	result, err := c.InsertMany(ctx, entries)
	if err != nil && !mongo.IsDuplicateKeyError(err) {
		slog.Error("Failed to seed collection.", "err", err)
	} else {
		slog.Info("Seed condenser collection.", "seed", result)
	}

	return utils.PromiseResolve[any](nil)
}

// Stop implements aggregate.Plugin.
func (c *Condenser) Stop() error {
	return nil
}

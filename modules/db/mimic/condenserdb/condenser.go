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
	accounts *mongo.Collection
	orders   *mongo.Collection
}

var condenserDb = &Condenser{nil, nil}

func Collection() *Condenser {
	return condenserDb
}

func New(d *mimic.MimicDb) *Condenser {
	condenserDb.accounts = db.NewCollection(d.DbInstance, "accounts")
	condenserDb.orders = db.NewCollection(d.DbInstance, "orders")
	return condenserDb
}

// Condenser implements `aggregate.Plugin`
func (c *Condenser) Init() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	db.CreateIndex(ctx, c.accounts, mongo.IndexModel{
		Keys:    bson.D{{Key: "name", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("name_unique"),
	})

	db.CreateIndex(ctx, c.orders, mongo.IndexModel{
		Keys:    bson.D{{Key: "id", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("id_unique"),
	})
	return nil
}

func (c *Condenser) Start() *promise.Promise[any] {
	data := make(map[string]Account)
	err := mock.GetMockData(&data, "condenser_api.get_accounts")
	if err != nil {
		panic(err)
	}

	entries := make([]any, 0, len(data))
	for _, v := range data {
		entries = append(entries, v)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	result, err := c.accounts.InsertMany(ctx, entries)
	if err != nil && !mongo.IsDuplicateKeyError(err) {
		slog.Error("Failed to seed collection.", "err", err)
	} else {
		slog.Info("Seed collection.", "collection", c.accounts.Name(), "new-record", len(result.InsertedIDs))
	}

	return utils.PromiseResolve[any](nil)
}

func (c *Condenser) Stop() error {
	return nil
}

// Queries

func (c *Condenser) QueryGetAccounts(accounts *[]Account, namedQueries []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	filter := bson.M{"name": bson.M{"$in": namedQueries}}

	cursor, err := c.accounts.Find(ctx, filter)
	if err != nil {
		return err
	}

	defer cursor.Close(ctx)

	return cursor.All(ctx, accounts)
}

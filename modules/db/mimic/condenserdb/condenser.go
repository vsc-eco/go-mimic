package condenserdb

import (
	"context"
	"log/slog"
	"mimic/lib/utils"
	"mimic/modules/db"
	"mimic/modules/db/mimic"
	"time"

	"github.com/chebyrash/promise"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Condenser struct {
	orders       *mongo.Collection
	transactions *mongo.Collection
}

var condenserDb = &Condenser{nil, nil}

func Collection() *Condenser {
	return condenserDb
}

func New(d *mimic.MimicDb) *Condenser {
	condenserDb.orders = db.NewCollection(d.DbInstance, "orders")
	condenserDb.transactions = db.NewCollection(d.DbInstance, "transactions")
	return condenserDb
}

// Condenser implements `aggregate.Plugin`
func (c *Condenser) Init() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	db.CreateIndex(ctx, c.orders, mongo.IndexModel{
		Keys:    bson.D{{Key: "id", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("id_unique"),
	})

	db.CreateIndex(ctx, c.transactions, mongo.IndexModel{
		Keys: bson.D{
			{Key: "ref_block_prefix", Value: 1},
			{Key: "ref_block_num", Value: 1},
		},
		Options: options.Index().SetUnique(true).SetName("trx_unique"),
	})

	return nil
}

func (c *Condenser) Start() *promise.Promise[any] {
	return utils.PromiseResolve[any](nil)
}

func (c *Condenser) Stop() error {
	return nil
}

// Read Queries

// Write Queries

// Save `trx` to database, then return the transaction number, and writes the
// inserted id to `trx`
func (c *Condenser) NewTransaction(
	ctx context.Context,
	trx *Transaction,
) (int64, error) {
	// write to db
	result, err := c.transactions.InsertOne(ctx, trx)
	if err != nil {
		return 0, err
	}

	trx.ObjectID = result.InsertedID.(primitive.ObjectID)

	slog.Debug(
		"Transaction created.",
		"block-prefix", trx.RefBlockPrefix,
		"block-num", trx.RefBlockNum,
		"trx-id", trx.ObjectID,
	)

	return c.transactions.CountDocuments(ctx, bson.M{})
}

package condenserdb

import (
	"context"
	"mimic/lib/utils"
	"mimic/modules/db"
	"mimic/modules/db/mimic"
	"time"

	"github.com/chebyrash/promise"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Condenser struct {
	orders *mongo.Collection
}

var condenserDb = &Condenser{nil}

func Collection() *Condenser {
	return condenserDb
}

func New(d *mimic.MimicDb) *Condenser {
	condenserDb.orders = db.NewCollection(d.DbInstance, "orders")
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

	return nil
}

func (c *Condenser) Start() *promise.Promise[any] {
	return utils.PromiseResolve[any](nil)
}

func (c *Condenser) Stop() error {
	return nil
}

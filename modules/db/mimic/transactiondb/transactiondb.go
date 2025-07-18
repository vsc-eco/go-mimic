package transactiondb

import (
	"context"
	"mimic/lib/utils"
	"mimic/modules/db"
	"time"

	"github.com/chebyrash/promise"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var trxDb = &TransactionCollection{nil}

type TransactionCollection struct {
	*mongo.Collection
}

// Runs initialization in order of how they are passed in to `Aggregate`
func (transactioncollection *TransactionCollection) Init() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	db.CreateIndex(ctx, trxDb.Collection, mongo.IndexModel{
		Keys: bson.D{
			{Key: "ref_block_prefix", Value: 1},
			{Key: "ref_block_num", Value: 1},
		},
		Options: options.Index().SetUnique(true).SetName("trx_unique"),
	})

	return nil
}

// Runs startup and should be non blocking
func (transactioncollection *TransactionCollection) Start() *promise.Promise[any] {
	return utils.PromiseResolve[any](nil)
}

// Runs cleanup once the `Aggregate` is finished
func (transactioncollection *TransactionCollection) Stop() error {
	return nil
}

func New(d *mongo.Database) *TransactionCollection {
	trxDb.Collection = db.NewCollection(d, "transactions")
	return trxDb
}

func Collection() *TransactionCollection {
	return trxDb
}

/*
// Save `trx` to database, then return the transaction number, and writes the
// inserted id to `trx`
func (c *TransactionCollection) NewTransaction(
	ctx context.Context,
	trx *Transaction,
) (int64, error) {
	// write to db
	result, err := trxDb.InsertOne(ctx, trx)
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

	return trxDb.CountDocuments(ctx, bson.M{})
}

*/

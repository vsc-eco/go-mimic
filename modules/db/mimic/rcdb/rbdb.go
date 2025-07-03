package rcdb

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

type RcCollection struct {
	rc *mongo.Collection
}

var rcDb = &RcCollection{nil}

func Collection() *RcCollection {
	return rcDb
}

func New(d *mimic.MimicDb) *RcCollection {
	rcDb.rc = db.NewCollection(d.DbInstance, "rc")
	return rcDb
}

// Condenser implements `aggregate.Plugin`
func (c *RcCollection) Init() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	db.CreateIndex(ctx, c.rc, mongo.IndexModel{
		Keys:    bson.D{{Key: "account", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("account_unique"),
	})

	return nil
}

func (c *RcCollection) Start() *promise.Promise[any] {
	var (
		buf      = make([]RCAccount, 0)
		mockFile = "rc_api.find_rc_accounts.json"
	)

	db.Seed(&buf, context.TODO(), c.rc, mockFile)
	return utils.PromiseResolve[any](nil)
}

func (c *RcCollection) Stop() error {
	return nil
}

// Queries

func (r *RcCollection) QueryFindRcAccounts(
	ctx context.Context,
	accountBuf *[]RCAccount,
	accountQueries []string,
) error {
	filter := bson.M{"account": bson.M{"$in": accountQueries}}

	cursor, err := r.rc.Find(ctx, filter)
	if err != nil {
		return err
	}

	defer cursor.Close(ctx)

	return cursor.All(ctx, accountBuf)
}

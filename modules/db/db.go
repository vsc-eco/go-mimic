package db

import (
	"context"
	"log/slog"
	"mimic/lib/utils"
	a "mimic/modules/aggregate"

	"github.com/chebyrash/promise"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Db interface {
	Database(name string, opts ...*options.DatabaseOptions) *mongo.Database
}
type db struct {
	conf   DbConfig
	cancel context.CancelFunc
	*mongo.Client
}

var _ a.Plugin = &db{}
var _ Db = &db{}

func New(conf DbConfig) *db {
	return &db{conf: conf}
}

func (db *db) Init() error {
	ctx, cancel := context.WithCancel(context.Background())
	db.cancel = cancel

	uri := db.conf.Get().DbURI
	c, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}
	err = c.Ping(ctx, nil)
	if err != nil {
		return err
	}
	db.Client = c

	slog.Info("Connected to MongoDB.", "url", uri)

	return nil
}

func (db *db) Start() *promise.Promise[any] {
	return utils.PromiseResolve[any](db)
}

func (db *db) Stop() error {
	db.cancel()
	return nil
}

// some helper functions

func CreateIndex(ctx context.Context, col *mongo.Collection, indexModel mongo.IndexModel) {
	indexName, err := col.Indexes().CreateOne(ctx, indexModel)

	if err != nil {
		slog.Info("Failed to create index.", "collection", col.Name(), "err", err)
	} else {
		slog.Info("Index created.", "collection", col.Name(), "index", indexName)
	}
}

package db

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"mimic/lib/utils"
	a "mimic/modules/aggregate"
	"os"
	"time"

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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
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

func CreateIndex(
	ctx context.Context,
	col *mongo.Collection,
	indexModel mongo.IndexModel,
) {
	indexName, err := col.Indexes().CreateOne(ctx, indexModel)

	if err != nil {
		slog.Error(
			"Failed to create index.",
			"collection",
			col.Name(),
			"err",
			err,
		)
	} else {
		slog.Debug("Index created.", "collection", col.Name(), "index", indexName)
	}
}

func Seed[T any](
	buf *[]T,
	ctx context.Context,
	collection *mongo.Collection,
	mockJsonFile string,
) {
	var (
		seedError error
		result    *mongo.InsertManyResult
		docs      []any
	)

	f, err := os.Open(fmt.Sprintf("mock/%s", mockJsonFile))
	if err != nil {
		seedError = err
		goto seedLogging
	}

	defer f.Close()

	if err := json.NewDecoder(f).Decode(buf); err != nil {
		seedError = err
		goto seedLogging
	}

	docs = make([]any, len(*buf))
	for i, v := range *buf {
		docs[i] = v
	}

	result, seedError = collection.InsertMany(ctx, docs)

seedLogging:
	if seedError != nil && !mongo.IsDuplicateKeyError(seedError) {
		slog.Error(
			"Failed to seed collection.",
			"collection", collection.Name(),
			"mock-file", mockJsonFile,
			"err", seedError,
		)
	} else {
		slog.Debug(
			"Seeded collection.",
			"collection", collection.Name(),
			"mock-file", mockJsonFile,
			"new-record", len(result.InsertedIDs),
		)
	}
}

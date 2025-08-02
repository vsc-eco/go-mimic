package db

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
)

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

	defer f.Close() // nolint:errcheck

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

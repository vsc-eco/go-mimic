package db

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewCollection(
	db *mongo.Database,
	name string,
	opts ...*options.CollectionOptions,
) *mongo.Collection {
	defaultOpts := []*options.CollectionOptions{
		options.Collection().SetBSONOptions(&options.BSONOptions{
			UseJSONStructTags: true,
		}),
	}
	defaultOpts = append(defaultOpts, opts...)
	return db.Collection(name, defaultOpts...)
}

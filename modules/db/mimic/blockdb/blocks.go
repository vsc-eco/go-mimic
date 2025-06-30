package blockdb

import (
	"context"
	"log/slog"
	"mimic/lib/utils"
	"mimic/modules/db/mimic"

	"github.com/chebyrash/promise"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Blocks struct {
	*mongo.Collection
}

var blockCollection = &Blocks{}

func New(d *mimic.MimicDb) *Blocks {
	blockCollection.Collection = d.DbInstance.Collection("blocks")
	return blockCollection
}

func (d *Blocks) Init() error {
	indexName, err := d.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    bson.D{{Key: "name", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("name_unique"),
	})

	if err != nil {
		slog.Info("Failed to create index.", "collection", d.Name(), "err", err)
		return err
	}

	slog.Info("Index created.", "collection", d.Name(), "index", indexName)
	return nil
}

func (d *Blocks) Start() *promise.Promise[any] {
	return utils.PromiseResolve[any](d)
}

func (d *Blocks) Stop() error {
	return nil
}

func (blks *Blocks) GetBlockRange(startHeight int64, endHeight int64) []HiveBlock {
	return nil
}

func (blks *Blocks) GetBlockById(id string) HiveBlock {
	return HiveBlock{}
}

func (blks *Blocks) GetBlockByHeight(height int64) (HiveBlock, error) {
	blk := HiveBlock{}

	result := blks.FindOne(context.Background(), bson.M{
		"height": height,
	})

	if result.Err() == nil {
		err := result.Decode(&blk)
		if err != nil {
			return HiveBlock{}, err
		}
		return blk, nil
	} else {
		return HiveBlock{}, result.Err()
	}
}

func (blks *Blocks) InsertBlock(blockData HiveBlock) {
	ctx := context.Background()
	options := options.FindOneAndUpdate().SetUpsert(true)
	blks.FindOneAndUpdate(ctx, bson.M{
		"height": blockData.Height,
	}, bson.M{
		"$set": blockData,
	}, options)
}

package blockdb

import (
	"context"
	"log/slog"
	"mimic/lib/utils"
	"mimic/modules/db"
	"time"

	"github.com/chebyrash/promise"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BlockQuery interface {
	QueryBlockByBlockNum(*HiveBlock, int64) error
	QueryBlockByRange(blocks *[]HiveBlock, start, end int) error
	QueryHeadBlock(context.Context, *HiveBlock) error
}

type BlockCollection struct {
	*mongo.Collection
}

var collection BlockQuery = nil

func New(d *mongo.Database) *BlockCollection {
	collection = &BlockCollection{
		db.NewCollection(d, "blocks"),
	}

	return collection.(*BlockCollection)
}

func Collection() *BlockCollection {
	return collection.(*BlockCollection)
}

// Blocks implement `aggregate.Plugin`
func (d *BlockCollection) Init() error {
	indexName, err := d.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    bson.D{{Key: "block_id", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("block_id_unique"),
	})

	if err != nil {
		return err
	}

	slog.Debug("Index created.", "collection", d.Name(), "index", indexName)
	return nil
}

func (d *BlockCollection) Start() *promise.Promise[any] {
	var blocks []HiveBlock
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	db.Seed(
		&blocks,
		ctx,
		collection.(*BlockCollection).Collection,
		"block_api.get_block.json",
	)
	return utils.PromiseResolve[any](d)
}

func (d *BlockCollection) Stop() error {
	return nil
}

// Queries

func (blks *BlockCollection) GetBlockRange(
	startHeight int64,
	endHeight int64,
) []HiveBlock {
	return nil
}

func (blks *BlockCollection) GetBlockById(id string) HiveBlock {
	return HiveBlock{}
}

func (blks *BlockCollection) GetBlockByHeight(height int64) (HiveBlock, error) {
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

func (blks *BlockCollection) InsertBlock(blockData *HiveBlock) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := blks.InsertOne(ctx, blockData)
	return err
}

package blockdb

import (
	"context"
	"log/slog"
	"mimic/lib/utils"
	"mimic/modules/db"
	"mimic/modules/db/mimic"
	"time"

	"github.com/chebyrash/promise"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BlockQuery interface {
	QueryBlockByBlockNum(*HiveBlock, int64) error
	QueryBlockByRange(blocks *[]HiveBlock, start, end int) error
}

type blockCollection struct {
	*mongo.Collection
}

var collection BlockQuery = nil

func New(d *mimic.MimicDb) *blockCollection {
	collection = &blockCollection{
		db.NewCollection(d.DbInstance, "blocks"),
	}

	return collection.(*blockCollection)
}

func Collection() *blockCollection {
	return collection.(*blockCollection)
}

// Blocks implement `aggregate.Plugin`
func (d *blockCollection) Init() error {
	indexName, err := d.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    bson.D{{Key: "block_id", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("block_id_unique"),
	})

	if err != nil {
		slog.Error(
			"Failed to create index.",
			"collection",
			d.Name(),
			"err",
			err,
		)
		return err
	}

	slog.Debug("Index created.", "collection", d.Name(), "index", indexName)
	return nil
}

func (d *blockCollection) Start() *promise.Promise[any] {
	var blocks []HiveBlock
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	db.Seed(
		&blocks,
		ctx,
		collection.(*blockCollection).Collection,
		"block_api.get_block.json",
	)
	return utils.PromiseResolve[any](d)
}

func (d *blockCollection) Stop() error {
	return nil
}

// Queries

func (blks *blockCollection) GetBlockRange(
	startHeight int64,
	endHeight int64,
) []HiveBlock {
	return nil
}

func (blks *blockCollection) GetBlockById(id string) HiveBlock {
	return HiveBlock{}
}

func (blks *blockCollection) GetBlockByHeight(height int64) (HiveBlock, error) {
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

func (blks *blockCollection) InsertBlock(blockData *HiveBlock) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := blks.InsertOne(ctx, blockData)
	return err
}

func (b *blockCollection) FindLatestBlock(
	ctx context.Context,
	buf *HiveBlock,
) error {
	// since timestamp is encoded with mongodb, can query for lastest inserted ID
	queryOpts := options.FindOne()
	queryOpts.SetSort(bson.M{"_id": -1})

	result := b.Collection.FindOne(ctx, bson.D{}, queryOpts)
	if result.Err() != nil {
		return result.Err()
	}

	return result.Decode(&buf)
}

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

type Blocks struct {
	*mongo.Collection
}

var blockCollection = &Blocks{}

func New(d *mimic.MimicDb) *Blocks {
	blockCollection.Collection = db.NewCollection(d.DbInstance, "blocks")
	return blockCollection
}

func Collection() *Blocks {
	return blockCollection
}

// Blocks implement `aggregate.Plugin`
func (d *Blocks) Init() error {
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

func (d *Blocks) Start() *promise.Promise[any] {
	var blocks []Block
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	db.Seed(
		&blocks,
		ctx,
		blockCollection.Collection,
		"block_api.get_block.json",
	)
	return utils.PromiseResolve[any](d)
}

func (d *Blocks) Stop() error {
	return nil
}

// Queries

func (b *Blocks) QueryBlockByBlockNum(blockBuf *Block, blockNum int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	filter := bson.M{"block_num": bson.M{"$eq": blockNum}}

	result := b.FindOne(ctx, filter)
	if result.Err() != nil {
		return result.Err()
	}

	return result.Decode(blockBuf)
}

func (b *Blocks) QueryBlockByRange(blocks *[]Block, start, end int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	filter := bson.M{"block_num": bson.M{"$gte": start, "$lte": end}}

	cursor, err := b.Find(ctx, filter)
	if err != nil {
		return err
	}

	defer cursor.Close(ctx)

	return cursor.All(ctx, blocks)
}

func (blks *Blocks) GetBlockRange(
	startHeight int64,
	endHeight int64,
) []HiveBlock {
	return nil
}

func (blks *Blocks) GetBlockById(id string) Block {
	return Block{}
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
	/*
		ctx := context.Background()
		options := options.FindOneAndUpdate().SetUpsert(true)
		blks.FindOneAndUpdate(ctx, bson.M{
			"height": blockData.Height,
		}, bson.M{
			"$set": blockData,
		}, options)
	*/
}

func (b *Blocks) FindLatestBlock(ctx context.Context, buf *HiveBlock) error {
	// since timestamp is encoded with mongodb, can query for lastest inserted ID
	queryOpts := options.FindOne()
	queryOpts.SetSort(bson.M{"_id": -1})

	result := blockCollection.FindOne(ctx, bson.D{}, queryOpts)
	if result.Err() != nil {
		return result.Err()
	}

	return result.Decode(&buf)
}

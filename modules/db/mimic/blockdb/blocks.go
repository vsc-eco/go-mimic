package blockdb

import (
	"context"
	"mimic/modules/db"
	"mimic/modules/db/mimic"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Blocks struct {
	*db.Collection
}

func New(d *mimic.MimicDb) Blocks {
	return Blocks{db.NewCollection(d.DbInstance, "blocks")}
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

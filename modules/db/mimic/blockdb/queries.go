package blockdb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func (b *blockCollection) QueryBlockByBlockNum(
	blockBuf *HiveBlock,
	blockNum int64,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	filter := bson.M{"block_num": bson.M{"$eq": blockNum}}

	result := b.FindOne(ctx, filter)
	if result.Err() != nil {
		return result.Err()
	}

	return result.Decode(blockBuf)
}

func (b *blockCollection) QueryBlockByRange(
	blocks *[]HiveBlock,
	start, end int,
) error {
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

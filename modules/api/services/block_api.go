package services

import (
	"mimic/modules/db/mimic/blockdb"

	"golang.org/x/exp/slog"
)

type BlockAPI struct {
}

type GetBlockRangeArgs struct {
	StartingBlockNum int `json:"starting_block_num"`
	Count            int `json:"count"`
}

type GetBlockRangeReply struct {
	Blocks []blockdb.HiveBlock `json:"blocks"`
}

type GetBlockArgs struct {
	BlockNum int64 `json:"block_num"`
}

type GetBlockReply struct {
	Block blockdb.HiveBlock `json:"block"`
}

func (BlockAPI) GetBlock(args *GetBlockArgs, reply *GetBlockReply) {
	blockCollection := blockdb.Collection()

	if err := blockCollection.QueryBlockByBlockNum(&reply.Block, args.BlockNum); err != nil {
		slog.Error("Failed to query block by block number.",
			"block-num", args.BlockNum, "err", err)
	}
}

func (BlockAPI) GetBlockRange(
	args *GetBlockRangeArgs,
	reply *GetBlockRangeReply,
) {
	start := args.StartingBlockNum
	end := start + args.Count

	blockCollection := blockdb.Collection()
	if err := blockCollection.QueryBlockByRange(&reply.Blocks, start, end); err != nil {
		slog.Error("Failed to query block by block range.",
			"start", start, "end", end, "err", err)
	}
}

func (BlockAPI) Expose(rm RegisterMethod) {
	rm("get_block_range", "GetBlockRange")
	rm("get_block", "GetBlock")
}

package services

import (
	"mimic/modules/db/mimic/blockdb"

	"github.com/sourcegraph/jsonrpc2"
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

func (*BlockAPI) GetBlock(
	args *GetBlockArgs,
) (*GetBlockReply, *jsonrpc2.Error) {
	reply := &GetBlockReply{}
	err := blockdb.Collection().QueryBlockByBlockNum(
		&reply.Block, args.BlockNum,
	)
	if err != nil {
		slog.Error("Failed to query block by block number.",
			"block-num", args.BlockNum, "err", err)

		return nil, &jsonrpc2.Error{
			Code:    jsonrpc2.CodeInvalidRequest,
			Message: "block not found",
		}
	}

	return reply, nil
}

func (BlockAPI) GetBlockRange(
	args *GetBlockRangeArgs,
) (*GetBlockRangeReply, *jsonrpc2.Error) {
	reply := &GetBlockRangeReply{}
	start := args.StartingBlockNum
	end := start + args.Count

	blockCollection := blockdb.Collection()

	err := blockCollection.QueryBlockByRange(&reply.Blocks, start, end)
	if err != nil {
		slog.Error("Failed to query block by block range.",
			"start", start, "end", end, "err", err)
		return nil, &jsonrpc2.Error{
			Code:    jsonrpc2.CodeInvalidRequest,
			Message: "blocks not found",
		}
	}

	return reply, nil
}

func (BlockAPI) Expose(rm RegisterMethod) {
	rm("get_block_range", "GetBlockRange")
	rm("get_block", "GetBlock")
}

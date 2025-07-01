package services

import (
	"fmt"
	"mimic/mock"
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
	Blocks []blockdb.Block `json:"blocks"`
}

func (BlockAPI) GetBlockRange(args *GetBlockRangeArgs, reply *GetBlockRangeReply) {
	data, err := mock.GetMockData[GetBlockReply]("mockdata/block_api.get_block.mock.json")
	if err != nil {
		panic(err)
	}

	for i := 0; i <= args.Count; i++ {
		blockIndex := fmt.Sprintf("%d", args.StartingBlockNum+i)
		block, ok := data[blockIndex]
		if !ok {
			continue
		}
		reply.Blocks = append(reply.Blocks, block.Block)
	}
}

type GetBlockArgs struct {
	BlockNum int64 `json:"block_num"`
}

type GetBlockReply struct {
	Block blockdb.Block `json:"block"`
}

func (BlockAPI) GetBlock(args *GetBlockArgs, reply *GetBlockReply) {
	data, err := mock.GetMockData[GetBlockReply]("mockdata/block_api.get_block.mock.json")
	if err != nil {
		panic(err)
	}

	block, ok := data[fmt.Sprintf("%d", args.BlockNum)]
	if !ok {
		slog.Error("Block not found.", "block_num", args.BlockNum)
	} else {
		*reply = block
	}
}

func (BlockAPI) Expose(rm RegisterMethod) {
	rm("get_block_range", "GetBlockRange")
	rm("get_block", "GetBlock")
}

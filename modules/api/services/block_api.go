package services

import (
	"fmt"

	"golang.org/x/exp/slog"
)

type BlockAPI struct {
}

type GetBlockRangeArgs struct {
}

type GetBlockRangeReply struct {
}

func (BlockAPI) GetBlockRange(args *GetBlockRangeArgs, reply *GetBlockRangeReply) {

}

type GetBlockArgs struct {
	BlockNum int64 `json:"block_num"`
}

type getBlockBlock struct {
	Previous              string   `json:"previous"`
	Timestamp             string   `json:"timestamp"`
	Witness               string   `json:"witness"`
	TransactionMerkleRoot string   `json:"transaction_merkle_root"`
	Extensions            []string `json:"extensions"`
	WitnessSignature      string   `json:"witness_signature"`
	Transactions          []string `json:"transactions"`
	BlockId               string   `json:"block_id"`
	SigningKey            string   `json:"signing_key"`
	TransactionIds        []string `json:"transaction_ids"`
}
type GetBlockReply struct {
	Block getBlockBlock `json:"block"`
}

func (BlockAPI) GetBlock(args *GetBlockArgs, reply *GetBlockReply) {
	data, err := getMockData[GetBlockReply]("mockdata/block_api.get_block.mock.json")
	if err != nil {
		panic(err)
	}
	fmt.Println(args)

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

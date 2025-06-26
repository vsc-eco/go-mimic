package services

type AccountHistoryApi struct {
}

type GetOpsInBlockArgs struct {
	StartingBlockRange uint64 `json:"starting_block_range"`
	Count              uint64 `json:"count"`
}

type GetOpsInBlockReply struct {
	Blocks []string `json:"blocks"`
}

func (ahapi AccountHistoryApi) GetOpsInBlock(args *GetOpsInBlockArgs, reply *GetOpsInBlockReply) {

}

func (t *AccountHistoryApi) Expose(rm RegisterMethod) {
	rm("get_ops_in_block", "GetOpsInBlock")
}

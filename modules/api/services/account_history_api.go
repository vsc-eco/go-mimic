package services

import (
	"mimic/mock"

	"golang.org/x/exp/slog"
)

type AccountHistoryApi struct {
}

type GetOpsInBlockArgs struct {
	BlockNum    uint64 `json:"block_num"`
	OnlyVirtual bool   `json:"only_virtual"`
	// ignoring the include_reversible value
}

type GetOpsInBlockReply struct {
	Ops []Operation `json:"ops"`
}

type Operation struct {
	TrxID      string `json:"trx_id"`
	Block      int64  `json:"block"`
	TrxInBlock int64  `json:"trx_in_block"`
	OpInTrx    int64  `json:"op_in_trx"`
	VirtualOp  int64  `json:"virtual_op"`
	Timestamp  string `json:"timestamp"`
	Op         any    `json:"op"`
}

func (a *AccountHistoryApi) GetOpsInBlock(args *GetOpsInBlockArgs, reply *GetOpsInBlockReply) {
	mockData := make(map[string][]Operation)

	err := mock.GetMockData(&mockData, "account_history_api.get_ops_in_block")
	if err != nil {
		slog.Error("Failed to serve mock data.", "err", err)
		return
	}

	var key string
	if !args.OnlyVirtual {
		key = "virtual-ops-true"
	} else {
		key = "virtual-ops-false"
	}

	reply.Ops = mockData[key]
}

func (t *AccountHistoryApi) Expose(rm RegisterMethod) {
	rm("get_ops_in_block", "GetOpsInBlock")
}

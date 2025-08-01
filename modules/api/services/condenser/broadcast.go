package condenser

import (
	"mimic/modules/producers"

	"github.com/sourcegraph/jsonrpc2"
)

func (c *Condenser) BroadcastTransaction(
	args *CondenserParam,
) (map[string]any, *jsonrpc2.Error) {
	go c.BroadcastTransactionSynchronous(args)
	return make(map[string]any), nil
}

func (c *Condenser) BroadcastTransactionSynchronous(
	args *CondenserParam,
) (*producers.BroadcastTransactionResponse, *jsonrpc2.Error) {
	trx := args.Trx
	if err := producers.ValidateTransaction(trx); err != nil {
		c.Logger.Error("failed to validate transaction", "err", err)
		return nil, &jsonrpc2.Error{
			Code:    jsonrpc2.CodeInvalidParams,
			Message: "invalid transaction",
		}
	}

	res := producers.BroadcastTransactions(trx)

	reply := res.Response()
	return &reply, nil
}

package condenser

import (
	"mimic/modules/producers"
)

func (c *Condenser) BroadcastTransaction(
	args *CondenserParam,
	reply *map[string]any,
) {
	go c.BroadcastTransactionSynchronous(
		args,
		&producers.BroadcastTransactionResponse{},
	)
	*reply = make(map[string]any)
}

func (c *Condenser) BroadcastTransactionSynchronous(
	args *CondenserParam,
	reply *producers.BroadcastTransactionResponse,
) {
	trx := args.Trx
	if err := producers.ValidateTransaction(trx); err != nil {
		c.Logger.Error("failed to validate transaction", "err", err)
		return
	}
}

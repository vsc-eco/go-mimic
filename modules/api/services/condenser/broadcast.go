package condenser

import (
	"encoding/json"
	"fmt"
	"mimic/modules/producers"
)

// broadcast_transaction
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

// broadcast_transaction_synchronous
func (c *Condenser) BroadcastTransactionSynchronous(
	args *CondenserParam,
	reply *producers.BroadcastTransactionResponse,
) {
	jj, err := json.MarshalIndent(args, "", "  ")
	fmt.Println(string(jj), err)
	// req := producers.BroadcastTransactions(*args)
	//*reply = req.Response()
}

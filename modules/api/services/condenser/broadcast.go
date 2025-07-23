package condenser

import (
	"encoding/json"
	"fmt"
	"log"
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
		log.Println(err) // TODO: log with module's slog
		return
	}
	jj, err := json.MarshalIndent(args, "", "  ")
	fmt.Println(string(jj), err)
	// req := producers.BroadcastTransactions(*args)
	//*reply = req.Response()
}

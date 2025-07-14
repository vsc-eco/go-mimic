package broadcastops

import (
	"fmt"
	"mimic/modules/api/services"
	"mimic/modules/db/mimic/accountdb"

	"github.com/chebyrash/promise"
)

type BroadcastOpsAccount struct {
	db accountdb.AccountDBQueries
}

// Runs initialization in order of how they are passed in to `Aggregate`
func (b *BroadcastOpsAccount) Init() error {
	b.db = accountdb.Collection()
	return nil
}

// Runs startup and should be non blocking
func (b *BroadcastOpsAccount) Start() *promise.Promise[any] {
	return nil
}

// Runs cleanup once the `Aggregate` is finished
func (b *BroadcastOpsAccount) Stop() error {
	return nil
}

func (b *BroadcastOpsAccount) AccountCreate(
	arg *BroadcastParam[AccountCreateParam],
	reply *any,
) {
	fmt.Println(arg)
}

func (b *BroadcastOpsAccount) Expose(rm services.RegisterMethod) {
	rm("account_create", "AccountCreate")
}

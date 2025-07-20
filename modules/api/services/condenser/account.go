package condenser

import (
	"context"
	"log/slog"
	"mimic/lib/utils"
	"mimic/modules/db/mimic/accountdb"
	"time"

	"github.com/chebyrash/promise"
)

// Runs initialization in order of how they are passed in to `Aggregate`
func (c *Condenser) Init() error {
	return nil
}

// Runs startup and should be non blocking
func (c *Condenser) Start() *promise.Promise[any] {
	return nil
}

// Runs cleanup once the `Aggregate` is finished
func (c *Condenser) Stop() error {
	return nil
}

func (c *Condenser) AccountCreate(
	arg *BroadcastParam[AccountCreateParam],
	reply *any,
) {
	timeStamp := time.Now().Format(utils.TimeFormat)
	a := arg.Param

	account := accountdb.Account{
		Id:   0,
		Name: a.NewAccountName,
		KeySet: accountdb.UserKeySet{
			Owner:   a.Owner,
			Active:  a.Active,
			Posting: a.Posting,
		},
		MemoKey:             a.MemoKey,
		JsonMeta:            a.JsonMetadata,
		LastOwnerUpdate:     timeStamp,
		LastAccountUpdate:   timeStamp,
		Created:             timeStamp,
		JsonPostingMetadata: "",
		Balance:             "",
		HbdBalance:          "",
		SavingsHbdBalance:   "",
		VestingShares:       "",
		Reputation:          0,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := c.AccountDB.InsertAccount(ctx, &account); err != nil {
		slog.Error("failed to create account", "error", err)
	}
}

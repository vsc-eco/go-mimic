package condenser

import (
	"context"
	"fmt"
	"log/slog"
	"mimic/lib/utils"
	"mimic/modules/db/mimic/accountdb"
	"time"

	"github.com/vsc-eco/hivego"
)

func (c *Condenser) AccountUpdate(
	arg *CondenserParam[hivego.AccountUpdateOperation],
	reply *any,
) {
	fmt.Println(arg, reply)
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
			Owner:   &a.Owner,
			Active:  &a.Active,
			Posting: &a.Posting,
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

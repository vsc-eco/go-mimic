package services

import (
	"context"
	"mimic/modules/db/mimic/rcdb"
	"time"

	"golang.org/x/exp/slog"
)

type RcAccount = rcdb.RCAccount

type RcApi struct{}

type rcArgs []string

type RcReply struct {
	RcAccounts []RcAccount `json:"rc_accounts"`
}

func (api *RcApi) FindRcAccounts(args *rcArgs, reply *RcReply) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	reply.RcAccounts = make([]RcAccount, 0)

	err := rcdb.Collection().
		QueryFindRcAccounts(ctx, &reply.RcAccounts, *args)

	if err != nil {
		slog.Error("Failed to queries for Rc Account.",
			"queries", args, "err", err)
	}
}

func (api *RcApi) Expose(mr RegisterMethod) {
	mr("find_rc_accounts", "FindRcAccounts")
}

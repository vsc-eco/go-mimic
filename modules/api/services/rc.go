package services

import "github.com/sourcegraph/jsonrpc2"

type RcApi struct {
}

type rcArgs struct {
	Accounts []string `json:"accounts"`
}

type RcReply struct {
	RcAccounts []RcAccount `json:"rc_accounts"`
}

type RcAccount struct {
	Account                 string `json:"account"`
	DelegatedRc             int    `json:"delegated_rc"`
	MaxRc                   int    `json:"max_rc"`
	MaxRcCreationAdjustment string `json:"max_rc_creation_adjustment"`
	RcManabar               struct {
		CurrentMana    int `json:"current_mana"`
		LastUpdateTime int `json:"last_update_time"`
	} `json:"rc_manabar"`
}

func (api RcApi) FindRcAccounts(args *rcArgs) (*RcReply, *jsonrpc2.Error) {
	reply := &RcReply{}
	for _, account := range args.Accounts {
		reply.RcAccounts = append(reply.RcAccounts, RcAccount{
			Account:                 account,
			DelegatedRc:             0,
			MaxRc:                   1000000000,
			MaxRcCreationAdjustment: "1000000000 VESTS",
			RcManabar: struct {
				CurrentMana    int "json:\"current_mana\""
				LastUpdateTime int "json:\"last_update_time\""
			}{
				CurrentMana:    1000000000,
				LastUpdateTime: 1550731380,
			},
		})
	}

	return reply, nil
}

func (api RcApi) Expose(mr RegisterMethod) {
	mr("find_rc_accounts", "FindRcAccounts")
}

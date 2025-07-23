package main

import "github.com/vsc-eco/hivego"

type accountUpdate struct{}

func (a *accountUpdate) params() any {
	op := &hivego.AccountUpdateOperation{
		Account:      "foo",
		Owner:        nil,
		Active:       nil,
		Posting:      nil,
		MemoKey:      "",
		JsonMetadata: "{}",
	}

	return []any{"account_update", op}
}

func (a *accountUpdate) OpName() string { return "account_update" }

func (a *accountUpdate) SerializeOp() ([]byte, error) {
	panic("not implemented")
}

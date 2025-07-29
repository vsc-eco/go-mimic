package main

import "github.com/vsc-eco/hivego"

type accountUpdate struct{}

func (a *accountUpdate) params() hivego.HiveOperation {
	return &hivego.AccountUpdateOperation{
		Account:      "foo",
		Owner:        nil,
		Active:       nil,
		Posting:      nil,
		MemoKey:      "",
		JsonMetadata: "{}",
	}
}

package main

import (
	"mimic/lib/hive"

	"github.com/vsc-eco/hivego"
)

type accountUpdate struct{}

func (a *accountUpdate) params() hivego.HiveOperation {
	keyset := hive.MakeHiveKeySet("foo", "bar")

	return &hivego.AccountUpdateOperation{
		Account:      "foo",
		MemoKey:      *keyset.MemoKey().GetPublicKeyString(),
		JsonMetadata: "{}",
	}
}

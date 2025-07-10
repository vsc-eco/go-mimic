package services

import (
	"fmt"
	"mimic/lib/encoder"
	"mimic/modules/db/mimic/accountdb"

	"github.com/chebyrash/promise"
)

type BroadcastOps struct {
	db accountdb.AccountDBQueries
}

// Runs initialization in order of how they are passed in to `Aggregate`
func (b *BroadcastOps) Init() error {
	b.db = accountdb.Collection()
	return nil
}

// Runs startup and should be non blocking
func (b *BroadcastOps) Start() *promise.Promise[any] {
	return nil
}

// Runs cleanup once the `Aggregate` is finished
func (b *BroadcastOps) Stop() error {
	return nil
}

type paramGeneric[T any] struct {
	Action string
	Param  T
}

func (p *paramGeneric[T]) UnmarshalJSON(v []byte) error {
	return encoder.JsonArrayDeserialize(p, v)
}

type AccountCreateParam struct {
	Fee struct {
		Amount    string `json:"amount,omitempty"`
		Precision int    `json:"precision,omitempty"`
		Nai       string `json:"nai,omitempty"`
	} `json:"fee"`

	Creator        string                     `json:"creator"`
	NewAccountName string                     `json:"new_account_name"`
	Owner          accountdb.AccountAuthority `json:"owner"`
	Active         accountdb.AccountAuthority `json:"active"`
	Posting        accountdb.AccountAuthority `json:"posting"`
	MemoKey        string                     `json:"memo_key"`
	JsonMetadata   string                     `json:"json_metadata"`
}

func (b *BroadcastOps) AccountCreate(
	arg *paramGeneric[AccountCreateParam],
	reply *any,
) {
	fmt.Println(arg)
}

func (b *BroadcastOps) Expose(rm RegisterMethod) {
	rm("account_create", "AccountCreate")
}

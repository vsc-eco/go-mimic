package condenser

import (
	"encoding/json"
	"mimic/lib/encoder"
	"mimic/modules/db/mimic/accountdb"
)

type BroadcastParam[T any] struct {
	Action string
	Param  T
}

func (p *BroadcastParam[T]) UnmarshalJSON(v []byte) error {
	return encoder.JsonArrayDeserialize(p, v)
}

func (p *BroadcastParam[T]) MarshalJSON() ([]byte, error) {
	buf := [2]any{p.Action, p.Param}
	return json.Marshal(&buf)
}

type AccountCreateParam struct {
	Fee            AccountCreateFee           `json:"fee"`
	Creator        string                     `json:"creator"`
	NewAccountName string                     `json:"new_account_name"`
	Owner          accountdb.AccountAuthority `json:"owner"`
	Active         accountdb.AccountAuthority `json:"active"`
	Posting        accountdb.AccountAuthority `json:"posting"`
	MemoKey        string                     `json:"memo_key"`
	JsonMetadata   string                     `json:"json_metadata"`
}

type AccountCreateFee struct {
	Amount    string `json:"amount,omitempty"`
	Precision int    `json:"precision,omitempty"`
	Nai       string `json:"nai,omitempty"`
}

package condenser

import (
	"encoding/json"
	"mimic/lib/encoder"
	"mimic/lib/validator"

	"github.com/vsc-eco/hivego"
)

type CondenserParam[T hivego.HiveOperation] struct {
	Trx *Transaction[T] `json:"trx,omitempty" validate:"required"`
}

type Transaction[T hivego.HiveOperation] struct {
	Expiration           string   `json:"expiration"`
	Extensions           []any    `json:"extensions"`
	Operations           []T      `json:"operations"`
	RefBlockNum          uint16   `json:"ref_block_num"`
	RefBlockPrefix       uint32   `json:"ref_block_prefix"`
	Signatures           []string `json:"signatures"`
	RequiredAuths        []string `json:"required_auths,omitempty"`
	RequiredPostingAuths []string `json:"required_posting_auths,omitempty"`
}

func (p *CondenserParam[T]) UnmarshalJSON(data []byte) error {
	// the call to json.Unmarshal will invoke the function
	// `json.Unmarshaler.UnmarshalJSON`, resulting in a infinite recursion.
	// to avoid the recusrive call to `UnmarshalJSON`, alias the type
	type Alias CondenserParam[T]
	aux := (*Alias)(p)
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	return validator.New().Struct(p)
}

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
	Fee            AccountCreateFee `json:"fee"`
	Creator        string           `json:"creator"`
	NewAccountName string           `json:"new_account_name"`
	Owner          hivego.Auths     `json:"owner"`
	Active         hivego.Auths     `json:"active"`
	Posting        hivego.Auths     `json:"posting"`
	MemoKey        string           `json:"memo_key"`
	JsonMetadata   string           `json:"json_metadata"`
}

type AccountCreateFee struct {
	Amount    string `json:"amount,omitempty"`
	Precision int    `json:"precision,omitempty"`
	Nai       string `json:"nai,omitempty"`
}

package condenser

import (
	"encoding/json"
	"fmt"
	"mimic/lib/encoder"
	"mimic/lib/validator"

	"github.com/vsc-eco/hivego"
)

type CondenserParam struct {
	Trx *hivego.HiveTransaction `json:"trx,omitempty" validate:"required"`
}

func (p *CondenserParam) UnmarshalJSON(data []byte) error {
	// the call to json.Unmarshal will invoke the function
	// `json.Unmarshaler.UnmarshalJSON()`, resulting in a infinite recursion.
	// to avoid the recusrive call to `UnmarshalJSON`, alias the type
	type Alias CondenserParam
	if err := json.Unmarshal(data, (*Alias)(p)); err != nil {
		return err
	}

	if err := validator.New().Struct(p); err != nil {
		return err
	}

	p.Trx.Operations = make([]hivego.HiveOperation, len(p.Trx.OperationsJs))
	for i, opRaw := range p.Trx.OperationsJs {
		trxBuf, err := parseTrxType(opRaw[0].(string))
		if err != nil {
			return err
		}

		if err := serializeToStruct(trxBuf, opRaw[1]); err != nil {
			return err
		}

		p.Trx.Operations[i] = trxBuf
	}

	return nil
}

func parseTrxType(opName string) (hivego.HiveOperation, error) {
	switch opName {
	case "custom_json":
		return &hivego.CustomJsonOperation{}, nil

	default:
		return nil, fmt.Errorf(
			"transaction typed %s does not implement hivego.HiveOperation interface",
			opName,
		)
	}
}

func serializeToStruct(buf hivego.HiveOperation, data any) error {
	j, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(j, buf); err != nil {
		return err
	}

	return validator.New().Struct(buf)
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

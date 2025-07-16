package transactions

import (
	"encoding/json"
	"mimic/modules/db/mimic/accountdb"
)

type AccountCreateOp struct {
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

// AccountCreateTRX implements hivego.HiveOperation
func (a *AccountCreateOp) SerializeOp() ([]byte, error) {
	return json.Marshal(a)
}

// AccountCreateTRX implements hivego.HiveOperation
func (a *AccountCreateOp) OpName() string {
	return "account_create"
}

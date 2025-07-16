package transactions

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"mimic/lib/utils"
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
	Precision uint8  `json:"precision,omitempty"` // TODO: why is this uint8?
	Nai       string `json:"nai,omitempty"`
}

// AccountCreateTRX implements hivego.HiveOperation
func (a *AccountCreateOp) SerializeOp() ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	// account create fee
	if err := binary.Write(buf, utils.HiveBinaryEndianess, a.Fee.Amount); err != nil {
		return nil, err
	}

	if err := buf.WriteByte(byte(a.Fee.Precision)); err != nil {
		return nil, err
	}

	if err := writeString(buf, a.Fee.Nai); err != nil {
		return nil, err
	}

	// creator, account name
	if err := writeString(buf, a.Creator); err != nil {
		return nil, err
	}

	if err := writeString(buf, a.NewAccountName); err != nil {
		return nil, err
	}

	// authorities
	if err := writeAuthority(buf, &a.Owner); err != nil {
		return nil, err
	}

	if err := writeAuthority(buf, &a.Active); err != nil {
		return nil, err
	}

	if err := writeAuthority(buf, &a.Posting); err != nil {
		return nil, err
	}

	// memo key + json metadata
	if err := writeString(buf, a.MemoKey); err != nil {
		return nil, err
	}

	if err := writeString(buf, a.JsonMetadata); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// AccountCreateTRX implements hivego.HiveOperation
func (a *AccountCreateOp) OpName() string {
	return "account_create"
}

func writeVarInt(buf *bytes.Buffer, value uint64) error {
	var b [8]byte
	n := binary.PutUvarint(b[:], value)
	_, err := buf.Write(b[:n])
	return err
}

func writeString(buf *bytes.Buffer, value string) error {
	err := writeVarInt(buf, uint64(len(value)))
	if err != nil {
		return err
	}
	_, err = buf.WriteString(value)
	return err
}

func writeAuthority(
	_ *bytes.Buffer,
	_ *accountdb.AccountAuthority,
) error {
	fmt.Println("TODO: implement writeAuthority")
	return nil
}

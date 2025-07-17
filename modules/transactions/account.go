package transactions

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"mimic/lib/utils"
	"sort"

	"github.com/vsc-eco/hivego"
)

type AccountCreateOp struct {
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
	Precision uint8  `json:"precision,omitempty"` // TODO: why is this uint8?
	Nai       string `json:"nai,omitempty"`
}

// AccountCreateTRX implements hivego.HiveOperation
func (a *AccountCreateOp) SerializeOp() ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	// account create fee
	err := binary.Write(buf, utils.HiveBinaryEndianess, a.Fee.Amount)
	if err != nil {
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

// ported from: https://github.com/vsc-eco/hivego/blob/fa6c9e2c8be757b260a9b48b7d206fa02f8cfde9/serializer.go#L325C1-L332C2
// TODO: just make this public at hivego?
func writeAuthority(buf *bytes.Buffer, auth *hivego.Auths) error {
	if auth == nil {
		return buf.WriteByte(0)
	}

	if err := buf.WriteByte(1); err != nil {
		return err
	}

	// write weight_threshold
	err := binary.Write(buf, binary.LittleEndian, uint32(auth.WeightThreshold))
	if err != nil {
		return fmt.Errorf("Error writing weight_threshold: %v", err)
	}

	// write account_auths
	err = hivego.WriteUvarint(buf, uint64(len(auth.AccountAuths)))
	if err != nil {
		return fmt.Errorf("error writing account_auths length: %v", err)
	}

	for _, accountAuth := range auth.AccountAuths {
		writeString(buf, accountAuth[0].(string))
		err = binary.Write(
			buf,
			utils.HiveBinaryEndianess,
			uint16(accountAuth[1].(uint16)),
		)
		if err != nil {
			return fmt.Errorf("error writing account_auth weight: %v", err)
		}
	}

	// write key_auths
	err = hivego.WriteUvarint(buf, uint64(len(auth.KeyAuths)))
	if err != nil {
		return fmt.Errorf("error writing key_auths length: %v", err)
	}

	// sorting pub keys by value?
	sort.SliceStable(auth.KeyAuths, func(i, j int) bool {
		return auth.KeyAuths[i][0].(string) < auth.KeyAuths[j][0].(string)
	})

	// serialize pub keys
	for _, keyAuth := range auth.KeyAuths {
		pk, err := hivego.DecodePublicKey(keyAuth[0].(string))
		if err != nil {
			return fmt.Errorf("error decoding public key: %v", err)
		}

		if err := binary.Write(buf, utils.HiveBinaryEndianess, pk.SerializeCompressed()); err != nil {
			return fmt.Errorf("error writing public key: %v", err)
		}

		err = binary.Write(buf, binary.LittleEndian, uint16(keyAuth[1].(int)))
		if err != nil {
			return fmt.Errorf("error writing key_auth weight: %v", err)
		}
	}

	return nil
}

func writeString(buf *bytes.Buffer, value string) error {
	err := hivego.WriteUvarint(buf, uint64(len(value)))
	if err != nil {
		return err
	}
	_, err = buf.WriteString(value)
	return err
}

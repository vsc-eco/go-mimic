package transactions

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"mimic/modules/crypto"
	"mimic/modules/db/mimic/condenserdb"
	"time"
)

type Transaction struct {
	RefBlockNum    uint32      `json:"ref_block_num"`
	RefBlockPrefix uint32      `json:"ref_block_prefix"`
	Expiration     string      `json:"expiration"`
	Operations     []Operation `json:"operations"`
	Extensions     []any       `json:"extensions"`
	Signatures     []string    `json:"signatures"`
}

type Operation struct {
	Action  string
	Payload PayloadSerialize
}

func (o *Operation) ActionCode() uint32 {
	switch o.Action {
	case "vote":
		return 0
	default:
		panic("Action not handled.")
	}
}

type PayloadSerialize interface {
	Serialize() []byte
}

type transactionBuilder struct {
	headBlock *condenserdb.GlobalProperties
	exp       *time.Time
	ext       []any
}

func TransactionBuilder(
	block *condenserdb.GlobalProperties,
) transactionBuilder {
	defaultExp := time.Now().Add(time.Hour * 24)
	return transactionBuilder{headBlock: block, exp: &defaultExp, ext: []any{}}
}

func (t *transactionBuilder) SetExpiration(exp *time.Time) *transactionBuilder {
	t.exp = exp
	return t
}

func (t *transactionBuilder) SetExtensions(
	extensions []any,
) *transactionBuilder {
	t.ext = extensions
	return t
}

func (t *transactionBuilder) Sign(
	operations []Operation,
	signingKey *crypto.HiveKey,
) (*Transaction, error) {

	refBlockNum := t.headBlock.HeadBlockNumber

	refBlockPrefix, err := hex.DecodeString(t.headBlock.HeadBlockID)
	if err != nil {
		return nil, err
	}
	refBlockPrefix = refBlockPrefix[:4]

	// write data to buf
	buf := make([]byte, 0, 1024)

	binary.LittleEndian.AppendUint16(buf, uint16(refBlockNum))
	buf = append(buf, refBlockPrefix...)

	binary.LittleEndian.AppendUint32(buf, uint32(t.exp.Unix()))
	binary.LittleEndian.AppendUint32(buf, uint32(len(operations)))

	for _, op := range operations {
		binary.LittleEndian.AppendUint32(buf, op.ActionCode())
		buf = append(buf, op.Payload.Serialize()...)
	}

	binary.LittleEndian.AppendUint32(buf, uint32(len(t.ext)))

	// sign buf
	sig, err := signingKey.Sign(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %v", err)
	}

	trx := &Transaction{
		RefBlockNum:    refBlockNum,
		RefBlockPrefix: binary.LittleEndian.Uint32(refBlockPrefix),
		Expiration:     t.exp.Format(time.RFC3339),
		Operations:     operations,
		Extensions:     []any{},
		Signatures:     []string{hex.EncodeToString(sig)},
	}
	return trx, nil
}

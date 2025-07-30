package producers

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"mimic/lib/hive"
	"mimic/lib/hive/hiveop"
	"mimic/lib/utils"
	"mimic/modules/db/mimic/accountdb"
	"mimic/modules/producers/opvalidator"
	"time"

	"github.com/decred/dcrd/dcrec/secp256k1/v2"
	"github.com/vsc-eco/hivego"
)

var (
	errMissingSignature = errors.New("missing signature")
	errMissingKey       = errors.New("missing key")
)

type (
	BroadcastTransactionResponse struct {
		ID       string `json:"id"`
		BlockNum uint32 `json:"block_num"`
		TrxNum   uint32 `json:"trx_num"`
		Expired  bool   `json:"expired"`
	}

	keyTypeCache    map[string]map[hive.KeyRole]*secp256k1.PublicKey
	pubKeyExtractor func(string) (*secp256k1.PublicKey, error)

	trxRequest struct {
		comm chan BroadcastTransactionResponse
		trx  *hivego.HiveTransaction
	}
)

func BroadcastTransactions(trx *hivego.HiveTransaction) trxRequest {
	req := trxRequest{
		comm: make(chan BroadcastTransactionResponse),
		trx:  trx,
	}
	producer.trxQueue <- req
	return req
}

func (t *trxRequest) Response() BroadcastTransactionResponse {
	return <-t.comm
}

func (t *trxRequest) Close() {
	close(t.comm)
}

func ValidateTransaction(trx *hivego.HiveTransaction) error {
	if len(trx.Signatures) == 0 {
		return errMissingSignature
	}

	// validate operations
	if err := utils.TryForEach(trx.Operations, validateOp); err != nil {
		return err
	}

	// validate signatures

	// serialize transaction
	trxBytes, err := hivego.SerializeTx(*trx)
	if err != nil {
		return err
	}

	trxBytes = hivego.HashTxForSig(trxBytes)
	extractor := makePubkeyExtractor(trxBytes)

	// extracted pub keys from signatures
	signedPks, err := utils.TryMap(trx.Signatures, extractor)
	if err != nil {
		return err
	}

	// get required pub keys
	keyBuf, err := getPubKeys(trx)
	if err != nil {
		return err
	}

	pubKeyBuf := make([]*secp256k1.PublicKey, 0)
	for _, pubKey := range keyBuf {
		for _, key := range pubKey {
			pubKeyBuf = append(pubKeyBuf, key)
		}
	}

	// validate the key exists
	for _, pk := range signedPks {
		if !validKey(pubKeyBuf, pk) {
			return errMissingKey
		}
	}

	return nil
}

func validKey(
	pubKeyBuf []*secp256k1.PublicKey,
	signedPk *secp256k1.PublicKey,
) bool {
	for _, pk := range pubKeyBuf {
		if pk.IsEqual(signedPk) {
			return true
		}
	}

	return false
}

func validateOp(op hivego.HiveOperation) error {
	v, err := opvalidator.NewValidator(op.OpName())
	if err != nil {
		return err
	}

	return v.ValidateOperation(op)
}

func makePubkeyExtractor(txDigest []byte) pubKeyExtractor {
	return func(sigStr string) (*secp256k1.PublicKey, error) {
		sigByte, err := hex.DecodeString(sigStr)
		if err != nil {
			return nil, err
		}

		pubKey, compacted, err := secp256k1.RecoverCompact(
			sigByte,
			txDigest,
		)
		if err != nil {
			return nil, err
		}

		if !compacted {
			return nil, errors.New("expected compacted signatures")
		}

		return pubKey, nil
	}
}

func getPubKeys(transaction *hivego.HiveTransaction) (keyTypeCache, error) {
	keyBuf := make(keyTypeCache)

	for _, opRaw := range transaction.OperationsJs {
		op, err := getOp(opRaw)
		if err != nil {
			panic(err)
		}

		for _, auth := range op.SigningAuthorities() {
			if _, ok := keyBuf[auth.Account]; !ok {
				keyBuf[auth.Account] = make(
					map[hive.KeyRole]*secp256k1.PublicKey,
				)
			}
			keyBuf[auth.Account][auth.KeyType] = nil
		}
	}

	signingAuths := []string{}
	for k := range keyBuf {
		signingAuths = append(signingAuths, k)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := accountdb.Collection().
		QueryPubKeysByAccount(ctx, keyBuf, signingAuths)
	if err != nil {
		return nil, fmt.Errorf("failed to query for public keys: %w", err)
	}

	// if any key is missing (ie, an empty string), then return an error
	for _, account := range keyBuf {
		for _, keyType := range account {
			if keyType == nil {
				return nil, errMissingKey
			}
		}
	}

	return keyBuf, nil
}

func getOp(operation [2]any) (hiveop.Operation, error) {
	var (
		op  hiveop.Operation = nil
		err error            = nil
	)

	opName, ok := operation[0].(string)
	if !ok {
		return nil, opvalidator.ErrInvalidOperation
	}

	switch opName {
	case "custom_json":
		op = &hiveop.CustomJson{}

	default:
		err = fmt.Errorf("unknown operation: %s", opName)
	}

	if err != nil {
		return nil, err
	}

	jBytes, err := json.Marshal(operation[1])
	if err != nil {
		return nil, err
	}

	return op, json.Unmarshal(jBytes, op)
}

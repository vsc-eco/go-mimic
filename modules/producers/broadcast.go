package producers

import (
	"context"
	"encoding/hex"
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

type keyTypeCache = map[string]map[hive.KeyRole]*secp256k1.PublicKey

func ValidateTransaction(transaction *hivego.HiveTransaction) error {
	if len(transaction.Signatures) == 0 {
		return errMissingSignature
	}

	// validate operations
	for _, op := range transaction.Operations {
		v, err := opvalidator.NewValidator(op.OpName())
		if err != nil {
			if errors.Is(err, opvalidator.ErrUnimplementedValidator) {
				panic(err)
			} else {
				return err
			}
		}

		if err := v.ValidateOperation(op); err != nil {
			return err
		}
	}

	// sereialize transaction
	txBytes, err := hivego.SerializeTx(*transaction)
	if err != nil {
		return err
	}

	// get required pub keys
	keyBuf := make(keyTypeCache)

	for _, opRaw := range transaction.OperationsJs {
		opName, ok := opRaw[0].(string)
		if !ok {
			return fmt.Errorf("invalid operation name: %v", opRaw[0])
		}

		op, err := getOp(opName)
		if err != nil {
			return err
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

	if err := getPubKeys(keyBuf); err != nil {
		return err
	}

	pubKeyBuf := make([]*secp256k1.PublicKey, 0)
	for _, pubKey := range keyBuf {
		for _, key := range pubKey {
			pubKeyBuf = append(pubKeyBuf, key)
		}
	}

	// verify each signature
	sigsBytes, err := utils.TryMap(transaction.Signatures, hex.DecodeString)
	for _, sig := range sigsBytes {
		pubKey, compact, err := secp256k1.RecoverCompact(sig, txBytes)
		if err != nil {
			return err
		}

		if !compact {
			return errors.New("uncompacted key not supported")
		}

		if !pubKeyIncluded(pubKeyBuf, pubKey) {
			return errMissingKey
		}
	}

	return nil
}

func pubKeyIncluded(
	pubKeyBuf []*secp256k1.PublicKey,
	key *secp256k1.PublicKey,
) bool {
	for _, k := range pubKeyBuf {
		if k.IsEqual(key) {
			return true
		}
	}
	return false
}

func getPubKeys(keyBuf keyTypeCache) error {
	signingAuths := []string{}
	for k := range keyBuf {
		signingAuths = append(signingAuths, k)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := accountdb.Collection().
		QueryPubKeysByAccount(ctx, keyBuf, signingAuths)
	if err != nil {
		return fmt.Errorf("failed to query for public keys: %w", err)
	}

	// if any key is missing (ie, an empty string), then return an error
	for _, account := range keyBuf {
		for _, keyType := range account {
			if keyType == nil {
				return errMissingKey
			}
		}
	}

	return nil
}

func getOp(opName string) (hiveop.Operation, error) {
	var (
		op  hiveop.Operation = nil
		err error            = nil
	)

	switch opName {
	case "custom_json":
		op = &hiveop.CustomJson{}

	default:
		err = fmt.Errorf("unknown operation: %s", opName)
	}

	return op, err
}

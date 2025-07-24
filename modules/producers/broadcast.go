package producers

import (
	"context"
	"encoding/hex"
	"errors"
	"mimic/modules/db/mimic/accountdb"
	"time"

	"github.com/decred/dcrd/dcrec/secp256k1/v2"
	"github.com/vsc-eco/hivego"
)

func ValidateTransaction(transaction *hivego.HiveTransaction) error {
	sig, err := hex.DecodeString(transaction.Signatures[0])
	if err != nil {
		return err
	}

	txDigest, err := hivego.SerializeTx(*transaction)
	if err != nil {
		return err
	}

	pubKey, _, err := secp256k1.RecoverCompact(sig, txDigest)
	if err != nil {
		return err
	}

	pubKeyWif := hivego.GetPublicKeyString(pubKey)
	if pubKeyWif == nil {
		return errors.New("bad public key")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	return accountdb.Collection().
		QueryAccountByPubKeyWIF(ctx, &accountdb.Account{}, *pubKeyWif)
}

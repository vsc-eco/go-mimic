package producers

import (
	"errors"
	"fmt"
	"mimic/lib/hive"
	"mimic/lib/hive/hiveop"

	"github.com/vsc-eco/hivego"
)

var errMissingSignature = errors.New("missing signature")

type keyTypeCache = map[string]map[hive.KeyRole]any

func ValidateTransaction(transaction *hivego.HiveTransaction) error {
	if len(transaction.Signatures) == 0 {
		return errMissingSignature
	}

	key := make(keyTypeCache)

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
			if _, ok := key[auth.Account]; !ok {
				key[auth.Account] = make(map[hive.KeyRole]any)
			}
			key[auth.Account][auth.KeyType] = struct{}{}
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

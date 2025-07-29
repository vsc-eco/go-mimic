package hive

import (
	"encoding/json"
	"fmt"

	"github.com/vsc-eco/hivego"
)

// Transaction implements json.Marshaler
type Transaction struct {
	*hivego.HiveTransaction
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	t.OperationsJs = make([][2]interface{}, 0, len(t.Operations))
	for _, op := range t.Operations {
		t.OperationsJs = append(t.OperationsJs, [2]interface{}{op.OpName(), op})
	}

	return json.Marshal(t.HiveTransaction)
}

func (t *Transaction) UnmarshalJSON(data []byte) error {
	trx := new(hivego.HiveTransaction)

	if err := json.Unmarshal(data, trx); err != nil {
		return err
	}

	opCount := len(trx.OperationsJs)
	trx.Operations = make([]hivego.HiveOperation, opCount)

	for i := range opCount {
		opRaw := trx.OperationsJs[i]

		opName := opRaw[0].(string)
		opData := opRaw[1]

		op, err := getHiveOp(opName)
		if err != nil {
			return err
		}

		// re-encode to json bytes, then decode again to struct
		j, err := json.Marshal(opData)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(j, op); err != nil {
			return err
		}
	}

	t.HiveTransaction = trx
	return nil
}

func getHiveOp(opName string) (hivego.HiveOperation, error) {
	var (
		hiveOp hivego.HiveOperation = nil
		err    error                = nil
	)

	switch opName {

	case "custom_json":
		hiveOp = new(hivego.CustomJsonOperation)

	default:
		err = fmt.Errorf("unimplemented operation type: %s", opName)
	}

	return hiveOp, err
}

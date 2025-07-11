package transactions

import (
	"mimic/modules/api/services"
)

type Transaction[Operations []any] struct {
	RefBlockNum    uint32     `json:"ref_block_num"`
	RefBlockPrefix uint32     `json:"ref_block_prefix"`
	Expiration     string     `json:"expiration"`
	Operations     Operations `json:"operations"`
	Extensions     []any      `json:"extensions"`
	Signatures     []string   `json:"signatures"`
}

type transactionBuilder struct {
}

func TransactionBuilder(block *services.GlobalProperties)

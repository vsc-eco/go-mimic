package producers

import "mimic/modules/db/mimic/transactiondb"

type transactionRequest struct {
	comm        chan BroadcastTransactionResponse
	transaction []transactiondb.Transaction
}

func BroadcastTransactions(trx []transactiondb.Transaction) transactionRequest {
	req := transactionRequest{
		comm:        make(chan BroadcastTransactionResponse),
		transaction: trx,
	}
	producer.trxQueue <- req
	return req
}

func (t *transactionRequest) Response() BroadcastTransactionResponse {
	return <-t.comm
}

type BroadcastTransactionResponse struct {
	ID       string `json:"id"`
	BlockNum uint32 `json:"block_num"`
	TrxNum   uint32 `json:"trx_num"`
	Expired  bool   `json:"expired"`
}

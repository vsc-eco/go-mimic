package producers

type transactionRequest struct {
	comm        chan BroadcastTransactionResponse
	transaction any // TODO: update this transaction type
}

func BroadcastTransactions(trx any) transactionRequest {
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

package producers

type transactionRequest struct {
	comm    chan BroadcastTransactionResponse
	payload any
}

func BroadcastTransaction(trx any) transactionRequest {
	req := transactionRequest{
		comm:    make(chan BroadcastTransactionResponse),
		payload: trx,
	}
	producer.trxQueue <- req
	return req
}

func (t *transactionRequest) Response() BroadcastTransactionResponse {
	return <-t.comm
}

type BroadcastTransactionResponse struct {
	ID       string `json:"id"`
	BlockNum int64  `json:"block_num"`
	TrxNum   int64  `json:"trx_num"`
	Expired  bool   `json:"expired"`
}

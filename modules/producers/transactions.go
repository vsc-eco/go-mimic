package producers

import "github.com/vsc-eco/hivego"

type transactionRequest struct {
	comm chan BroadcastTransactionResponse
	trx  *hivego.HiveTransaction
}

func BroadcastTransactions(trx *hivego.HiveTransaction) transactionRequest {
	req := transactionRequest{
		comm: make(chan BroadcastTransactionResponse),
		trx:  trx,
	}
	producer.trxQueue <- req
	return req
}

func (t *transactionRequest) Response() BroadcastTransactionResponse {
	return <-t.comm
}

func (t *transactionRequest) Close() {
	close(t.comm)
}

type BroadcastTransactionResponse struct {
	ID       string `json:"id"`
	BlockNum uint32 `json:"block_num"`
	TrxNum   uint32 `json:"trx_num"`
	Expired  bool   `json:"expired"`
}

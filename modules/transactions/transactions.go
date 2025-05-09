package transactions

import "github.com/vsc-eco/hivego"

type TransactionCore struct {
	TransactionQueue []hivego.HiveTransaction
}

// Fill out
func (t *TransactionCore) PostTransaction() {

}

func (t *TransactionCore) ValidateTransaction() {

}

func (t *TransactionCore) ValidationExpiration() {

}

func (t *TransactionCore) ClearQueue() {
	t.TransactionQueue = []hivego.HiveTransaction{}
}

func New() *TransactionCore {
	return &TransactionCore{
		TransactionQueue: []hivego.HiveTransaction{},
	}
}

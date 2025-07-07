package producers

import "sync"

type transactionQueue struct {
	mtx *sync.Mutex
	buf []transactionRequest
}

func (q *transactionQueue) push(item transactionRequest) {
	q.mtx.Lock()
	defer q.mtx.Unlock()
	q.buf = append(q.buf, item)
}

func (q *transactionQueue) collectBatch() []transactionRequest {
	q.mtx.Lock()
	defer q.mtx.Unlock()

	buf := make([]transactionRequest, len(q.buf))

	copy(buf, q.buf)
	q.buf = q.buf[:0]

	return buf
}

package producers

import (
	"context"
	"fmt"
	"mimic/lib/utils"
	"mimic/modules/db/mimic/blockdb"
	"sync"
	"time"

	"github.com/chebyrash/promise"
)

const (
	blockProducer = "go-mimic-producer"
	blockIdLen    = 16

	merkleRootBlockSize = 32
)

var producer *Producer = nil

type Producer struct {
	stop     context.CancelFunc
	ctx      context.Context
	trxQueue chan transactionRequest
}

func New() *Producer {
	producer = new(Producer)
	return producer
}

// Runs initialization in order of how they are passed in to `Aggregate`
func (p *Producer) Init() error {
	p.ctx, p.stop = context.WithCancel(context.Background())
	p.trxQueue = make(chan transactionRequest)

	return nil
}

// Runs startup and should be non blocking
func (p *Producer) Start() *promise.Promise[any] {
	return utils.PromiseResolve[any](nil)
}

// Runs cleanup once the `Aggregate` is finished
func (p *Producer) Stop() error {
	p.stop()
	return nil
}

func (p *Producer) Produce(interval time.Duration) {
	tick := time.NewTicker(interval)

	latestBlock := &Block{&blockdb.HiveBlock{}}
	err := blockdb.Collection().FindLatestBlock(p.ctx, latestBlock.HiveBlock)
	if err != nil {
		panic(err)
	}

	trxQueue := transactionQueue{
		mtx: new(sync.Mutex),
		buf: make([]transactionRequest, 0, 100),
	}

	for {
		select {
		case <-p.ctx.Done():
			requests := trxQueue.collectBatch()
			p.makeBlock(requests, latestBlock.NextBlock())
			return

		case req := <-p.trxQueue:
			trxQueue.push(req)
			fmt.Printf("Request queued: %v.\n", req)

		case <-tick.C:
			requests := trxQueue.collectBatch()
			latestBlock = p.makeBlock(requests, latestBlock.NextBlock())
		}
	}
}

func (p *Producer) makeBlock(
	requests []transactionRequest,
	block Block,
) *Block {
	fmt.Printf("Making block with with %d requests.\n", len(requests))

	for _, req := range requests {
		req.comm <- BroadcastTransactionResponse{}
		close(req.comm)
		fmt.Println("TODO: make block and send back response.")
	}

	return &block
}

type Witness struct {
	name string
}

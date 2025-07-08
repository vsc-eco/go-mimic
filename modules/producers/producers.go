package producers

import (
	"context"
	"fmt"
	"log/slog"
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	latestBlock := &Block{&blockdb.HiveBlock{}, 0}
	err := blockdb.Collection().FindLatestBlock(ctx, latestBlock.HiveBlock)
	if err != nil {
		panic(err)
	}

	latestBlock.blockNum, err = blockdb.Collection().FindBlockCount(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println()
	fmt.Println("latest block.", "block", latestBlock)
	fmt.Println()

	trxQueue := transactionQueue{
		mtx: new(sync.Mutex),
		buf: make([]transactionRequest, 0, 100),
	}

	for {
		select {
		case <-p.ctx.Done():
			requests := trxQueue.collectBatch()
			if _, err := p.makeBlock(requests, latestBlock.NextBlock()); err != nil {
				slog.Error(
					"Failed to create block.",
					"block",
					latestBlock.HiveBlock,
				)
			}
			return

		case req := <-p.trxQueue:
			trxQueue.push(req)
			fmt.Printf("Request queued: %v.\n", req)

		case <-tick.C:
			requests := trxQueue.collectBatch()

			lastBlock, err := p.makeBlock(requests, latestBlock.NextBlock())
			if err != nil {
				slog.Error(
					"Failed to create block.",
					"block", latestBlock.HiveBlock,
					"err", err,
				)
			} else {
				latestBlock = lastBlock
			}
		}
	}
}

var stubWitness = Witness{
	name: "hive-io-witness",
}

func (p *Producer) makeBlock(
	requests []transactionRequest,
	block Block,
) (*Block, error) {
	fmt.Printf("Making block with with %d requests.\n", len(requests))

	trx := make([]any, len(requests))
	for i := range requests {
		trx[i] = requests[i].transaction
	}

	if err := block.MakeBlock(trx, stubWitness); err != nil {
		return nil, err
	}

	for _, req := range requests {
		req.comm <- BroadcastTransactionResponse{
			BlockNum: block.blockNum,
			// TODO: fill these out
			ID:      "",
			TrxNum:  0,
			Expired: false,
		}
		close(req.comm)
	}

	return &block, nil
}

type Witness struct {
	name string
}

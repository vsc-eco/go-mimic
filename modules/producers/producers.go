package producers

import (
	"context"
	"log/slog"
	"mimic/lib/utils"
	"mimic/modules/db/mimic/blockdb"
	"time"

	"github.com/chebyrash/promise"
	"github.com/vsc-eco/hivego"
)

const (
	blockProducer       = "go-mimic-producer"
	blockIdLen          = 16
	merkleRootBlockSize = 32
)

var producer *Producer = nil

type Producer struct {
	trxQueue chan transactionRequest
}

func New() *Producer {
	producer = new(Producer)
	return producer
}

// Runs initialization in order of how they are passed in to `Aggregate`
func (p *Producer) Init() error {
	p.trxQueue = make(chan transactionRequest, 100) // bufferred

	return nil
}

// Runs startup and should be non blocking
func (p *Producer) Start() *promise.Promise[any] {
	go p.produceBlocks(time.Second * 3)
	return utils.PromiseResolve[any](nil)
}

// Runs cleanup once the `Aggregate` is finished
func (p *Producer) Stop() error {
	return nil
}

func (p *Producer) produceBlocks(interval time.Duration) {
	slog.Debug("Producing blocks.", "interval", interval)

	// get latest block
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	latestBlock := &producerBlock{&blockdb.HiveBlock{}}
	err := blockdb.Collection().QueryHeadBlock(ctx, latestBlock.HiveBlock)
	if err != nil {
		panic(err)
	}

	if _, err := latestBlock.getBlockNum(); err != nil {
		panic(err)
	}

	tick := time.NewTicker(interval)
	for range tick.C {
		requests := p.batchTransactions()

		lastBlock, err := p.makeBlock(requests, latestBlock.next())
		if err != nil {
			slog.Error(
				"Failed to create block.",
				"block", latestBlock.HiveBlock,
				"err", err,
			)
			continue
		}

		latestBlock = lastBlock
	}
}

func (p *Producer) batchTransactions() []transactionRequest {
	requests := make([]transactionRequest, len(p.trxQueue))
	for i := range requests {
		requests[i] = <-p.trxQueue
	}
	return requests
}

var stubWitness = Witness{
	name: "hive-io-witness",
}

func (p *Producer) makeBlock(
	requests []transactionRequest,
	block producerBlock,
) (*producerBlock, error) {
	transactions := make([]hivego.HiveTransaction, 0, len(requests))
	if err := block.sign(transactions, stubWitness); err != nil {
		return nil, err
	}

	if _, err := block.getBlockNum(); err != nil {
		return nil, err
	}

	if err := blockdb.Collection().InsertBlock(block.HiveBlock); err != nil {
		return nil, err
	}

	for _, req := range requests {
		req.comm <- BroadcastTransactionResponse{
			BlockNum: block.BlockNum,
			// TODO: fill these out
			ID:      "",
			TrxNum:  0,
			Expired: false,
		}
		close(req.comm)
	}

	slog.Debug("New block produced.",
		"transactions", len(requests),
		"block-num", block.BlockNum)

	return &block, nil
}

type Witness struct {
	name string
}

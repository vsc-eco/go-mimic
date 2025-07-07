package producers

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"mimic/lib/utils"
	"mimic/modules/db/mimic/blockdb"
	"slices"
	"sync"
	"time"

	"github.com/chebyrash/promise"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/sha3"
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

type Block struct {
	*blockdb.HiveBlock
}

func (b *Block) NextBlock() Block {
	nextBlock := &blockdb.HiveBlock{
		ObjectID: primitive.NilObjectID,
		Previous: b.HiveBlock.BlockID,
	}
	return Block{nextBlock}
}

func (b *Block) MakeBlock(
	transactions []any,
	witness Witness,
) error {
	b.Timestamp = time.Now().Format(time.RFC3339)

	// get block number
	blockCtrBuf, err := hex.DecodeString(b.Previous[:8])
	if err != nil {
		return err
	}

	// add 1 to the block number
	binary.BigEndian.PutUint32(
		blockCtrBuf[:],
		binary.BigEndian.Uint32(blockCtrBuf)+1,
	)

	// get previous block bytes
	previousBlockHash, err := hex.DecodeString(b.Previous[8:])
	if err != nil {
		return err
	}

	// calculate the merkle root
	merkleRoot, err := generateMerkleRoot(transactions)
	if err != nil {
		return err
	}

	// generating ID
	buf := slices.Concat(
		previousBlockHash,
		[]byte(b.Timestamp),
		[]byte(witness.name),
		[]byte(merkleRoot),
	)

	blockDigest := make([]byte, blockIdLen+len(blockCtrBuf))

	// write the incremented block number
	copy(blockDigest[:len(blockCtrBuf)], blockCtrBuf[:])

	// write the digest buf
	sha3.ShakeSum256(blockDigest[len(blockCtrBuf):], buf)

	// write to new block
	b.BlockID = hex.EncodeToString(blockDigest)
	b.Witness = witness.name
	b.Transactions = transactions
	b.MerkleRoot = hex.EncodeToString(merkleRoot)
	// TODO: generate valid transaction id
	b.TransactionIDs = make([]any, len(transactions))

	return nil
}

func generateMerkleRoot(transactions []any) ([]byte, error) {
	// empty merkle tree
	if len(transactions) == 0 {
		return make([]byte, merkleRootBlockSize), nil
	}

	// merkle tree with 1 transaction, just hash the transaction.
	if len(transactions) == 1 {
		bytes, err := encode(transactions[0])
		if err != nil {
			return nil, err
		}

		digest := sha3.Sum256(bytes)
		return digest[:], nil
	}

	// with 2+ transactions
	digests := make(
		[][merkleRootBlockSize]byte,
		len(transactions),
		len(transactions)+1,
	)

	for i, transaction := range transactions {
		bytes, err := encode(transaction)
		if err != nil {
			return nil, err
		}

		digests[i] = sha3.Sum256(bytes)
	}

	// incase the length of transactions is odd, duplicate the last transaction
	if len(digests)&1 == 1 {
		digests = append(digests, digests[len(digests)-1])
	}

	// pair + hash
	for len(digests) > 1 {
		buf := make([][merkleRootBlockSize]byte, len(digests)/2)

		for i := range buf {
			left := digests[i<<1][:]
			right := digests[(i<<1)+1][:]
			buf[i] = sha3.Sum256(slices.Concat(left, right))
		}

		digests = buf
	}

	return digests[0][:], nil
}

func encode(v any) ([]byte, error) {
	buf := bytes.Buffer{}
	if err := gob.NewEncoder(&buf).Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

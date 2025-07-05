package producers

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"mimic/modules/db/mimic/blockdb"
	"slices"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/sha3"
)

const (
	blockProducer = "go-mimic-producer"
	blockIdLen    = 16
)

func MakeBlockInterval(interval time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	latestBlock := &Block{}
	err := blockdb.Collection().FindLatestBlock(ctx, latestBlock.HiveBlock)
	if err != nil {
		panic(err)
	}

	fmt.Println()
	fmt.Println(latestBlock)
	fmt.Println()

	for {
		time.Sleep(interval)
	}
}

type Witness struct {
	name string
}

type Block struct {
	*blockdb.HiveBlock
}

func (b *Block) NextBlock() *Block {
	nextEmptyBlock := &blockdb.HiveBlock{
		ObjectID: primitive.NilObjectID,
		Previous: b.HiveBlock.BlockID,
	}

	return &Block{nextEmptyBlock}
}

func (b *Block) MakeBlock(transactions []any, witness Witness) error {
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

	// TODO: need to append transaction here before calling this function
	merkleRoot := b.generateMerkleRoot()

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
	// TODO: write to `b.TransactionIDs`

	return nil
}

func (b *Block) generateMerkleRoot() []byte {
	merkleRoot := sha3.Sum256([]byte{}) // TODO: calculate merkle root
	b.MerkleRoot = hex.EncodeToString(merkleRoot[:])
	return merkleRoot[:]
}

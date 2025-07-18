package producers

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"mimic/lib/utils"
	"mimic/modules/db/mimic/blockdb"
	"slices"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type producerBlock struct {
	*blockdb.HiveBlock
}

func (b *producerBlock) next() producerBlock {
	nextBlock := &blockdb.HiveBlock{
		ObjectID: primitive.NilObjectID,
		Previous: b.HiveBlock.BlockID,
	}
	return producerBlock{nextBlock}
}

func (b *producerBlock) sign(transactions []any, witness Witness) error {
	b.Timestamp = time.Now().Format(utils.TimeFormat)

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
	bufHash := checksum(buf)
	copy(blockDigest[len(blockCtrBuf):], bufHash[:16])

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

		digest := checksum(bytes)
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

		digests[i] = checksum(bytes)
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
			buf[i] = checksum(slices.Concat(left, right))
		}

		digests = buf
	}

	return digests[0][:], nil
}

func (p *producerBlock) getBlockNum() (uint32, error) {
	if len(p.BlockID) < 8 {
		return 0, errors.New("Invalid block id.")
	}

	blockNumBytes, err := hex.DecodeString(p.BlockID[:8])
	if err != nil {
		return 0, err
	}

	p.BlockNum = binary.BigEndian.Uint32(blockNumBytes)
	return p.BlockNum, nil
}

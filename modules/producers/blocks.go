package producers

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"mimic/lib/utils"
	"mimic/modules/db/mimic/blockdb"
	"slices"
	"time"

	"github.com/decred/dcrd/dcrec/secp256k1/v2"
	"github.com/vsc-eco/hivego"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type producerBlock struct {
	*blockdb.HiveBlock
}

func (b *producerBlock) next() producerBlock {
	nextBlock := &blockdb.HiveBlock{
		ObjectID: primitive.NilObjectID,
		Previous: b.BlockID,
	}
	return producerBlock{nextBlock}
}

func (b *producerBlock) sign(
	transactions []*hivego.HiveTransaction,
	witness *Witness,
) error {
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

	// sign block
	sig, err := secp256k1.SignCompact(
		witness.keyPair.PrivateKey,
		blockDigest,
		true,
	)
	if err != nil {
		return err
	}

	// write to new block
	blockTrxs := utils.Map(
		transactions,
		func(trx *hivego.HiveTransaction) hivego.HiveTransaction { return *trx },
	)

	// transaction IDs
	trxIds, err := utils.TryMap(
		transactions,
		func(trx *hivego.HiveTransaction) (string, error) {
			return trx.GenerateTrxId()
		},
	)
	block := blockdb.HiveBlock{
		ObjectID:         primitive.NilObjectID,
		BlockNum:         b.BlockNum,
		BlockID:          hex.EncodeToString(blockDigest),
		Previous:         b.Previous,
		Timestamp:        time.Now().Format(utils.TimeFormat),
		Witness:          witness.name,
		MerkleRoot:       hex.EncodeToString(merkleRoot),
		Extensions:       []any{},
		WitnessSignature: hex.EncodeToString(sig),
		Transactions:     blockTrxs,
		TransactionIDs:   trxIds,
		SigningKey:       *witness.keyPair.GetPublicKeyString(),
	}

	*b.HiveBlock = block

	return err
}

func generateMerkleRoot(trxs []*hivego.HiveTransaction) ([]byte, error) {
	// empty merkle tree
	if len(trxs) == 0 {
		return make([]byte, merkleRootBlockSize), nil
	}

	// merkle tree with 1 transaction, just hash the transaction.
	if len(trxs) == 1 {
		bytes, err := encode(trxs[0])
		if err != nil {
			return nil, err
		}

		digest := checksum(bytes)
		return digest[:], nil
	}

	// with 2+ transactions
	digests, err := utils.TryMap(trxs, trxCheckSum)
	if err != nil {
		return nil, err
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

func trxCheckSum(
	trx *hivego.HiveTransaction,
) ([merkleRootBlockSize]byte, error) {
	bytes, err := encode(trx)
	if err != nil {
		return [merkleRootBlockSize]byte{}, err
	}
	return checksum(bytes), nil
}

func (p *producerBlock) getBlockNum() (uint32, error) {
	if len(p.BlockID) < 8 {
		return 0, errors.New("invalid block id")
	}

	blockNumBytes, err := hex.DecodeString(p.BlockID[:8])
	if err != nil {
		return 0, err
	}

	p.BlockNum = binary.BigEndian.Uint32(blockNumBytes)
	return p.BlockNum, nil
}

package producers

import (
	"encoding/hex"
	"encoding/json"
	"mimic/modules/db/mimic/blockdb"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const hiveBlockIDLen = 40

var testTransactions = []any{
	&testTransaction{
		ID:     "1",
		From:   "hive-io-from",
		To:     "hive-io-to",
		Amount: 12.0,
	},
	&testTransaction{
		ID:     "2",
		From:   "hive-io-from",
		To:     "hive-io-to",
		Amount: 14.0,
	},
	&testTransaction{
		ID:     "3",
		From:   "hive-io-from",
		To:     "hive-io-to",
		Amount: 16.0,
	},
}

type testTransaction struct {
	ID     string  `json:"id"`
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}

func TestMakeBlock(t *testing.T) {
	var buf []blockdb.HiveBlock

	f, err := os.Open("../../mock/block_api.get_block.json")
	if err != nil {
		panic(err)
	}

	if err := json.NewDecoder(f).Decode(&buf); err != nil {
		panic(err)
	}

	witness := Witness{name: "go-mimic-test"}

	// test for empty merkle tree generation
	firstBlock := producerBlock{&buf[0]}
	err = firstBlock.sign([]any{}, witness)
	assert.Nil(t, err)
	assert.Equal(
		t,
		hex.EncodeToString(make([]byte, merkleRootBlockSize)),
		firstBlock.MerkleRoot,
		"It should contain an emtpy merkle root with no transactions.",
	)

	// test for second block derivation
	trx := make([]any, len(testTransactions))
	copy(trx, testTransactions)

	secondBlock := firstBlock.next()
	assert.Equal(t, firstBlock.BlockID, secondBlock.Previous)

	err = secondBlock.sign(trx, witness)
	assert.Nil(t, err)

	assert.Nil(t, err)
	assert.Equal(
		t,
		"00000001",
		secondBlock.Previous[:8],
		"Reference the previous block corrently.",
	)
	assert.Equal(
		t,
		"00000002",
		secondBlock.BlockID[:8],
		"Incremented the block correctly.",
	)
	assert.Equal(
		t,
		hiveBlockIDLen,
		len(secondBlock.BlockID),
		"Valid length for Hive's Block ID.",
	)
	assert.Equal(
		t,
		witness.name,
		secondBlock.Witness,
		"Witness name is propagated.",
	)
	assert.Equal(
		t,
		len(trx),
		len(secondBlock.Transactions),
		"Transactions are propagated.",
	)
	assert.Equal(
		t,
		len(trx),
		len(secondBlock.TransactionIDs),
		"TransactionIDs are propagated",
	)
	assert.NotEqual(
		t,
		secondBlock.BlockID[8:],
		firstBlock.BlockID[8:],
		"Block IDs should diff.",
	)
	assert.NotEqual(
		t,
		secondBlock.MerkleRoot,
		firstBlock.MerkleRoot,
		"Merkle root should diff.",
	)

	// the merkle root is calculated
	thirdBlock := secondBlock.next()
	trxs := make([]any, len(testTransactions))
	copy(trxs, testTransactions)

	err = thirdBlock.sign(trxs, witness)
	assert.Nil(t, err)
	assert.NotEqual(
		t,
		hex.EncodeToString(make([]byte, merkleRootBlockSize)),
		thirdBlock.MerkleRoot,
		"Merkle root is calculated.",
	)
}

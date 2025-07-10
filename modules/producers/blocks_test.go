package producers

import (
	"encoding/json"
	"mimic/modules/db/mimic/blockdb"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBlockNum(t *testing.T) {
	f, err := os.Open("../../mock/block_api.get_block.json")
	if err != nil {
		panic(err)
	}

	buf := make([]blockdb.HiveBlock, 1)
	if err := json.NewDecoder(f).Decode(&buf); err != nil {
		panic(err)
	}

	block := producerBlock{&buf[0], 0}
	blockNum, err := block.getBlockNum()
	assert.Nil(t, err)
	assert.Equal(t, uint32(0), blockNum)
	assert.Equal(t, blockNum, block.blockNum)

	block = block.next()
	assert.Equal(t, uint32(1), block.blockNum)
}

func TestMerkleRoot(t *testing.T) {
	testTRXs := make([]any, len(testTransactions))
	copy(testTRXs, testTransactions)

	merkleRoot1, err := generateMerkleRoot(testTRXs)
	assert.Nil(t, err)

	merkleRoot2, err := generateMerkleRoot(testTRXs)
	assert.Nil(t, err)

	assert.Equal(
		t,
		merkleRoot1,
		merkleRoot2,
		"merkle roots should be the same with identical transactions.",
	)

	// merkle root should diff with modified transaction
	trx4 := testTransaction{
		ID:     "4",
		From:   "hive-io-from",
		To:     "hive-io-to",
		Amount: 18.0,
	}

	testTRXs[0] = &trx4

	merkleRoot3, err := generateMerkleRoot(testTRXs)
	assert.Nil(t, err)

	assert.NotEqual(
		t,
		merkleRoot1,
		merkleRoot3,
		"merkle roots should diff with modified transactions.",
	)

	// merkle root should diff with more transactions
	testTRXs = append(testTRXs, trx4)
	merkleRoot4, err := generateMerkleRoot(testTRXs)
	assert.Nil(t, err)

	assert.NotEqual(
		t,
		merkleRoot3,
		merkleRoot4,
		"merkle root should diff with more transactions.",
	)
}

package producers

import (
	"encoding/json"
	"mimic/modules/db/mimic/blockdb"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const hiveBlockIDLen = 40

func TestBlockID(t *testing.T) {
	var buf []blockdb.HiveBlock

	f, err := os.Open("../../mock/block_api.get_block.json")
	if err != nil {
		panic(err)
	}

	if err := json.NewDecoder(f).Decode(&buf); err != nil {
		panic(err)
	}

	firstBlock := Block{&buf[0]}
	secondBlock := firstBlock.NextBlock()
	assert.Equal(t, firstBlock.BlockID, secondBlock.Previous)

	witness := Witness{name: "go-mimic-test"}
	transactions := []any{} // TODO: mock some data
	err = secondBlock.MakeBlock(transactions, witness)
	assert.Nil(t, err)

	t.Logf("New block: %s", secondBlock)

	assert.Nil(t, err)
	assert.Equal(t, "00000000", secondBlock.Previous[:8])
	assert.Equal(t, "00000001", secondBlock.BlockID[:8])
	assert.Equal(t, hiveBlockIDLen, len(secondBlock.BlockID))
	assert.Equal(t, witness.name, secondBlock.Witness)
	assert.Equal(t, len(transactions), len(secondBlock.Transactions))
	assert.Equal(t, len(transactions), len(secondBlock.TransactionIDs))
	assert.NotEqual(t, secondBlock.BlockID[8:], firstBlock.BlockID[8:])
	assert.NotEqual(t, secondBlock.Timestamp, firstBlock.Timestamp)
	assert.NotEqual(t, secondBlock.MerkleRoot, firstBlock.MerkleRoot)
}

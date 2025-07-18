package condenser

import (
	"context"
	"mimic/modules/db/mimic/blockdb"
	"mimic/modules/db/mimic/condenserdb"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockDB struct{}

// QueryBlockByBlockNum implements blockdb.BlockQuery.
func (m *mockDB) QueryBlockByBlockNum(*blockdb.HiveBlock, int64) error {
	panic("unimplemented")
}

// QueryBlockByRange implements blockdb.BlockQuery.
func (m *mockDB) QueryBlockByRange(
	blocks *[]blockdb.HiveBlock,
	start int,
	end int,
) error {
	panic("unimplemented")
}

// QueryHeadBlock implements blockdb.BlockQuery.
func (m *mockDB) QueryHeadBlock(
	_ context.Context,
	buf *blockdb.HiveBlock,
) error {
	*buf = blockdb.HiveBlock{
		BlockNum:         1,
		BlockID:          "1",
		Previous:         "0",
		Timestamp:        "2025-07-11T14:27:00-07:00",
		Witness:          "go-mimic-witness",
		MerkleRoot:       "merkleroot",
		Extensions:       []any{},
		WitnessSignature: "go-mimic-witness-sig",
		Transactions:     []any{},
		TransactionIDs:   []any{},
		SigningKey:       "signingkey",
	}

	return nil
}

func TestGetDynamicGlobalProperties(t *testing.T) {
	srv := &Condenser{
		BlockDB:   &mockDB{},
		AccountDB: nil,
	}

	args := make([]string, 0)
	response := &condenserdb.GlobalProperties{}
	headBlock := blockdb.HiveBlock{}

	srv.GetDynamicGlobalProperties(&args, response)

	t.Run("it propagates the correct data.", func(t *testing.T) {
		err := srv.BlockDB.QueryHeadBlock(context.TODO(), &headBlock)
		assert.Nil(t, err)
		assert.Equal(t, headBlock.BlockID, response.HeadBlockID)
		assert.Equal(t, headBlock.BlockNum, response.HeadBlockNumber)
		assert.Equal(t, headBlock.Timestamp, response.Time)
		assert.Equal(t, headBlock.Witness, response.CurrentWitness)
	})
}

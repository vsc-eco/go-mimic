package blockdb

import (
	"encoding/json"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Block struct {
	BlockId               string   `json:"block_id"`
	Previous              string   `json:"previous"`
	Timestamp             string   `json:"timestamp"`
	Witness               string   `json:"witness"`
	TransactionMerkleRoot string   `json:"transaction_merkle_root"`
	Extensions            []string `json:"extensions"`
	WitnessSignature      string   `json:"witness_signature"`
	Transactions          []string `json:"transactions"`
	SigningKey            string   `json:"signing_key"`
	TransactionIds        []string `json:"transaction_ids"`
}

type HiveBlock struct {
	ObjectID primitive.ObjectID `json:"-" bson:"_id,omitempty"`

	BlockNum         uint32 `json:"-"                       bson:"block_num"`
	BlockID          string `json:"block_id"`
	Previous         string `json:"previous"`
	Timestamp        string `json:"time"`
	Witness          string `json:"witness"`
	MerkleRoot       string `json:"transaction_merkle_root"`
	Extensions       []any  `json:"extensions"`
	WitnessSignature string `json:"witness_signature"`
	Transactions     []any  `json:"transactions"`
	TransactionIDs   []any  `json:"transaction_ids"`
	SigningKey       string `json:"signing_key"`
}

func (h *HiveBlock) String() string {
	jsonString, _ := json.MarshalIndent(h, "", "  ")
	return fmt.Sprintf("[%s] %s", h.ObjectID.String(), string(jsonString))
}

package blocks

type HiveBlock struct {
	BlockId        string `json:"block_id" bson:"id"`
	Witness        string `json:"witness" bson:"witness"`
	Timestamp      string `json:"timestamp" bson:"ts"`
	MerkleRoot     string `json:"merkle_root" bson:"merkle_root"`
	Previous       string `json:"previous" bson:"previous"`
	TransactionIds string `json:"transaction_ids" bson:"tx_ids"`

	Height int64 `json:"height" bson:"height"`
}

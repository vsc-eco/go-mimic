package blockdb

type Block struct {
	BlockNum              uint64   `json:"-" bson:"block_num"`
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
	BlockId        string `json:"block_id" bson:"id"`
	Witness        string `json:"witness" bson:"witness"`
	Timestamp      string `json:"timestamp" bson:"ts"`
	MerkleRoot     string `json:"merkle_root" bson:"merkle_root"`
	Previous       string `json:"previous" bson:"previous"`
	TransactionIds string `json:"transaction_ids" bson:"tx_ids"`

	Height int64 `json:"height" bson:"height"`
}

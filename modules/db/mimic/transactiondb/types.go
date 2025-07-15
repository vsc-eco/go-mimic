package transactiondb

import "go.mongodb.org/mongo-driver/bson/primitive"

type Transaction struct {
	ObjectID primitive.ObjectID `json:"-" bson:"_id,omitempty"`

	RefBlockNum    uint16   `json:"ref_block_num"`
	RefBlockPrefix uint32   `json:"ref_block_prefix"`
	Expiration     string   `json:"expiration"`
	Operations     []any    `json:"operations"`
	Extensions     []any    `json:"extensions"`
	Signatures     []string `json:"signatures"`
}

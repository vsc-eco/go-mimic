package producers

import (
	"crypto/sha256"
	"encoding/json"
)

var checksum = sha256.Sum256

func encode(v any) ([]byte, error) {
	return json.Marshal(v)
}

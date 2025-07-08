package producers

import (
	"encoding/json"
)

func encode(v any) ([]byte, error) {
	return json.Marshal(v)
}

package producers

import (
	"bytes"
	"encoding/gob"
)

func encode(v any) ([]byte, error) {
	buf := bytes.Buffer{}
	if err := gob.NewEncoder(&buf).Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

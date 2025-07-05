package mock

import (
	"encoding/json"
	"fmt"
	"os"
)

func GetMockData(buf any, apiMethod string) error {
	f, err := os.Open(fmt.Sprintf("mock/%s.json", apiMethod))
	if err != nil {
		return err
	}

	return json.NewDecoder(f).Decode(buf)
}

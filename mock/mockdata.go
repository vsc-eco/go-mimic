package mock

import (
	"encoding/json"
	"fmt"
	"os"
)

func GetMockData(buf any, mockJsonFile string) error {
	f, err := os.Open(fmt.Sprintf("mock/%s", mockJsonFile))
	if err != nil {
		return err
	}

	return json.NewDecoder(f).Decode(buf)
}

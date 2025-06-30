package mock

import (
	"encoding/json"
	"os"
)

func GetMockData[T any](mockPath string) (map[string]T, error) {
	// TODO: save to mock data to db
	f, err := os.ReadFile(mockPath)
	if err != nil {
		return nil, err
	}

	data := make(map[string]T)
	if err := json.Unmarshal(f, &data); err != nil {
		return nil, err
	}

	return data, nil
}

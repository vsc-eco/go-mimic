package services

import (
	"encoding/json"
	"os"
)

func getMockData[T any](mockPath string) (map[string]T, error) {
	// TODO: propagate this into db
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

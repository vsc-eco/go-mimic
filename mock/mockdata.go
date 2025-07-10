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
	defer f.Close()

	return json.NewDecoder(f).Decode(buf)
}

type seedUserCredentials struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func LoadSeedUserCredentials() ([]seedUserCredentials, error) {
	f, err := os.Open("mock/account_seed.json")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := make([]seedUserCredentials, 0, 3)
	return buf, json.NewDecoder(f).Decode(&buf)
}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"mimic/modules/db/mimic/accountdb"
	"sync"
)

func main() {
	wg := new(sync.WaitGroup)

	accounts, err := accountdb.GetSeedAccounts()
	if err != nil {
		log.Fatal(err)
	}

	wg.Add(len(accounts))

	for _, account := range accounts {
		go func(wg *sync.WaitGroup, account *accountdb.Account) {
			defer wg.Done()

			requestBody := [2]any{
				"account_create",
				map[string]any{
					"fee": map[string]any{
						"amount":    "0",
						"precision": 3,
						"nai":       "@@000000021",
					},
					"creator":          "go-mimic-root",
					"new_account_name": account.Name,
					"owner":            account.Owner,
					"active":           account.Active,
					"posting":          account.Posting,
					"memo_key":         account.MemoKey,
					"json_metadata":    "{}",
				},
			}

			// TODO: send these in the POST request
			jsonEncoded, err := json.MarshalIndent(requestBody, "", "  ")
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(string(jsonEncoded))
		}(wg, &account)
	}

	wg.Wait()
}

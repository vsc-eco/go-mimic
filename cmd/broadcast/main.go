package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mimic/modules/api/services/broadcastops"
	"mimic/modules/db/mimic/accountdb"
	"net/http"
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

			trx := broadcastops.BroadcastParam[broadcastops.AccountCreateParam]{
				Action: "account_create",
				Param: broadcastops.AccountCreateParam{
					Fee: broadcastops.AccountCreateFee{
						Amount:    "0",
						Precision: 3,
						Nai:       "@@000000021",
					},
					Creator:        "go-mimic-root",
					NewAccountName: account.Name,
					Owner:          account.Owner,
					Active:         account.Active,
					Posting:        account.Posting,
					MemoKey:        account.MemoKey,
					JsonMetadata:   "{}",
				},
			}

			req := map[string]any{
				"jsonrpc": "2.0",
				"method":  "condenser_api.account_create",
				"params":  [2]any{trx.Action, trx.Param},
				"id":      1,
			}

			// TODO: send these in the POST request
			jsonEncoded, err := json.Marshal(&req)
			if err != nil {
				log.Fatal(err)
			}

			reqBody := bytes.NewBuffer(jsonEncoded)

			cx := http.Client{}
			res, err := cx.Post("http://localhost:3000", "application/json", reqBody)
			if err != nil {
				log.Fatal(err)
			}
			defer res.Body.Close()

			var buf []byte
			if _, err := io.ReadFull(res.Body, buf); err != nil {
				log.Fatal(err)
			}

			fmt.Println(string(buf))
		}(wg, &account)
	}

	wg.Wait()
}

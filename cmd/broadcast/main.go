package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"mimic/modules/api/services/condenser"
	"mimic/modules/db/mimic/accountdb"
	"mimic/modules/db/mimic/condenserdb"
	"mimic/modules/db/mimic/transactiondb"
	"net/http"
	"time"
)

var httpClient = http.Client{}

const url = "http://localhost:3000"

func main() {
	accounts, err := accountdb.GetSeedAccounts()
	if err != nil {
		log.Fatal(err)
	}

	for _, account := range accounts {
		err := createAccount(account)
		if err != nil {
			log.Fatal(err)
		}

		// trx := condenser.BroadcastParam[condenser.AccountCreateParam]{
		// 	Action: "account_create",
		// 	Param: condenser.AccountCreateParam{
		// 		Fee: condenser.AccountCreateFee{
		// 			Amount:    "0",
		// 			Precision: 3,
		// 			Nai:       "@@000000021",
		// 		},
		// 		Creator:        "go-mimic-root",
		// 		NewAccountName: account.Name,
		// 		Owner:          account.Owner,
		// 		Active:         account.Active,
		// 		Posting:        account.Posting,
		// 		MemoKey:        account.MemoKey,
		// 		JsonMetadata:   "{}",
		// 	},
		// }

		// req := map[string]any{
		// 	"jsonrpc": "2.0",
		// 	"method":  "condenser_api.account_create",
		// 	"params":  [2]any{trx.Action, trx.Param},
		// 	"id":      1,
		// }

		// // TODO: send these in the POST request
		// jsonEncoded, err := json.Marshal(&req)
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// reqBody := bytes.NewBuffer(jsonEncoded)

		// cx := http.Client{}
		// res, err := cx.Post(
		// 	"http://localhost:3000",
		// 	"application/json",
		// 	reqBody,
		// )
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// defer res.Body.Close()

		// var buf []byte
		// if _, err := io.ReadFull(res.Body, buf); err != nil {
		// 	log.Fatal(err)
		// }

		// fmt.Println(string(buf))
	}
}

func createAccount(account accountdb.Account) error {
	headBlock, err := sendRequest[condenserdb.GlobalProperties](
		[]any{},
		"condenser_api.get_dynamic_global_properties",
	)
	if err != nil {
		return err
	}

	// making a transaction, based on this implementation
	// https://github.com/vsc-eco/hivego/blob/fa6c9e2c8be757b260a9b48b7d206fa02f8cfde9/signer.go#L22
	trx := transactiondb.Transaction{}

	trx.RefBlockNum = uint16(headBlock.HeadBlockNumber & 0xffff)
	hbidB, err := hex.DecodeString(headBlock.HeadBlockID)
	if err != nil {
		return err
	}
	trx.RefBlockPrefix = binary.LittleEndian.Uint32(hbidB[4:])

	ts, err := time.Parse(time.RFC3339, headBlock.Time)
	if err != nil {
		return err
	}
	trx.Expiration = ts.Add(time.Second * 30).Format(time.RFC3339)

	trx.Operations = make([]any, 1)
	trx.Operations[0] = [2]any{
		"account_create",
		condenser.AccountCreateParam{
			Fee: condenser.AccountCreateFee{
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

	// TODO: sign the transaction then make a request

	// send request
	req := map[string]any{
		"jsonrpc": "2.0",
		"method":  "condenser_api.account_create",
		"id":      1,
		"params":  map[string]transactiondb.Transaction{"trx": trx},
	}

	// TODO: send these in the POST request
	jsonEncoded, err := json.Marshal(&req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(jsonEncoded))

	// reqBody := bytes.NewBuffer(jsonEncoded)

	return nil
}

func sendRequest[T any](payload any, jsonrpcMethod string) (*T, error) {
	reqBody := map[string]any{
		"jsonrpc": "2.0",
		"method":  jsonrpcMethod,
		"params":  payload,
		"id":      1,
	}

	v, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	body := bytes.NewBuffer(v)
	response, err := httpClient.Post(url, "application/json", body)
	if err != nil {
		return nil, err
	}

	var buf struct {
		Result T `json:"result"`
	}
	if err := json.NewDecoder(response.Body).Decode(&buf); err != nil {
		return nil, err
	}

	return &buf.Result, nil
}

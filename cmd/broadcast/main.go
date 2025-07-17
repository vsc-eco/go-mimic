package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"log"
	"mimic/lib/utils"
	"mimic/modules/db/mimic/accountdb"
	"mimic/modules/db/mimic/condenserdb"
	"mimic/modules/transactions"
	"net/http"
	"time"

	"github.com/vsc-eco/hivego"
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
	trx := hivego.HiveTransaction{}

	trx.RefBlockNum = uint16(headBlock.HeadBlockNumber & 0xffff)
	hbidB, err := hex.DecodeString(headBlock.HeadBlockID)
	if err != nil {
		return err
	}
	trx.RefBlockPrefix = binary.LittleEndian.Uint32(hbidB[4:])

	ts, err := time.Parse(utils.TimeFormat, headBlock.Time)
	if err != nil {
		return err
	}
	trx.Expiration = ts.Add(time.Second * 30).Format(utils.TimeFormat)

	// trx.Operations = make([]hivego.HiveOperation, 1)
	// trx.OperationsJs = make([][2]any, 1)
	trx.Signatures = make([]string, 1)
	op := transactions.AccountCreateOp{
		Fee: transactions.AccountCreateFee{
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
	}
	appendTrx(&trx, &op)

	// TODO: sign the transaction
	serializedTx, err := hivego.SerializeTx(trx)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(len(serializedTx))

	/*
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

	*/
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

func appendTrx(trx *hivego.HiveTransaction, op hivego.HiveOperation) {
	trx.Operations = append(trx.Operations, op)
	trx.OperationsJs = append(trx.OperationsJs, [2]any{op.OpName(), op})
}

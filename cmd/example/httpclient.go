package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mimic/modules/db/mimic/condenserdb"
	"net/http"

	"github.com/vsc-eco/hivego"
)

type apiClient struct {
	httpClient *http.Client
	buf        *bytes.Buffer
}

type broadcastTransactionRequest struct {
	JsonRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Method  string `json:"method"`
	Params  struct {
		Trx hivego.HiveTransaction `json:"trx"`
	} `json:"params"`
}

func (a *apiClient) broadcastSync(
	trx *hivego.HiveTransaction,
	keyPair *hivego.KeyPair,
) ([]byte, error) {
	sig, err := trx.Sign(*keyPair)
	if err != nil {
		return nil, err
	}
	trx.AddSig(sig)

	a.buf.Reset()
	defer a.buf.Reset()

	req := &broadcastTransactionRequest{
		JsonRPC: "2.0",
		ID:      1,
		Method:  "condenser_api.broadcast_transaction_synchronous",
		Params: struct {
			Trx hivego.HiveTransaction `json:"trx"`
		}{
			Trx: *trx,
		},
	}

	if err := json.NewEncoder(a.buf).Encode(req); err != nil {
		return nil, err
	}

	res, err := a.httpClient.Post(apiUrl, apiContentType, a.buf)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if _, err := io.Copy(a.buf, res.Body); err != nil {
		return nil, err
	}

	if res.StatusCode >= 300 || res.StatusCode < 200 {
		return nil, errors.New(a.buf.String())
	}

	out := make([]byte, len(a.buf.Bytes()))
	copy(out, a.buf.Bytes())

	return out, nil
}

func (a *apiClient) queryGlobalProps() (*condenserdb.GlobalProperties, error) {
	a.buf.Reset()
	defer a.buf.Reset()

	request := jsonrpcRequest{
		JsonRPC: "2.0",
		ID:      1,
		Method:  "condenser_api.get_dynamic_global_properties",
		Params:  []string{},
	}

	if err := json.NewEncoder(a.buf).Encode(&request); err != nil {
		return nil, err
	}

	res, err := a.httpClient.Post(apiUrl, apiContentType, a.buf)
	if err != nil {
		return nil, err
	}

	var responseBodyRaw struct {
		Result condenserdb.GlobalProperties `json:"result"`
	}

	if err := json.NewDecoder(res.Body).Decode(&responseBodyRaw); err != nil {
		return nil, err
	}

	return &responseBodyRaw.Result, nil
}

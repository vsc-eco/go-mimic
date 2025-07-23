package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mimic/lib/utils"
	"mimic/modules/api/services/condenser"
	"mimic/modules/db/mimic/condenserdb"
	"net/http"
	"os"
	"time"

	"github.com/vsc-eco/hivego"
)

var (
	errUnsupportedMethod = errors.New("unsupported method")

	httpClient = new(http.Client)
)

const (
	apiUrl         = "http://0.0.0.0:3000"
	apiContentType = "application/json"
)

type jsonrpcMethod interface {
	hivego.HiveOperation
	params() any
}

var methods = map[string]jsonrpcMethod{
	"condenser_api.custom_json":    &customJson{},
	"condenser_api.account_update": &accountUpdate{},
}

func main() {
	args := os.Args

	requestsExec := make([]string, 0, len(args)-1)
	if len(args) == 1 {
		for k := range methods {
			requestsExec = append(requestsExec, k)
		}
	} else {
		requestsExec = append(requestsExec, args[1:]...)
	}

	for _, r := range requestsExec {
		log.Println("executing:", r)
		content, err := executeMethod(r)
		if err != nil {
			log.Println("method failed:", err)
			continue
		}

		log.Println("method ok:")
		fmt.Println(string(content))
	}
}

func executeMethod(jsonrpcMethod string) ([]byte, error) {
	b, ok := methods[jsonrpcMethod]
	if !ok {
		return nil, errUnsupportedMethod
	}

	buf := new(bytes.Buffer)
	if err := makeRequest(buf, jsonrpcMethod, b); err != nil {
		return nil, err
	}

	res, err := httpClient.Post(apiUrl, apiContentType, buf)
	if err != nil {
		return nil, err
	}

	buf.Reset()
	if _, err := io.Copy(buf, res.Body); err != nil {
		return nil, err
	}

	if res.StatusCode >= 300 || res.StatusCode < 200 {
		return nil, errors.New(buf.String())
	}

	return buf.Bytes(), nil
}

func makeRequest[T hivego.HiveOperation](
	buf *bytes.Buffer,
	method string,
	b T,
) error {
	buf.Reset()

	// get global prop
	request := map[string]any{

		"jsonrpc": "2.0",
		"id":      1,
		"method":  "condenser_api.get_dynamic_global_properties",
		"params":  []any{},
	}

	if err := json.NewEncoder(buf).Encode(&request); err != nil {
		return err
	}

	res, err := httpClient.Post(apiUrl, apiContentType, buf)
	if err != nil {
		return err
	}

	responseBodyRaw := struct {
		Result condenserdb.GlobalProperties `json:"result"`
	}{}

	if err := json.NewDecoder(res.Body).Decode(&responseBodyRaw); err != nil {
		return err
	}

	// make transaction
	props := responseBodyRaw.Result

	refBlockNum := uint16(props.HeadBlockNumber & 0xffff)

	hbidB, err := hex.DecodeString(props.HeadBlockID)
	if err != nil {
		return err
	}

	refBlockPrefix := binary.LittleEndian.Uint32(hbidB[4:])

	exp, err := time.Parse(utils.TimeFormat, props.Time)
	if err != nil {
		return err
	}

	expiration := exp.Add(30 * time.Second).Format(utils.TimeFormat)

	trx := condenser.Transaction[T]{
		Expiration:           expiration,
		RefBlockNum:          refBlockNum,
		RefBlockPrefix:       refBlockPrefix,
		Extensions:           []any{},
		Operations:           []T{b},
		Signatures:           []string{},
		RequiredAuths:        []string{},
		RequiredPostingAuths: []string{},
	}

	// broadcast op
	request = map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  method,
		"params":  map[string]any{"trx": trx},
	}
	return json.NewEncoder(buf).Encode(&request)
}

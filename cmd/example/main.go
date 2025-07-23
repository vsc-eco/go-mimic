package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"mimic/lib/utils"
	"net/http"
	"os"
	"time"

	"github.com/vsc-eco/hivego"
)

const (
	apiUrl         = "http://0.0.0.0:3000"
	apiContentType = "application/json"
)

var (
	testUserKey *hivego.KeyPair

	errUnsupportedMethod = errors.New("unsupported method")
	httpClient           = &apiClient{new(http.Client), new(bytes.Buffer)}
)

type jsonrpcMethod interface {
	params() hivego.HiveOperation
}

var transactionType = map[string]jsonrpcMethod{
	"custom_json":    &customJson{},
	"account_update": &accountUpdate{},
}

func init() {
	privKeyWIF := utils.EnvOrPanic("TEST_OWNER_PRIVATE_KEY")

	var err error
	testUserKey, err = hivego.KeyPairFromWif(privKeyWIF)
	if err != nil {
		panic(err)
	}
}

func main() {
	args := os.Args

	requestsExec := make([]string, 0, len(args)-1)
	if len(args) == 1 {
		for k := range transactionType {
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
	b, ok := transactionType[jsonrpcMethod]
	if !ok {
		return nil, errUnsupportedMethod
	}

	trx, err := makeTransaction(b)
	if err != nil {
		return nil, err
	}

	return httpClient.broadcastSync(trx, testUserKey)
}

type jsonrpcRequest struct {
	JsonRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Method  string `json:"method"`
	Params  any    `json:"params"`
}

func makeTransaction(b jsonrpcMethod) (*hivego.HiveTransaction, error) {
	// query global props
	props, err := httpClient.queryGlobalProps()
	if err != nil {
		return nil, err
	}

	// make transaction
	refBlockNum := uint16(props.HeadBlockNumber & 0xffff)

	hbidB, err := hex.DecodeString(props.HeadBlockID)
	if err != nil {
		return nil, err
	}

	refBlockPrefix := binary.LittleEndian.Uint32(hbidB[4:])

	exp, err := time.Parse(utils.TimeFormat, props.Time)
	if err != nil {
		return nil, err
	}

	expiration := exp.Add(30 * time.Second).Format(utils.TimeFormat)

	trx := hivego.HiveTransaction{
		Expiration:     expiration,
		RefBlockNum:    refBlockNum,
		RefBlockPrefix: refBlockPrefix,
		Extensions:     []string{},
		Operations:     []hivego.HiveOperation{b.params()},
		Signatures:     []string{},
	}

	return &trx, nil
}

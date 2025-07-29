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

/**
- [Main Consumer]
- binary for main process
- Separate binary for tests/consumer application
-- hits the API
-- post transactions
- [Tests]
- Creating an account with generated keys -- keys pregenerated
- Block streaming (showing it can stream blocks from the RPC)
- Creating various transaction with proper validation
-- Test for known correct transactions
-- Test for negative cases (i.e incorrect transactions); Signature validation; Format validation
- Verify streaming can receive said transactions after 3-6s
- [Possibly later]
- Test get_account API
-- account_update transaction should update get_account API response (i.e modify database records)
- Balance updates for transfers




[validatioen]
- format?
- signature
- state transition
-- ie user transfer, need to verify they have the sufficient balance

*/

func init() {
	privKeyWIF := utils.EnvOrPanic("TEST_POSTING_KEY_PRIVATE")

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

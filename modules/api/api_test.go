package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"github.com/vsc-eco/go-mimic/modules/api"
)

func TestSimplePost(t *testing.T) {
	apiServer := api.NewAPIServer()

	apiServer.Init()
	apiServer.Start()

	reqBody := api.Request{
		JsonRpc: "2.0",
		Method:  "condenser_api.get_block",
		Params: map[string]interface{}{
			"a": 10,
			"b": 20,
		},
	}
	data, err := json.Marshal(reqBody)
	if err != nil {
		log.Fatalf("failed to marshal request: %v", err)
	}

	resp, err := http.Post("http://localhost:3000/", "application/json", bytes.NewReader(data))
	if err != nil {
		log.Fatalf("RPC call failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println("body", string(body))
	var rpcResp api.Response
	if err := json.Unmarshal(body, &rpcResp); err != nil {
		log.Fatalf("invalid response: %v", err)
	}
	if rpcResp.Error != "" {
		log.Fatalf("RPC error: %s", rpcResp.Error)
	}

	// client, err := rpc.DialHTTPPath("tcp", "localhost"+":3000", "/")

	// if err != nil {
	// 	t.Fatalf("Failed to connect to server: %v", err)
	// }

	// args := &api.Args{}
	// var reply int64

	// err = client.Call("condenser.get_block", args, &reply)
	// fmt.Println("err", err)

	// fmt.Println("reply", reply)
}

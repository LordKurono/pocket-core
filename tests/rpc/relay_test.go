package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/pokt-network/pocket-core/config"
	"github.com/pokt-network/pocket-core/node"
	"github.com/pokt-network/pocket-core/rpc/relay"
	"github.com/pokt-network/pocket-core/rpc/shared"
	"github.com/pokt-network/pocket-core/service"
)

/*
Unit test for the relay functionality
*/
func TestRelay(t *testing.T) {
	node.DWL().Add("DEVID1")
	// grab the hosted chains via file
	if err := node.CFIle(config.GlobalConfig().CFile); err != nil {
		t.Fatalf(err.Error())
	}
	node.TestChains()
	fmt.Println(node.Chains())
	// Start server instance
	go http.ListenAndServe(":"+config.GlobalConfig().RRPCPort, shared.Router(relay.Routes()))
	// @ Url
	u := "http://localhost:" + config.GlobalConfig().RRPCPort + "/v1/relay/"
	// Setup relay
	r := service.Relay{}
	// add blockchain value
	r.Blockchain = "ethereum"
	// add netid value
	r.NetworkID = "1"
	// add version value
	r.Version = "1.0"
	// add data value
	r.Data = "{\"jsonrpc\":\"2.0\",\"method\":\"web3_clientVersion\",\"params\":[],\"id\":67}"
	// add developer id
	r.DevID = "DEVID1"
	// convert structure to json
	j, err := json.Marshal(r)
	// handle error
	if err != nil {
		t.Fatalf("Cannot convert struct to json " + err.Error())
	}
	// create new post request
	req, err := http.NewRequest("POST", u, bytes.NewBuffer(j))
	// hanlde error
	if err != nil {
		t.Fatalf("Cannot create post request " + err.Error())
	}
	// setup header for json data
	req.Header.Set("Content-Type", "application/json")
	// setup http client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Unable to do post request " + err.Error())
	}
	// get body of response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Unable to unmarshal response: " + err.Error())
	}
	t.Log(string(body))
}

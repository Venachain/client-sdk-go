package client

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/PlatONE_Network/PlatONE-SDK-Go/common"
	"github.com/PlatONE_Network/PlatONE-SDK-Go/platone/rpc"
	"github.com/stretchr/testify/assert"
)

func InitContractClient() (common.TxParams, ContractClient) {
	txparam := common.TxParams{}
	codePath1 := "/Users/cxh/Downloads/example/example.wasm"
	keyfile := "/Users/cxh/go/src/github.com/PlatONE_Network/PlatONE-Go/release/linux/conf/keyfile.json"
	abiPath := "/Users/cxh/Downloads/example/example.cpp.abi.json"
	PassPhrase := "0"
	vm := "wasm"
	url := URL{
		IP:      "127.0.0.1",
		RPCPort: 6791,
	}
	rpc, _ := rpc.DialContext(context.Background(), url.GetEndpoint())
	pc := Client{
		RpcClient:   rpc,
		Passphrase:  PassPhrase,
		KeyfilePath: keyfile,
		URL:         &url,
	}
	contract := ContractClient{
		Client:   &pc,
		CodePath: codePath1,
		AbiPath:  abiPath,
		Vm:       vm,
	}
	return txparam, contract
}

func TestContractClientDeploy(t *testing.T) {
	txparam, contract := InitContractClient()
	var consParams []string
	result, _ := contract.Deploy(txparam, consParams)
	assert.True(t, result != nil)
}

func TestContractClient_ListContractMethods(t *testing.T) {
	_, contract := InitContractClient()
	result, _ := contract.ListContractMethods()
	assert.True(t, result != nil)
}

func TestContractClient_Execute(t *testing.T) {
	txparam, contract := InitContractClient()
	funcname := "setEvidence"
	funcparam := []string{"1", "23"}
	addr := "0x24de7156a5e973a5d1a7ee82a27772ca0a22fdb5"
	result, _ := contract.Execute(txparam, funcname, funcparam, addr)
	assert.True(t, result != nil)
}

func TestContractClient_GetReceipt(t *testing.T) {
	txhash := "0x35972a847e8c29148976e8a1884665732c862706c71bbaaf573e8cbd432ba921"
	_, contractClient := InitContractClient()
	result, _ := contractClient.GetReceipt(txhash)
	if result != nil {
		resultBytes, _ := json.MarshalIndent(result, "", "\t")
		fmt.Printf("result:\n%s\n", resultBytes)
	}
	assert.True(t, result != nil)
}

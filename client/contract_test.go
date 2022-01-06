package client

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"git-c.i.wxblockchain.com/PlatONE/src/node/client-sdk-go/common"
	"git-c.i.wxblockchain.com/PlatONE/src/node/client-sdk-go/platone/rpc"
	"github.com/stretchr/testify/assert"
)

func InitContractClient() (common.TxParams, ContractClient) {
	txparam := common.TxParams{}
	keyfile := "/Users/cxh/go/src/github.com/PlatONE_Network/PlatONE-Go/release/linux/conf/keyfile.json"
	//abiPath := "/Users/cxh/go/src/PlatONE-Go/release/linux/conf/contracts/evidenceManager.cpp.abi.json"
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
		Client:  &pc,
		//AbiPath: abiPath,
		Vm:      vm,
	}
	return txparam, contract
}

func TestContractClient_Deploy(t *testing.T) {
	codePath := "/Users/cxh/Downloads/example/example.wasm"
	txparam, contract := InitContractClient()
	var consParams []string
	result, _ := contract.Deploy(context.Background(), txparam, codePath, consParams)
	fmt.Println(result)
	assert.True(t, result != nil)
}

func TestContractClient_ListContractMethods(t *testing.T) {
	_, contract := InitContractClient()
	result, _ := contract.ListContractMethods()
	fmt.Println(result.ListAbiFuncName())
	assert.True(t, result != nil)
}

func TestContractClient_Execute(t *testing.T) {
	txparam, contract := InitContractClient()
	funcname := "setEvidence"
	funcparam := []string{"1", "data"}
	addr := "0x35853e5643104cd96bd4590f5d4466c577786cfe"
	result, _ := contract.Execute(context.Background(), txparam, funcname, funcparam, addr)
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

func TestContractClient_CnsExecute(t *testing.T) {
	txparam, contract := InitContractClient()
	funcname := "saveEvidence"
	funcparam := []string{"2", "23"}
	//cns := "wxbc1"
	result, _ := contract.Execute(context.Background(), txparam, funcname, funcparam, "0x0000000000000000000000000000000000000099")
	fmt.Println(result)
	contract.Client.RpcClient.Close()

	assert.True(t, result != nil)
}

func TestContractClient_GetReceipt2Execute(t *testing.T) {
	txparam, contract := InitContractClient()
	funcname := "getEvidence"
	funcparam := []string{"2"}
	result, _ := contract.Execute(context.Background(), txparam, funcname, funcparam, "0x0000000000000000000000000000000000000099")
	fmt.Println(result)
	contract.Client.RpcClient.Close()
	assert.True(t, result != nil)
}


func TestContractClient_IsFuncNameInContract(t *testing.T) {
	_, contract := InitContractClient()
	funcname := "setEvidence"
	result, _ := contract.IsFuncNameInContract(funcname)
	fmt.Println(result)
	assert.True(t, result != false)
}

func TestContractClient_NFTMint(t *testing.T) {
	txparam, contract := InitContractClient()
	funcname := "mint"
	funcparam := []string{}
	funcparam = append(funcparam,"{\"method\":\"mint\", \"data\":[{\"name\":\"abcd\",\"symbol\":\"ab\",\"description\":\"abcdf1\",\"iprice\": 100, \"price\":100,\"url\":\"ww.qwe.com\",\"property\":\"p11210\",\"others\":\"123\"}]}")
	funcparam = append(funcparam,"")
	result, _ := contract.Execute(context.Background(), txparam, funcname, funcparam, "0x0000000000000000000000000000000000000012")
	fmt.Println(result)
	contract.Client.RpcClient.Close()
	assert.True(t, result != nil)
}

func TestContractClient_NFTGet(t *testing.T) {
	txparam, contract := InitContractClient()
	funcname := "getNFTById"
	funcparam := []string{}
	funcparam = append(funcparam,"acd088948164483a37f73989f1c27b47a5586ebd1d004d9f6dfdcef17082f275")
	result, _ := contract.Execute(context.Background(), txparam, funcname, funcparam, "0x0000000000000000000000000000000000000012")
	fmt.Println(result)
	contract.Client.RpcClient.Close()
	assert.True(t, result != nil)
}


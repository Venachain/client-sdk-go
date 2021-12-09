package client

import (
	"context"
	"fmt"
	"testing"

	common_sdk "git-c.i.wxblockchain.com/PlatONE/src/node/client-sdk-go/common"
	"git-c.i.wxblockchain.com/PlatONE/src/node/client-sdk-go/precompiled/syscontracts"
	"github.com/stretchr/testify/assert"
)

func InitNodeWallClient() (common_sdk.TxParams, NodeClient) {
	txparam, contract := InitContractClient()
	contract.AbiPath = ""
	//contract.CodePath = ""
	client := NodeClient{
		ContractClient: contract,
		NodeName:       "test",
	}
	return txparam, client
}

func TestNodeClient_NodeAdd(t *testing.T) {
	txparam, client := InitNodeWallClient()
	requestNodeInfo := syscontracts.NodeInfo{
		Name:       client.NodeName,
		ExternalIP: "127.0.0.1",
		InternalIP: "127.0.0.1",
		PublicKey:  "feffe2938d427088f5fcce94a9245760b92c468d3ca25ab5ef2b1cdccf0ed911963b74ca2dffef20ef135966e34ebcc905d1f12c1df09f05974a617cf8afe8e8",
	}
	result, _ := client.NodeAdd(context.Background(), txparam, requestNodeInfo)
	fmt.Println(result)
	assert.True(t, result != "")
}

func TestNodeClient_NodeDelete(t *testing.T) {
	txparam, client := InitNodeWallClient()
	result, _ := client.NodeDelete(context.Background(), txparam)
	fmt.Println(result)
	assert.True(t, result != "")
}

func TestNodeClient_NodeUpdate(t *testing.T) {
	txparam, client := InitNodeWallClient()
	request := syscontracts.NodeUpdateInfo{
		Desc: "this is a desc",
	}
	result, _ := client.NodeUpdate(context.Background(), txparam, request)
	fmt.Println(result)
	assert.True(t, result != "")
}

func TestNodeClient_NodeQuery(t *testing.T) {
	txparam, client := InitNodeWallClient()
	request := syscontracts.NodeQueryInfo{
		Name: "test",
	}
	result, _ := client.NodeQuery(context.Background(), txparam, &request)
	fmt.Println(result)
	assert.True(t, result != "")
}

func TestNodeClient_NodeStat(t *testing.T) {
	txparam, client := InitNodeWallClient()
	request := syscontracts.NodeStatInfo{
		Status: 1,
		Type:   1,
	}
	result, _ := client.NodeStat(context.Background(), txparam, &request)
	fmt.Println(result)
	assert.True(t, result == 1)
}

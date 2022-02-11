package client

import (
	"context"
	"testing"

	"git-c.i.wxblockchain.com/PlatONE/src/node/client-sdk-go/log"
	"git-c.i.wxblockchain.com/PlatONE/src/node/client-sdk-go/precompiled/syscontracts"
	"github.com/stretchr/testify/assert"
)

func InitNodeWallClient() (*NodeClient, error) {
	keyfile := "/Users/cxh/go/src/github.com/PlatONE_Network/PlatONE-Go/release/linux/conf/keyfile.json"
	PassPhrase := "0"
	url := URL{
		IP:      "127.0.0.1",
		RPCPort: 6791,
	}
	contractName := "test"
	return NewNodeClient(context.Background(), url, keyfile, PassPhrase, contractName)
}

func TestNodeClient_NodeAdd(t *testing.T) {
	client, err := InitNodeWallClient()
	if err != nil {
		log.Error("error:%v", err)
		return
	}
	defer client.RpcClient.Close()
	requestNodeInfo := syscontracts.NodeInfo{
		Name:       client.NodeName,
		ExternalIP: "127.0.0.1",
		InternalIP: "127.0.0.1",
		PublicKey:  "feffe2938d427088f5fcce94a9245760b92c468d3ca25ab5ef2b1cdccf0ed911963b74ca2dffef20ef135966e34ebcc905d1f12c1df09f05974a617cf8afe8e8",
	}
	result, _ := client.NodeAdd(context.Background(), requestNodeInfo)
	log.Info("result:%v", result)
	assert.True(t, result != "")
}

func TestNodeClient_NodeDelete(t *testing.T) {
	client, err := InitNodeWallClient()
	if err != nil {
		log.Error("error:%v", err)
		return
	}
	defer client.RpcClient.Close()
	result, _ := client.NodeDelete(context.Background())
	log.Info("result:%v", result)
	assert.True(t, result != "")
}

func TestNodeClient_NodeUpdate(t *testing.T) {
	client, err := InitNodeWallClient()
	if err != nil {
		log.Error("error:%v", err)
		return
	}
	defer client.RpcClient.Close()
	request := syscontracts.NodeUpdateInfo{
		Desc: "this is a desc",
	}
	result, _ := client.NodeUpdate(context.Background(), request)
	log.Info("result:%v", result)
	assert.True(t, result != "")
}

func TestNodeClient_NodeQuery(t *testing.T) {
	client, err := InitNodeWallClient()
	if err != nil {
		log.Error("error:%v", err)
		return
	}
	defer client.RpcClient.Close()
	request := syscontracts.NodeQueryInfo{
		Name: "test",
	}
	result, _ := client.NodeQuery(context.Background(), &request)
	log.Info("result:%v", result)
	assert.True(t, result != "")
}

func TestNodeClient_NodeStat(t *testing.T) {
	client, err := InitNodeWallClient()
	if err != nil {
		log.Error("error:%v", err)
		return
	}
	defer client.RpcClient.Close()
	request := syscontracts.NodeStatInfo{
		Status: 1,
		Type:   1,
	}
	result, _ := client.NodeStat(context.Background(), &request)
	log.Info("result:%v", result)
	assert.True(t, result == 1)
}

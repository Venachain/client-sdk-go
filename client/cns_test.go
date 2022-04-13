package client

import (
	"context"
	"testing"

	"git-c.i.wxblockchain.com/vena/src/client-sdk-go/log"
	"github.com/stretchr/testify/assert"
)

// 如果没有abipath 和codepath 的话，可以设置为空
func InitCnsClient() (*CnsClient, error) {
	keyfile := "/Users/cxh/go/src/github.com/PlatONE_Network/PlatONE-Go/release/linux/conf/keyfile.json"
	PassPhrase := "0"
	url := URL{
		IP:      "127.0.0.1",
		RPCPort: 6791,
	}
	name := "wxbc"
	return NewCnsClient(context.Background(), url, keyfile, PassPhrase, name)
}

func TestCnsClient_CnsRegister(t *testing.T) {
	client, err := InitCnsClient()
	if err != nil {
		log.Error("error:%v", err)
		return
	}
	defer client.RpcClient.Close()

	result, _ := client.CnsRegister(context.Background(), "1.0.0.0", "0x6988decc03a2d38888534ad0b4a33a267b34807d")
	log.Info("result:%v", result)
	assert.True(t, result != "")
}

func TestCnsClient_CnsResolve(t *testing.T) {
	client, err := InitCnsClient()
	if err != nil {
		log.Error("error:%v", err)
		return
	}
	defer client.RpcClient.Close()

	result, _ := client.CnsResolve(context.Background(), "1.0.0.0")
	log.Info("result:%v", result)
	assert.True(t, result != "")
}

func TestCnsClient_CnsQueryAll(t *testing.T) {
	client, err := InitCnsClient()
	if err != nil {
		log.Error("error:%v", err)
		return
	}
	defer client.RpcClient.Close()

	result, _ := client.CnsQueryAll(context.Background())
	log.Info("result:%v", result)
	assert.True(t, result != "")
}

func TestCnsClient_CnsQueryByName(t *testing.T) {
	client, err := InitCnsClient()
	if err != nil {
		log.Error("error:%v", err)
		return
	}
	defer client.RpcClient.Close()

	result, _ := client.CnsQueryByName(context.Background())
	log.Info("result:%v", result)
	assert.True(t, result != "")
}

func TestCnsClient_CnsQueryByAddress(t *testing.T) {
	client, err := InitCnsClient()
	if err != nil {
		log.Error("error:%v", err)
		return
	}
	defer client.RpcClient.Close()
	address := "0x6988decc03a2d38888534ad0b4a33a267b34807d"
	result, _ := client.CnsQueryByAddress(context.Background(), address)
	log.Info("result:%v", result)
	assert.True(t, result != "")
}

func TestCnsClient_CnsQueryByAccount(t *testing.T) {
	client, err := InitCnsClient()
	if err != nil {
		log.Error("error:%v", err)
		return
	}
	defer client.RpcClient.Close()
	address := "0xdbd41e01e0e4a51fdb03c6152c50df071207a04b"
	result, _ := client.CnsQueryByAccount(context.Background(), address)
	log.Info("result:%v", result)
	assert.True(t, result != "")
}

func TestCnsClient_CnsStateByAddress(t *testing.T) {
	client, err := InitCnsClient()
	if err != nil {
		log.Error("error:%v", err)
		return
	}
	defer client.RpcClient.Close()
	address := "0x6988decc03a2d38888534ad0b4a33a267b34807d"
	result, _ := client.CnsStateByAddress(context.Background(), address)
	log.Info("result:%v", result)
	assert.True(t, result == 1)
}

func TestCnsClient_CnsState(t *testing.T) {
	client, err := InitCnsClient()
	if err != nil {
		log.Error("error:%v", err)
		return
	}
	defer client.RpcClient.Close()
	result, _ := client.CnsState(context.Background())
	log.Info("result:%v", result)
	assert.True(t, result == 1)
}

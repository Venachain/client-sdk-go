package client

import (
	"context"
	"testing"

	"git-c.i.wxblockchain.com/PlatONE/src/node/client-sdk-go/log"
	"github.com/stretchr/testify/assert"
)

// 如果没有abipath 和codepath 的话，可以设置为空
func InitFireWallClient() (*FireWallClient, error) {
	keyfile := "/Users/cxh/go/src/github.com/PlatONE_Network/PlatONE-Go/release/linux/conf/keyfile.json"
	PassPhrase := "0"
	url := URL{
		IP:      "127.0.0.1",
		RPCPort: 6791,
	}
	accountAddress := "0x6988decc03a2d38888534ad0b4a33a267b34807d"
	return NewFireWallClient(context.Background(), url, keyfile, PassPhrase, accountAddress)
}

func TestAccountClient_FwStart(t *testing.T) {
	client, err := InitFireWallClient()
	if err != nil {
		log.Error("error:%v", err)
		return
	}
	defer client.RpcClient.Close()
	result, err := client.FwStart(context.Background())
	if err != nil {
		log.Error("error:%v", err)
		return
	}
	log.Info("result:%v", result)
	assert.True(t, result != "")
}

func TestAccountClient_FwStatusFwStatus(t *testing.T) {
	client, err := InitFireWallClient()
	if err != nil {
		log.Error("error:%v", err)
		return
	}
	defer client.RpcClient.Close()

	result, err := client.FwStatus(context.Background())
	if err != nil {
		log.Error("error:%v", err)
		return
	}
	log.Info("result:%v", result)
	assert.True(t, result != "")
}

func TestAccountClient_FwClose(t *testing.T) {
	client, err := InitFireWallClient()
	if err != nil {
		log.Error("error:%v", err)
		return
	}
	defer client.RpcClient.Close()
	result, _ := client.FwClose(context.Background())
	log.Info("result:%v", result)
	assert.True(t, result != "")
}

func TestFireWallClient_FwExport(t *testing.T) {
	client, err := InitFireWallClient()
	if err != nil {
		log.Error("error:%v", err)
		return
	}
	defer client.RpcClient.Close()
	result, _ := client.FwExport(context.Background(), "./config1")
	log.Info("result:%v", result)
	assert.True(t, result == true)
}

func TestFireWallClient_FwNew(t *testing.T) {
	client, err := InitFireWallClient()
	if err != nil {
		log.Error("error:%v", err)
		return
	}
	defer client.RpcClient.Close()
	result, _ := client.FwNew(context.Background(), "accept", "0x3fcaa0a86dfbbe105c7ed73ca505c7a59c579667", "atransfer")
	log.Info("result:%v", result)
	assert.True(t, result != "")
}

func TestFireWallClient_FwDelete(t *testing.T) {
	client, err := InitFireWallClient()
	if err != nil {
		log.Error("error:%v", err)
		return
	}
	defer client.RpcClient.Close()
	result, _ := client.FwDelete(context.Background(), "accept", "0x3fcaa0a86dfbbe105c7ed73ca505c7a59c579667", "atransfer")
	log.Info("result:%v", result)
	assert.True(t, result != "")
}

func TestFireWallClient_FwClear(t *testing.T) {
	client, err := InitFireWallClient()
	if err != nil {
		log.Error("error:%v", err)
		return
	}
	defer client.RpcClient.Close()
	result, _ := client.FwClear(context.Background(), "accept")
	log.Info("result:%v", result)
	assert.True(t, result != "")
}

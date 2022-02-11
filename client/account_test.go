package client

import (
	"context"
	"testing"

	"git-c.i.wxblockchain.com/PlatONE/src/node/client-sdk-go/log"
	"github.com/stretchr/testify/assert"
)

// 如果没有abipath 和codepath 的话，可以设置为空
func InitAccountClient() (*AccountClient, error) {
	keyfile := "/Users/cxh/go/src/github.com/PlatONE_Network/PlatONE-Go/release/linux/conf/keyfile.json"
	PassPhrase := "0"
	url := URL{
		IP:      "127.0.0.1",
		RPCPort: 6791,
	}
	accountAddress := "0xdbd41e01e0e4a51fdb03c6152c50df071207a04b"
	return NewAccountClient(context.Background(), url, keyfile, PassPhrase, accountAddress)
}

func TestAccountClient_UserAdd(t *testing.T) {
	client, err := InitAccountClient()
	if err != nil {
		log.Info("error:%v", err)
		return
	}
	defer client.RpcClient.Close()
	result, err := client.UserAdd(context.Background(), "Alice", "110", "", "")
	if err != nil {
		log.Info("err:%v", err)
	}
	log.Info(result)
	assert.True(t, result != "")
}

func TestAccountClient_UserUpdate(t *testing.T) {
	client, err := InitAccountClient()
	if err != nil {
		log.Info("error:%v", err)
		return
	}
	defer client.RpcClient.Close()
	result, _ := client.UserUpdate(context.Background(), "13556672653", "test@163.com", "wxbc2")
	log.Info(result)
	assert.True(t, result != "")
}

func TestAccountClient_QueryUser(t *testing.T) {
	client, err := InitAccountClient()
	if err != nil {
		log.Info("error:%v", err)
		return
	}
	defer client.RpcClient.Close()
	result, err := client.QueryUser(context.Background(), "Alice")
	if err != nil {
		log.Info("error:%v", err)
		return
	}
	log.Info(result)
	assert.True(t, result != "")
}

func TestAccountClient_Lock(t *testing.T) {
	client, err := InitAccountClient()
	if err != nil {
		log.Info("error:%v", err)
		return
	}
	defer client.RpcClient.Close()
	result, _ := client.Lock(context.Background())
	log.Info("result:%v", result)
	assert.True(t, result == true)
}

func TestAccountClient_UnLock(t *testing.T) {
	client, err := InitAccountClient()
	if err != nil {
		log.Info("error:%v", err)
		return
	}
	defer client.RpcClient.Close()
	passphrase := "0"
	result, _ := client.UnLock(context.Background(), passphrase)
	log.Info("result:%v", result)
	assert.True(t, result == true)
}

func TestAccountClient_CreateAccountk(t *testing.T) {
	client, err := InitAccountClient()
	if err != nil {
		log.Info("error:%v", err)
		return
	}
	defer client.RpcClient.Close()
	result, err := client.CreateAccount(context.Background(), "0")
	if err != nil {
		log.Info("error:%v", err)
		return
	}
	log.Info(result.Hex())
	assert.True(t, result != nil)
}

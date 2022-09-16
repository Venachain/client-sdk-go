package client

import (
	"context"
	"testing"

	"github.com/Venachain/client-sdk-go/log"
)

// 如果没有abipath 和codepath 的话，可以设置为空
func InitAccountClient() (*AccountClient, error) {
	keyfile := "/Users/cxh/go/src/VenaChain/venachain/release/linux/conf/keyfile.json"
	PassPhrase := "0"
	url := URL{
		IP:      "127.0.0.1",
		RPCPort: 6791,
	}
	accountAddress := "0x85a4b8ad3a023fab30146fed114ea7cd6f8a4193"
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
}

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/Venachain/client-sdk-go/log"
	"github.com/Venachain/client-sdk-go/venachain/common"
	"github.com/stretchr/testify/assert"
)

func InitClient() (*Client, error) {
	keyfile := "/Users/cxh/go/src/VenaChain/venachain/release/linux/conf/keyfile.json"
	PassPhrase := "0"
	url := URL{
		IP:      "127.0.0.1",
		RPCPort: 6791,
	}
	return NewClient(context.Background(), url, keyfile, PassPhrase)
}

func TestRpcCall_newAccount(t *testing.T) {
	client, err := InitClient()
	if err != nil {
		log.Error("error:%v", err)
		return
	}
	funcName := "personal_newAccount"
	params := "0"
	result, err := client.RpcCall(context.Background(), funcName, params)
	if err != nil {
		log.Error("error:%v", err)
		return
	}
	log.Info("result:%v", result)
	// 结果转换
	var res common.Address
	if err = json.Unmarshal(result, &res); err != nil {
		log.Error("error:%v", err)
		return
	}
	assert.True(t, result != nil)
}

func TestRpcCall_lockAccount(t *testing.T) {
	client, err := InitClient()
	if err != nil {
		log.Error("error:%v", err)
		return
	}
	funcName := "personal_lockAccount"
	params := "0xdbd41e01e0e4a51fdb03c6152c50df071207a04b"
	result, err := client.RpcCall(context.Background(), funcName, params)
	if err != nil {
		log.Error("error:%v", err)
		return
	}
	log.Info("result:%v", result)
	// 结果转换，如果是go通用数据类型，可用getRpcResult函数获取
	res := getRpcResult(result, "bool")
	log.Info("result:%v", res)
	assert.True(t, result != nil)
}

func TestRpcCall_getBlockByHash(t *testing.T) {
	client, err := InitClient()
	if err != nil {
		log.Error("error:%v", err)
		return
	}
	res, err := client.GetBlockByHash("0xe8413de3f95aa00fb219e3f62af483e05ebb7c90d11bb33a144b6dc3b5e491cf")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("res", res.Transactions)
}

func TestRpcCall_getBlockByNumber(t *testing.T) {
	client, err := InitClient()
	if err != nil {
		log.Error("error:%v", err)
		return
	}
	funcName := "eth_getBlockByNumber"
	result, err := client.RpcClient.Call(context.Background(), funcName, "latest", false)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(result))
}

func TestClient_GetBlockAllByHash(t *testing.T) {
	client, err := InitClient()
	if err != nil {
		log.Error("error:%v", err)
		return
	}
	blockHash := "0xf8f3b61ee5d739d460f527a640938ab8a231525685f29e7c203a1ce92410ecfa"

	res, err := client.GetBlockAllByHash(blockHash)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}

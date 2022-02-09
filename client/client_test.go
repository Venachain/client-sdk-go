package client

import (
	"context"
	"encoding/json"
	"testing"

	"git-c.i.wxblockchain.com/PlatONE/src/node/client-sdk-go/log"
	"git-c.i.wxblockchain.com/PlatONE/src/node/client-sdk-go/platone/common"
	"github.com/stretchr/testify/assert"
)

func InitClient() (*Client, error) {
	keyfile := "/Users/cxh/go/src/github.com/PlatONE_Network/PlatONE-Go/release/linux/conf/keyfile.json"
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

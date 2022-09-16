package client

import (
	"context"
	"fmt"
	"testing"

	"github.com/Venachain/client-sdk-go/log"
	"github.com/stretchr/testify/assert"
)

func InitSysconfigClient() (*SysConfigClient, error) {
	keyfile := "/Users/cxh/go/src/github.com/PlatONE_Network/PlatONE-Go/release/linux/conf/keyfile.json"
	PassPhrase := "0"
	url := URL{
		IP:      "127.0.0.1",
		RPCPort: 6791,
	}
	return NewSysConfigClient(context.Background(), url, keyfile, PassPhrase)
}

func TestSysConfigClient_SetSysConfig(t *testing.T) {
	client, err := InitSysconfigClient()
	if err != nil {
		log.Error("error:%v", err)
		return
	}
	defer client.RpcClient.Close()

	request := SysConfigParam{
		Tx_gaslimit:    "1900000000",
		Block_gaslimit: "20000000000",
		Empty_block:    "notallow-empty",
	}
	result, _ := client.SetSysConfig(context.Background(), request)
	log.Info("result:%v", result)
	assert.True(t, result != nil)
}

func TestSysConfigClient_GetIsApproveDeployedContract(t *testing.T) {
	client, err := InitSysconfigClient()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.RpcClient.Close()
	result, _ := client.GetIsApproveDeployedContract(context.Background())
	log.Info("result:%v", result)
	assert.True(t, result == 0)
}

func TestSysConfigClient_GetVRFParams(t *testing.T) {
	client, err := InitSysconfigClient()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.RpcClient.Close()
	result, _ := client.GetVRFParams(context.Background())
	log.Info("result:%v", result)
	assert.True(t, result != "")
}

func TestSysConfigClient_GetIsTxUseGas(t *testing.T) {
	client, err := InitSysconfigClient()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.RpcClient.Close()
	result, _ := client.GetIsTxUseGas(context.Background())
	log.Info("result:%v", result)
	assert.True(t, result == 0)
}

func TestSysConfigClient_GetIsProduceEmptyBlock(t *testing.T) {
	client, err := InitSysconfigClient()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.RpcClient.Close()
	result, _ := client.GetIsProduceEmptyBlock(context.Background())
	log.Info("result:%v", result)
	assert.True(t, result == 0)
}

func TestSysConfigClient_GetBlockGasLimit(t *testing.T) {
	client, err := InitSysconfigClient()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.RpcClient.Close()
	result, _ := client.GetBlockGasLimit(context.Background())
	log.Info("result:%v", result)
	assert.True(t, result != 0)
}

func TestSysConfigClient_GetTxGasLimit(t *testing.T) {
	client, err := InitSysconfigClient()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.RpcClient.Close()
	result, _ := client.GetTxGasLimit(context.Background())
	log.Info("result:%v", result)
	assert.True(t, result != 0)
}

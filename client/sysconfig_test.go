package client

import (
	"context"
	"fmt"
	"testing"

	common_sdk "git-c.i.wxblockchain.com/PlatONE/src/node/client-sdk-go/common"
	"github.com/stretchr/testify/assert"
)

// 如果没有abipath 和codepath 的话，可以设置为空
func InitSysconfigClient() (common_sdk.TxParams, SysConfigClient) {
	txparam, contract := InitContractClient()
	contract.AbiPath = ""
	//contract.CodePath = ""
	client := SysConfigClient{
		ContractClient: contract,
	}
	return txparam, client
}

func TestSysConfigClient_SetSysConfig(t *testing.T) {
	txparam, client := InitSysconfigClient()
	request := SysConfigParam{
		Tx_gaslimit:    "1900000000",
		Block_gaslimit: "20000000000",
		Empty_block:    "notallow-empty",
	}
	result, _ := client.SetSysConfig(context.Background(), txparam, request)
	fmt.Println(result)
	assert.True(t, result != nil)
}

func TestSysConfigClient_GetIsApproveDeployedContract(t *testing.T) {
	txparam, client := InitSysconfigClient()
	result, _ := client.GetIsApproveDeployedContract(context.Background(), txparam)
	fmt.Println(result)
	assert.True(t, result == 0)
}

func TestSysConfigClient_GetVRFParams(t *testing.T) {
	txparam, client := InitSysconfigClient()
	result, _ := client.GetVRFParams(context.Background(), txparam)
	fmt.Println(result)
	assert.True(t, result != "")
}

func TestSysConfigClient_GetIsTxUseGas(t *testing.T) {
	txparam, client := InitSysconfigClient()
	result, _ := client.GetIsTxUseGas(context.Background(), txparam)
	fmt.Println(result)
	assert.True(t, result == 0)
}

func TestSysConfigClient_GetIsProduceEmptyBlock(t *testing.T) {
	txparam, client := InitSysconfigClient()
	result, _ := client.GetIsProduceEmptyBlock(context.Background(), txparam)
	fmt.Println(result)
	assert.True(t, result == 0)
}

func TestSysConfigClient_GetBlockGasLimit(t *testing.T) {
	txparam, client := InitSysconfigClient()
	result, _ := client.GetBlockGasLimit(context.Background(), txparam)
	fmt.Println(result)
	assert.True(t, result != 0)
}

func TestSysConfigClient_GetTxGasLimit(t *testing.T) {
	txparam, client := InitSysconfigClient()
	result, _ := client.GetTxGasLimit(context.Background(), txparam)
	fmt.Println(result)
	assert.True(t, result != 0)
}

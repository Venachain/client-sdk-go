package client

import (
	"context"
	"fmt"
	"testing"

	common_sdk "git-c.i.wxblockchain.com/PlatONE/src/node/client-sdk-go/common"
	"github.com/stretchr/testify/assert"
)

// 如果没有abipath 和codepath 的话，可以设置为空
func InitCnsClient() (common_sdk.TxParams, CnsClient) {
	txparam, contract := InitContractClient()
	contract.AbiPath = ""
	client := CnsClient{
		ContractClient: contract,
		name:           "wxbc1",
	}
	return txparam, client
}

func TestCnsClient_CnsRegister(t *testing.T) {
	txparam, client := InitCnsClient()

	result, _ := client.CnsRegister(context.Background(), txparam, "1.0.0.0", "0xf2aa70bfcfbc6095f4f9e19d01b79de3604c4447")
	fmt.Println(result)
	assert.True(t, result != "")
}

func TestCnsClient_CnsResolve(t *testing.T) {
	txparam, client := InitCnsClient()

	result, _ := client.CnsResolve(context.Background(), txparam, "1.0.0.0")
	fmt.Println(result)
	assert.True(t, result != "")
}

func TestCnsClient_CnsQueryAll(t *testing.T) {
	txparam, client := InitCnsClient()

	result, _ := client.CnsQueryAll(context.Background(), txparam)
	fmt.Println(result)
	assert.True(t, result != "")
}

func TestCnsClient_CnsQueryByName(t *testing.T) {
	txparam, client := InitCnsClient()

	result, _ := client.CnsQueryByName(context.Background(), txparam)
	fmt.Println(result)
	assert.True(t, result != "")
}

func TestCnsClient_CnsQueryByAddress(t *testing.T) {
	txparam, client := InitCnsClient()
	address := "0x7311adfe02f1d027c7c896ceeb45e59ec7282a80"
	result, _ := client.CnsQueryByAddress(context.Background(), txparam, address)
	fmt.Println(result)
	assert.True(t, result != "")
}

func TestCnsClient_CnsQueryByAccount(t *testing.T) {
	txparam, client := InitCnsClient()
	address := "0x3fcaa0a86dfbbe105c7ed73ca505c7a59c579667"
	result, _ := client.CnsQueryByAccount(context.Background(), txparam, address)
	fmt.Println(result)
	assert.True(t, result != "")
}

func TestCnsClient_CnsStateByAddress(t *testing.T) {
	txparam, client := InitCnsClient()
	address := "0x3fcaa0a86dfbbe105c7ed73ca505c7a59c579667"
	result, _ := client.CnsStateByAddress(context.Background(), txparam, address)
	fmt.Println(result)
	assert.True(t, result == 0)
}

func TestCnsClient_CnsState(t *testing.T) {
	txparam, client := InitCnsClient()
	result, _ := client.CnsState(context.Background(), txparam)
	fmt.Println(result)
	assert.True(t, result == 1)
}

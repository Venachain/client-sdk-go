package client

import (
	"context"
	"fmt"
	"testing"

	common_sdk "github.com/PlatONE_Network/PlatONE-SDK-Go/common"

	"github.com/stretchr/testify/assert"
)

// 如果没有abipath 和codepath 的话，可以设置为空
func InitFireWallClient() (common_sdk.TxParams, FireWallClient) {
	txparam, contract := InitContractClient()
	contract.AbiPath = ""
	//contract.CodePath = ""
	client := FireWallClient{
		ContractClient:  contract,
		ContractAddress: "0x26527b41f4a5d9a1e0652c97fd629ced6f7a2263",
	}
	return txparam, client
}

func TestAccountClient_FwStatusFwStatus(t *testing.T) {
	txparam, client := InitFireWallClient()
	result, _ := client.FwStatus(context.Background(), txparam)
	fmt.Println(result)
	assert.True(t, result != "")
}

func TestAccountClient_FwStart(t *testing.T) {
	txparam, client := InitFireWallClient()
	result, _ := client.FwStart(context.Background(), txparam)
	fmt.Println(result)
	assert.True(t, result != "")
}

func TestAccountClient_FwClose(t *testing.T) {
	txparam, client := InitFireWallClient()
	result, _ := client.FwClose(context.Background(), txparam)
	fmt.Println(result)
	assert.True(t, result != "")
}

func TestFireWallClient_FwExport(t *testing.T) {
	txparam, client := InitFireWallClient()
	result, _ := client.FwExport(context.Background(), txparam, "./config1")
	fmt.Println(result)
	assert.True(t, result == true)
}

func TestFireWallClient_FwNew(t *testing.T) {
	txparam, client := InitFireWallClient()
	result, _ := client.FwNew(context.Background(), txparam, "accept", "0x3fcaa0a86dfbbe105c7ed73ca505c7a59c579667", "atransfer")
	fmt.Println(result)
	assert.True(t, result != "")
}

func TestFireWallClient_FwDelete(t *testing.T) {
	txparam, client := InitFireWallClient()
	result, _ := client.FwDelete(context.Background(), txparam, "accept", "0x3fcaa0a86dfbbe105c7ed73ca505c7a59c579667", "atransfer")
	fmt.Println(result)
	assert.True(t, result != "")
}

func TestFireWallClient_FwClear(t *testing.T) {
	txparam, client := InitFireWallClient()
	result, _ := client.FwClear(context.Background(), txparam, "accept")
	fmt.Println(result)
	assert.True(t, result != "")
}

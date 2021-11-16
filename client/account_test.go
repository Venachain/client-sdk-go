package client

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	common_sdk "github.com/PlatONE_Network/PlatONE-SDK-Go/common"
	"github.com/PlatONE_Network/PlatONE-SDK-Go/platone/common"
)

// 如果没有abipath 和codepath 的话，可以设置为空
func InitAccountClient() (common_sdk.TxParams, AccountClient) {
	txparam, contract := InitContractClient()
	contract.AbiPath = ""
	client := AccountClient{
		ContractClient: contract,
		Address:        common.HexToAddress("3fcaa0a86dfbbe105c7ed73ca505c7a59c579667"),
	}
	return txparam, client
}

func TestAccountClient_UserAdd(t *testing.T) {
	txparam, client := InitAccountClient()
	result, _ := client.UserAdd(context.Background(), txparam, "Alice1", "110", "", "")
	fmt.Println(result)
	assert.True(t, result != "")
}

func TestAccountClient_UserUpdate(t *testing.T) {
	txparam, client := InitAccountClient()
	result, _ := client.UserUpdate(context.Background(), txparam, "13556672653", "test@163.com", "wxbc2")
	fmt.Println(result)
	assert.True(t, result != "")
}

func TestAccountClient_QueryUser(t *testing.T) {
	txparam, client := InitAccountClient()
	result, _ := client.QueryUser(context.Background(), txparam, "Alice")
	fmt.Println(result)
	assert.True(t, result != "")
}

func TestAccountClient_Lock(t *testing.T) {
	_, client := InitAccountClient()
	result, _ := client.Lock(context.Background())
	fmt.Println(result)
	assert.True(t, result == true)
}

func TestAccountClient_UnLock(t *testing.T) {
	_, client := InitAccountClient()
	result, _ := client.UnLock(context.Background())
	fmt.Println(result)
	assert.True(t, result == false)
}

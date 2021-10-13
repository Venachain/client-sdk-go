package client

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	common_sdk "github.com/PlatONE_Network/PlatONE-SDK-Go/common"
)

// 如果没有abipath 和codepath 的话，可以设置为空
func InitRoleClient() (common_sdk.TxParams, RoleClient) {
	txparam, contract := InitContractClient()
	contract.AbiPath = ""
	contract.CodePath = ""
	client := RoleClient{
		ContractClient: contract,
	}
	return txparam, client
}

func TestRoleClient_SetSuperAdmin(t *testing.T) {
	txparam, client := InitRoleClient()
	result, _ := client.SetSuperAdmin(txparam)
	fmt.Println(result)
	assert.True(t, result != "")
}

func TestRoleClient_TransferSuperAdmin(t *testing.T) {
	txparam, client := InitRoleClient()
	addr := "0x8d4d2ed9ca6c6279bab46be1624cf7adbab89e18"
	result, _ := client.TransferSuperAdmin(txparam, addr)
	fmt.Println(result)
	assert.True(t, result != "")
}

func TestRoleClient_AddChainAdmin(t *testing.T) {
	txparam, client := InitRoleClient()
	addr := "0x8d4d2ed9ca6c6279bab46be1624cf7adbab89e18"
	result, _ := client.AddChainAdmin(txparam, addr)
	fmt.Println(result)
	assert.True(t, result != "")
}

func TestRoleClient_DelChainAdmin(t *testing.T) {
	txparam, client := InitRoleClient()
	addr := "0x8d4d2ed9ca6c6279bab46be1624cf7adbab89e18"
	result, _ := client.DelChainAdmin(txparam, addr)
	fmt.Println(result)
	assert.True(t, result != "")
}

func TestRoleClient_GetAddrListOfRole(t *testing.T) {
	txparam, client := InitRoleClient()
	role := "SUPER_ADMIN"
	result, _ := client.GetAddrListOfRole(txparam, role)
	fmt.Println(result)
	assert.True(t, result != "")
}

func TestRoleClient_GetRoles(t *testing.T) {
	txparam, client := InitRoleClient()
	addr := "0x3fcaa0a86dfbbe105c7ed73ca505c7a59c579667"
	result, _ := client.GetRoles(txparam, addr)
	fmt.Println(result)
	assert.True(t, result != "")
}

func TestRoleClient_HasRole(t *testing.T) {
	txparam, client := InitRoleClient()
	addr := "0x3fcaa0a86dfbbe105c7ed73ca505c7a59c579667"
	role := "CHAIN_ADMIN"
	result, _ := client.HasRole(txparam, addr, role)
	fmt.Println(result)
	assert.True(t, result != 0)
}

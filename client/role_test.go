package client

import (
	"context"
	"testing"

	"git-c.i.wxblockchain.com/PlatONE/src/node/client-sdk-go/log"
	"github.com/stretchr/testify/assert"
)

func InitRoleClient() (*RoleClient, error) {
	keyfile := "/Users/cxh/go/src/github.com/PlatONE_Network/PlatONE-Go/release/linux/conf/keyfile.json"
	PassPhrase := "0"
	url := URL{
		IP:      "127.0.0.1",
		RPCPort: 6791,
	}
	return NewRoleClient(context.Background(), url, keyfile, PassPhrase)
}

func TestRoleClient_SetSuperAdmin(t *testing.T) {
	client, err := InitRoleClient()
	if err != nil {
		log.Error("error:%v", err)
		return
	}
	defer client.RpcClient.Close()

	result, _ := client.SetSuperAdmin(context.Background())
	log.Info("result:%v", result)
	assert.True(t, result != "")
}

func TestRoleClient_TransferSuperAdmin(t *testing.T) {
	client, err := InitRoleClient()
	if err != nil {
		log.Error("error:%v", err)
		return
	}
	defer client.RpcClient.Close()
	addr := "0x8d4d2ed9ca6c6279bab46be1624cf7adbab89e18"
	result, _ := client.TransferSuperAdmin(context.Background(), addr)
	log.Info("result:%v", result)
	assert.True(t, result != "")
}

func TestRoleClient_AddChainAdmin(t *testing.T) {
	client, err := InitRoleClient()
	if err != nil {
		log.Error("error:%v", err)
		return
	}
	defer client.RpcClient.Close()
	addr := "0x8d4d2ed9ca6c6279bab46be1624cf7adbab89e18"
	result, _ := client.AddChainAdmin(context.Background(), addr)
	log.Info("result:%v", result)
	assert.True(t, result != "")
}

func TestRoleClient_DelChainAdmin(t *testing.T) {
	client, err := InitRoleClient()
	if err != nil {
		log.Error("error:%v", err)
		return
	}
	defer client.RpcClient.Close()
	addr := "0x8d4d2ed9ca6c6279bab46be1624cf7adbab89e18"
	result, _ := client.DelChainAdmin(context.Background(), addr)
	log.Info("result:%v", result)
	assert.True(t, result != "")
}

func TestRoleClient_GetAddrListOfRole(t *testing.T) {
	client, err := InitRoleClient()
	if err != nil {
		log.Error("error:%v", err)
		return
	}
	defer client.RpcClient.Close()
	role := "SUPER_ADMIN"
	result, _ := client.GetAddrListOfRole(context.Background(), role)
	log.Info("result:%v", result)
	assert.True(t, result != "")
}

func TestRoleClient_GetRoles(t *testing.T) {
	client, err := InitRoleClient()
	if err != nil {
		log.Error("error:%v", err)
		return
	}
	defer client.RpcClient.Close()
	addr := "0x3fcaa0a86dfbbe105c7ed73ca505c7a59c579667"
	result, _ := client.GetRoles(context.Background(), addr)
	log.Info("result:%v", result)
	assert.True(t, result != "")
}

func TestRoleClient_HasRole(t *testing.T) {
	client, err := InitRoleClient()
	if err != nil {
		log.Error("error:%v", err)
		return
	}
	defer client.RpcClient.Close()
	addr := "0x3fcaa0a86dfbbe105c7ed73ca505c7a59c579667"
	role := "CHAIN_ADMIN"
	result, _ := client.HasRole(context.Background(), addr, role)
	log.Info("result:%v", result)
	assert.True(t, result != 0)
}

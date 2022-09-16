package client

import (
	"context"

	"github.com/Venachain/client-sdk-go/venachain/keystore"
	precompile "github.com/Venachain/client-sdk-go/precompiled"
)

type RoleClient struct {
	ContractClient
}

func NewRoleClient(ctx context.Context, url URL, keyfilePath string, passphrase string) (*RoleClient, error) {
	client, err := NewContractClient(ctx, url, keyfilePath, passphrase, precompile.UserManagementAddress, "wasm")
	if err != nil {
		return nil, err
	}
	roleClient := &RoleClient{
		*client,
	}
	return roleClient, nil
}

// 传入key 构造Role客户端
func NewRoleClientWithKey(ctx context.Context, url URL, key *keystore.Key) (*RoleClient, error) {
	client, err := NewContractClientWithKey(ctx, url, key, precompile.UserManagementAddress, "wasm")
	if err != nil {
		return nil, err
	}
	roleClient := &RoleClient{
		*client,
	}
	return roleClient, nil
}

// 设置当前client 对应的账户为超级管理员
func (roleClient RoleClient) SetSuperAdmin(ctx context.Context) (string, error) {
	funcName := "setSuperAdmin"

	result, err := roleClient.contractCallWithParams(ctx, nil, funcName, precompile.UserManagementAddress)
	if err != nil {
		return "", err
	}
	res := result.([]interface{})
	return res[0].(string), nil
}

// <address>: 转移后的超级管理员地址
func (roleClient RoleClient) TransferSuperAdmin(ctx context.Context, address string) (string, error) {
	funcName := "transferSuperAdminByAddress"
	funcParams := []string{address}

	result, err := roleClient.contractCallWithParams(ctx, funcParams, funcName, precompile.UserManagementAddress)
	if err != nil {
		return "", err
	}
	res := result.([]interface{})
	return res[0].(string), nil
}

// 为某个账户地址添加指定角色的权限。<address>: 被赋予角色权限的账户地址
func (roleClient RoleClient) AddChainAdmin(ctx context.Context, address string) (string, error) {
	funcName := "addChainAdminByAddress"
	funcParams := []string{address}

	result, err := roleClient.contractCallWithParams(ctx, funcParams, funcName, precompile.UserManagementAddress)
	if err != nil {
		return "", err
	}
	res := result.([]interface{})
	return res[0].(string), nil
}

// 为某个账户地址删除指定角色的权限。<address>: 被赋予角色权限的账户地址
func (roleClient RoleClient) DelChainAdmin(ctx context.Context, address string) (string, error) {
	funcName := "delChainAdminByAddress"
	funcParams := []string{address}

	result, err := roleClient.contractCallWithParams(ctx, funcParams, funcName, precompile.UserManagementAddress)
	if err != nil {
		return "", err
	}
	res := result.([]interface{})
	return res[0].(string), nil
}

func (roleClient RoleClient) AddGroupAdmin(ctx context.Context, address string) (string, error) {
	funcName := "addGroupAdminByAddress"
	funcParams := []string{address}

	result, err := roleClient.contractCallWithParams(ctx, funcParams, funcName, precompile.UserManagementAddress)
	if err != nil {
		return "", err
	}
	res := result.([]interface{})
	return res[0].(string), nil
}

func (roleClient RoleClient) DelGroupAdmin(ctx context.Context, address string) (string, error) {
	funcName := "delGroupAdminByAddress"
	funcParams := []string{address}

	result, err := roleClient.contractCallWithParams(ctx, funcParams, funcName, precompile.UserManagementAddress)
	if err != nil {
		return "", err
	}
	res := result.([]interface{})
	return res[0].(string), nil
}

func (roleClient RoleClient) AddNodeAdmin(ctx context.Context, address string) (string, error) {
	funcName := "addNodeAdminByAddress"
	funcParams := []string{address}

	result, err := roleClient.contractCallWithParams(ctx, funcParams, funcName, precompile.UserManagementAddress)
	if err != nil {
		return "", err
	}
	res := result.([]interface{})
	return res[0].(string), nil
}

func (roleClient RoleClient) DelNodeAdmin(ctx context.Context, address string) (string, error) {
	funcName := "delNodeAdminByAddress"
	funcParams := []string{address}

	result, err := roleClient.contractCallWithParams(ctx, funcParams, funcName, precompile.UserManagementAddress)
	if err != nil {
		return "", err
	}
	res := result.([]interface{})
	return res[0].(string), nil
}

func (roleClient RoleClient) AddContractAdmin(ctx context.Context, address string) (string, error) {
	funcName := "addContractAdminByAddress"
	funcParams := []string{address}

	result, err := roleClient.contractCallWithParams(ctx, funcParams, funcName, precompile.UserManagementAddress)
	if err != nil {
		return "", err
	}
	res := result.([]interface{})
	return res[0].(string), nil
}

func (roleClient RoleClient) DelContractAdmin(ctx context.Context, address string) (string, error) {
	funcName := "delContractAdminByAddress"
	funcParams := []string{address}

	result, err := roleClient.contractCallWithParams(ctx, funcParams, funcName, precompile.UserManagementAddress)
	if err != nil {
		return "", err
	}
	res := result.([]interface{})
	return res[0].(string), nil
}

func (roleClient RoleClient) AddContractDeployer(ctx context.Context, address string) (string, error) {
	funcName := "addContractDeployerByAddress"
	funcParams := []string{address}

	result, err := roleClient.contractCallWithParams(ctx, funcParams, funcName, precompile.UserManagementAddress)
	if err != nil {
		return "", err
	}
	res := result.([]interface{})
	return res[0].(string), nil
}

func (roleClient RoleClient) DelContractDeployer(ctx context.Context, address string) (string, error) {
	funcName := "delContractDeployerByAddress"
	funcParams := []string{address}

	result, err := roleClient.contractCallWithParams(ctx, funcParams, funcName, precompile.UserManagementAddress)
	if err != nil {
		return "", err
	}
	res := result.([]interface{})
	return res[0].(string), nil
}

// <role>: 角色可以且只能为"SUPER_ADMIN", "CHAIN_ADMIN", "GROUP_ADMIN", "NODE_ADMIN", "CONTRACT_ADMIN" ， "CONTRACT_DEPLOYER"其中之一
func (roleClient RoleClient) GetAddrListOfRole(ctx context.Context, role string) (string, error) {
	funcName := "getAddrListOfRole"
	funcParams := []string{role}

	result, err := roleClient.contractCallWithParams(ctx, funcParams, funcName, precompile.UserManagementAddress)
	if err != nil {
		return "", err
	}
	res := result.([]interface{})
	return res[0].(string), nil
}

// 获取某账户地址用户权限。
func (roleClient RoleClient) GetRoles(ctx context.Context, address string) (string, error) {
	funcName := "getRolesByAddress"
	funcParams := []string{address}

	result, err := roleClient.contractCallWithParams(ctx, funcParams, funcName, precompile.UserManagementAddress)
	if err != nil {
		return "", err
	}
	res := result.([]interface{})
	return res[0].(string), nil
}

// 检查某账户地址是否拥有指定用户权限。
// <role>: 角色可以且只能为"SUPER_ADMIN", "CHAIN_ADMIN", "GROUP_ADMIN", "NODE_ADMIN", "CONTRACT_ADMIN" ，"CONTRACT_DEPLOYER"其中之一
//返回结果 有权限 result= 1，无权限 result = 0
func (roleClient RoleClient) HasRole(ctx context.Context, address, role string) (int32, error) {
	funcName := "hasRole"
	funcParams := []string{address, role}

	result, err := roleClient.contractCallWithParams(ctx, funcParams, funcName, precompile.UserManagementAddress)
	if err != nil {
		return 0, err
	}
	res := result.([]interface{})
	return res[0].(int32), nil
}

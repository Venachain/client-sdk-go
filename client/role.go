package client

import (
	common_sdk "github.com/PlatONE_Network/PlatONE-SDK-Go/common"
	precompile "github.com/PlatONE_Network/PlatONE-SDK-Go/precompiled"
)

type RoleClient struct {
	ContractClient
}

func (roleClient RoleClient) SetSuperAdmin(txparam common_sdk.TxParams) (string, error) {
	funcName := "setSuperAdmin"

	result, err := roleClient.contractCallWrap(txparam, nil, funcName, precompile.UserManagementAddress)
	if err != nil {
		return "", err
	}
	res := result[0].([]interface{})
	return res[0].(string), nil
}

// <address>: 转移后的超级管理员地址
func (roleClient RoleClient) TransferSuperAdmin(txparam common_sdk.TxParams, address string) (string, error) {
	funcName := "transferSuperAdminByAddress"
	funcParams := []string{address}

	result, err := roleClient.contractCallWrap(txparam, funcParams, funcName, precompile.UserManagementAddress)
	if err != nil {
		return "", err
	}
	res := result[0].([]interface{})
	return res[0].(string), nil
}

// 为某个账户地址添加指定角色的权限。<address>: 被赋予角色权限的账户地址
func (roleClient RoleClient) AddChainAdmin(txparam common_sdk.TxParams, address string) (string, error) {
	funcName := "addChainAdminByAddress"
	funcParams := []string{address}

	result, err := roleClient.contractCallWrap(txparam, funcParams, funcName, precompile.UserManagementAddress)
	if err != nil {
		return "", err
	}
	res := result[0].([]interface{})
	return res[0].(string), nil
}

// 为某个账户地址删除指定角色的权限。<address>: 被赋予角色权限的账户地址
func (roleClient RoleClient) DelChainAdmin(txparam common_sdk.TxParams, address string) (string, error) {
	funcName := "delChainAdminByAddress"
	funcParams := []string{address}

	result, err := roleClient.contractCallWrap(txparam, funcParams, funcName, precompile.UserManagementAddress)
	if err != nil {
		return "", err
	}
	res := result[0].([]interface{})
	return res[0].(string), nil
}

func (roleClient RoleClient) AddGroupAdmin(txparam common_sdk.TxParams, address string) (string, error) {
	funcName := "addGroupAdminByAddress"
	funcParams := []string{address}

	result, err := roleClient.contractCallWrap(txparam, funcParams, funcName, precompile.UserManagementAddress)
	if err != nil {
		return "", err
	}
	res := result[0].([]interface{})
	return res[0].(string), nil
}

func (roleClient RoleClient) DelGroupAdmin(txparam common_sdk.TxParams, address string) (string, error) {
	funcName := "delGroupAdminByAddress"
	funcParams := []string{address}

	result, err := roleClient.contractCallWrap(txparam, funcParams, funcName, precompile.UserManagementAddress)
	if err != nil {
		return "", err
	}
	res := result[0].([]interface{})
	return res[0].(string), nil
}

func (roleClient RoleClient) AddNodeAdmin(txparam common_sdk.TxParams, address string) (string, error) {
	funcName := "addNodeAdminByAddress"
	funcParams := []string{address}

	result, err := roleClient.contractCallWrap(txparam, funcParams, funcName, precompile.UserManagementAddress)
	if err != nil {
		return "", err
	}
	res := result[0].([]interface{})
	return res[0].(string), nil
}

func (roleClient RoleClient) DelNodeAdmin(txparam common_sdk.TxParams, address string) (string, error) {
	funcName := "delNodeAdminByAddress"
	funcParams := []string{address}

	result, err := roleClient.contractCallWrap(txparam, funcParams, funcName, precompile.UserManagementAddress)
	if err != nil {
		return "", err
	}
	res := result[0].([]interface{})
	return res[0].(string), nil
}

func (roleClient RoleClient) AddContractAdmin(txparam common_sdk.TxParams, address string) (string, error) {
	funcName := "addContractAdminByAddress"
	funcParams := []string{address}

	result, err := roleClient.contractCallWrap(txparam, funcParams, funcName, precompile.UserManagementAddress)
	if err != nil {
		return "", err
	}
	res := result[0].([]interface{})
	return res[0].(string), nil
}

func (roleClient RoleClient) DelContractAdmin(txparam common_sdk.TxParams, address string) (string, error) {
	funcName := "delContractAdminByAddress"
	funcParams := []string{address}

	result, err := roleClient.contractCallWrap(txparam, funcParams, funcName, precompile.UserManagementAddress)
	if err != nil {
		return "", err
	}
	res := result[0].([]interface{})
	return res[0].(string), nil
}

func (roleClient RoleClient) AddContractDeployer(txparam common_sdk.TxParams, address string) (string, error) {
	funcName := "addContractDeployerByAddress"
	funcParams := []string{address}

	result, err := roleClient.contractCallWrap(txparam, funcParams, funcName, precompile.UserManagementAddress)
	if err != nil {
		return "", err
	}
	res := result[0].([]interface{})
	return res[0].(string), nil
}

func (roleClient RoleClient) DelContractDeployer(txparam common_sdk.TxParams, address string) (string, error) {
	funcName := "delContractDeployerByAddress"
	funcParams := []string{address}

	result, err := roleClient.contractCallWrap(txparam, funcParams, funcName, precompile.UserManagementAddress)
	if err != nil {
		return "", err
	}
	res := result[0].([]interface{})
	return res[0].(string), nil
}

// <role>: 角色可以且只能为"SUPER_ADMIN", "CHAIN_ADMIN", "GROUP_ADMIN", "NODE_ADMIN", "CONTRACT_ADMIN" ， "CONTRACT_DEPLOYER"其中之一
func (roleClient RoleClient) GetAddrListOfRole(txparam common_sdk.TxParams, role string) (string, error) {
	funcName := "getAddrListOfRole"
	funcParams := []string{role}

	result, err := roleClient.contractCallWrap(txparam, funcParams, funcName, precompile.UserManagementAddress)
	if err != nil {
		return "", err
	}
	res := result[0].([]interface{})
	return res[0].(string), nil
}

// 获取某账户地址用户权限。
func (roleClient RoleClient) GetRoles(txparam common_sdk.TxParams, address string) (string, error) {
	funcName := "getRolesByAddress"
	funcParams := []string{address}

	result, err := roleClient.contractCallWrap(txparam, funcParams, funcName, precompile.UserManagementAddress)
	if err != nil {
		return "", err
	}
	res := result[0].([]interface{})
	return res[0].(string), nil
}

// 检查某账户地址是否拥有指定用户权限。
// <role>: 角色可以且只能为"SUPER_ADMIN", "CHAIN_ADMIN", "GROUP_ADMIN", "NODE_ADMIN", "CONTRACT_ADMIN" ，"CONTRACT_DEPLOYER"其中之一
//返回结果 有权限 result= 1，无权限 result = 0
func (roleClient RoleClient) HasRole(txparam common_sdk.TxParams, address, role string) (int32, error) {
	funcName := "hasRole"
	funcParams := []string{address, role}

	result, err := roleClient.contractCallWrap(txparam, funcParams, funcName, precompile.UserManagementAddress)
	if err != nil {
		return 0, err
	}
	res := result[0].([]interface{})
	return res[0].(int32), nil
}

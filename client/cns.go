package client

import (
	common_sdk "github.com/PlatONE_Network/PlatONE-SDK-Go/common"
	precompile "github.com/PlatONE_Network/PlatONE-SDK-Go/precompiled"
)

type CnsClient struct {
	ContractClient
	name string
}

// 通过合约名称调用合约（cns服务）

// 将合约注册到cns平台中，注册后的合约不仅可以通过合约账户地址进行调用执行，还可以通过其对应的合约名称进行执行。
func (cnsClient CnsClient) CnsRegister(txparam common_sdk.TxParams, version, address string) (string, error) {
	funcName := "cnsRegister"
	funcParams := []string{cnsClient.name, version, address}

	result, err := cnsClient.contractCallWrap(txparam, funcParams, funcName, precompile.CnsManagementAddress)
	if err != nil {
		return "", err
	}
	res := result[0].([]interface{})
	return res[0].(string), nil
}

// 通过合约名称以及版本号（默认为"latest"）解析出对应的账户地址。一个合约名可以对应多个（在注册的）合约地址，
// 通过版本号解析出对应的合约地址，但在cns平台中已注销的合约名对应的版本号无法解析出相应的账户地址。
func (cnsClient CnsClient) CnsResolve(txparam common_sdk.TxParams, version string) (string, error) {
	funcName := "getContractAddress"
	if version == "" {
		version = "latest"
	}
	funcParams := []string{cnsClient.name, version}

	result, err := cnsClient.contractCallWrap(txparam, funcParams, funcName, precompile.CnsManagementAddress)
	if err != nil {
		return "", err
	}
	res := result[0].([]interface{})
	return res[0].(string), nil
}

// 制定cns名称对应的合约版本
func (cnsClient CnsClient) CnsRedirect(txparam common_sdk.TxParams, version string) (string, error) {
	funcName := "cnsRedirect"
	funcParams := []string{cnsClient.name, version}

	result, err := cnsClient.contractCallWrap(txparam, funcParams, funcName, precompile.CnsManagementAddress)
	if err != nil {
		return "", err
	}
	res := result[0].([]interface{})
	return res[0].(string), nil
}

// 显示所有的cns 合约信息
func (cnsClient CnsClient) CnsQueryAll(txparam common_sdk.TxParams) (string, error) {
	funcName := "getRegisteredContracts"
	funcParams := []string{"0", "0"}

	result, err := cnsClient.contractCallWrap(txparam, funcParams, funcName, precompile.CnsManagementAddress)
	if err != nil {
		return "", err
	}
	res := result[0].([]interface{})
	return res[0].(string), nil
}

// 通过cns 名字来查询合约信息
func (cnsClient CnsClient) CnsQueryByName(txparam common_sdk.TxParams) (string, error) {
	funcName := "getRegisteredContractsByName"
	funcParams := []string{cnsClient.name}

	result, err := cnsClient.contractCallWrap(txparam, funcParams, funcName, precompile.CnsManagementAddress)
	if err != nil {
		return "", err
	}
	res := result[0].([]interface{})
	return res[0].(string), nil
}

// 通过cns 名字来查询合约信息, name 可以为空
func (cnsClient CnsClient) CnsQueryByAddress(txparam common_sdk.TxParams, address string) (string, error) {
	funcName := "getRegisteredContractsByAddress"
	funcParams := []string{address}

	result, err := cnsClient.contractCallWrap(txparam, funcParams, funcName, precompile.CnsManagementAddress)
	if err != nil {
		return "", err
	}
	res := result[0].([]interface{})
	return res[0].(string), nil
}

// 通过账户地址来查询合约信息, name 可以为空
func (cnsClient CnsClient) CnsQueryByAccount(txparam common_sdk.TxParams, account string) (string, error) {
	funcName := "getRegisteredContractsByOrigin"
	funcParams := []string{account}

	result, err := cnsClient.contractCallWrap(txparam, funcParams, funcName, precompile.CnsManagementAddress)
	if err != nil {
		return "", err
	}
	res := result[0].([]interface{})
	return res[0].(string), nil
}

// 根据合约地址查询该地址是否注册了cns
func (cnsClient CnsClient) CnsStateByAddress(txparam common_sdk.TxParams, address string) (int32, error) {
	funcName := "ifRegisteredByAddress"
	funcParams := []string{address}

	result, err := cnsClient.contractCallWrap(txparam, funcParams, funcName, precompile.CnsManagementAddress)
	if err != nil {
		return 0, err
	}
	res := result[0].([]interface{})
	return res[0].(int32), nil
}

// 查询该cns 客户端的名称是否注册了cns
func (cnsClient CnsClient) CnsState(txparam common_sdk.TxParams) (int32, error) {
	funcName := "ifRegisteredByName"
	funcParams := []string{cnsClient.name}

	result, err := cnsClient.contractCallWrap(txparam, funcParams, funcName, precompile.CnsManagementAddress)
	if err != nil {
		return 0, err
	}
	res := result[0].([]interface{})
	return res[0].(int32), nil
}

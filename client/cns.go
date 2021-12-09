package client

import (
	"context"

	common_sdk "git-c.i.wxblockchain.com/PlatONE/src/node/client-sdk-go/common"
	precompile "git-c.i.wxblockchain.com/PlatONE/src/node/client-sdk-go/precompiled"
)

type CnsClient struct {
	ContractClient
	name string
}

// 将合约注册到cns平台中，注册后的合约不仅可以通过合约账户地址进行调用执行，还可以通过其对应的合约名称进行执行。
func (cnsClient CnsClient) CnsRegister(ctx context.Context, txparam common_sdk.TxParams, version, address string) (string, error) {
	funcName := "cnsRegister"
	funcParams := []string{cnsClient.name, version, address}

	result, err := cnsClient.contractCallWrap(ctx, txparam, funcParams, funcName, precompile.CnsManagementAddress)
	if err != nil {
		return "", err
	}
	res := result.([]interface{})
	return res[0].(string), nil
}

// 通过cns 名字调用合约
func (cnsClient CnsClient) CnsExecute(ctx context.Context, txparam common_sdk.TxParams, funcName string, funcParams []string, cns string) (interface{}, error) {
	//var res []interface{}
	isListMethods, err := cnsClient.IsFuncNameInContract(funcName)
	if !isListMethods {
		return nil, err
	}
	funcName, funcParams = common_sdk.FuncParse(funcName, funcParams)

	result, err := cnsClient.contractCallWrap(ctx, txparam, funcParams, funcName, cns)
	return result, nil
}

// 通过合约名称以及版本号（默认为"latest"）解析出对应的账户地址。一个合约名可以对应多个（在注册的）合约地址，
// 通过版本号解析出对应的合约地址，但在cns平台中已注销的合约名对应的版本号无法解析出相应的账户地址。
func (cnsClient CnsClient) CnsResolve(ctx context.Context, txparam common_sdk.TxParams, version string) (string, error) {
	funcName := "getContractAddress"
	if version == "" {
		version = "latest"
	}
	funcParams := []string{cnsClient.name, version}

	result, err := cnsClient.contractCallWrap(ctx, txparam, funcParams, funcName, precompile.CnsManagementAddress)
	if err != nil {
		return "", err
	}
	res := result.([]interface{})
	return res[0].(string), nil
}

// 制定cns名称对应的合约版本
func (cnsClient CnsClient) CnsRedirect(ctx context.Context, txparam common_sdk.TxParams, version string) (string, error) {
	funcName := "cnsRedirect"
	funcParams := []string{cnsClient.name, version}

	result, err := cnsClient.contractCallWrap(ctx, txparam, funcParams, funcName, precompile.CnsManagementAddress)
	if err != nil {
		return "", err
	}
	res := result.([]interface{})
	return res[0].(string), nil
}

// 显示所有的cns 合约信息
func (cnsClient CnsClient) CnsQueryAll(ctx context.Context, txparam common_sdk.TxParams) (string, error) {
	funcName := "getRegisteredContracts"
	funcParams := []string{"0", "0"}

	result, err := cnsClient.contractCallWrap(ctx, txparam, funcParams, funcName, precompile.CnsManagementAddress)
	if err != nil {
		return "", err
	}
	res := result.([]interface{})
	return res[0].(string), nil
}

// 通过cns 名字来查询合约信息
func (cnsClient CnsClient) CnsQueryByName(ctx context.Context, txparam common_sdk.TxParams) (string, error) {
	funcName := "getRegisteredContractsByName"
	funcParams := []string{cnsClient.name}

	result, err := cnsClient.contractCallWrap(ctx, txparam, funcParams, funcName, precompile.CnsManagementAddress)
	if err != nil {
		return "", err
	}
	res := result.([]interface{})
	return res[0].(string), nil
}

// 通过cns 名字来查询合约信息, name 可以为空
func (cnsClient CnsClient) CnsQueryByAddress(ctx context.Context, txparam common_sdk.TxParams, address string) (string, error) {
	funcName := "getRegisteredContractsByAddress"
	funcParams := []string{address}

	result, err := cnsClient.contractCallWrap(ctx, txparam, funcParams, funcName, precompile.CnsManagementAddress)
	if err != nil {
		return "", err
	}
	res := result.([]interface{})
	return res[0].(string), nil
}

// 通过账户地址来查询合约信息, name 可以为空
func (cnsClient CnsClient) CnsQueryByAccount(ctx context.Context, txparam common_sdk.TxParams, account string) (string, error) {
	funcName := "getRegisteredContractsByOrigin"
	funcParams := []string{account}

	result, err := cnsClient.contractCallWrap(ctx, txparam, funcParams, funcName, precompile.CnsManagementAddress)
	if err != nil {
		return "", err
	}
	res := result.([]interface{})
	return res[0].(string), nil
}

// 根据合约地址查询该地址是否注册了cns
func (cnsClient CnsClient) CnsStateByAddress(ctx context.Context, txparam common_sdk.TxParams, address string) (int32, error) {
	funcName := "ifRegisteredByAddress"
	funcParams := []string{address}

	result, err := cnsClient.contractCallWrap(ctx, txparam, funcParams, funcName, precompile.CnsManagementAddress)
	if err != nil {
		return 0, err
	}
	res := result.([]interface{})
	return res[0].(int32), nil
}

// 查询该cns 客户端的名称是否注册了cns
func (cnsClient CnsClient) CnsState(ctx context.Context, txparam common_sdk.TxParams) (int32, error) {
	funcName := "ifRegisteredByName"
	funcParams := []string{cnsClient.name}

	result, err := cnsClient.contractCallWrap(ctx, txparam, funcParams, funcName, precompile.CnsManagementAddress)
	if err != nil {
		return 0, err
	}
	res := result.([]interface{})
	return res[0].(int32), nil
}

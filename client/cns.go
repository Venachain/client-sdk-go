package client

import (
	"context"

	"github.com/Venachain/client-sdk-go/packet"
	"github.com/Venachain/client-sdk-go/venachain/keystore"
	precompile "github.com/Venachain/client-sdk-go/precompiled"
)

type CnsClient struct {
	ContractClient
	name string
}

func NewCnsClient(ctx context.Context, url URL, keyfilePath string, passphrase string, name string) (*CnsClient, error) {
	client, err := NewContractClient(ctx, url, keyfilePath, passphrase, precompile.CnsManagementAddress, "wasm")
	if err != nil {
		return nil, err
	}

	cnsClient := &CnsClient{
		*client,
		name,
	}
	return cnsClient, nil
}

// 传入key 构造Cns客户端
func NewCnsClientWithKey(ctx context.Context, url URL, key *keystore.Key, name string) (*CnsClient, error) {
	client, err := NewContractClientWithKey(ctx, url, key, precompile.CnsManagementAddress, "wasm")
	if err != nil {
		return nil, err
	}
	cnsClient := &CnsClient{
		*client,
		name,
	}
	return cnsClient, nil
}

// 将合约注册到cns平台中，注册后的合约不仅可以通过合约账户地址进行调用执行，还可以通过其对应的合约名称进行执行。
func (cnsClient CnsClient) CnsRegister(ctx context.Context, version string, ContractAddress string) (string, error) {
	funcName := "cnsRegister"
	funcParams := []string{cnsClient.name, version, ContractAddress}

	result, err := cnsClient.contractCallWithParams(ctx, funcParams, funcName, precompile.CnsManagementAddress)
	if err != nil {
		return "", err
	}
	res := result.([]interface{})
	return res[0].(string), nil
}

// 通过cns 名字调用合约
func (cnsClient CnsClient) CnsExecute(ctx context.Context, funcName string, funcParams []string, cns string) (interface{}, error) {
	isListMethods, err := cnsClient.IsFuncNameInContract(funcName)
	if !isListMethods {
		return nil, err
	}
	funcName, funcParams = packet.FuncParse(funcName, funcParams)

	result, err := cnsClient.contractCallWithParams(ctx, funcParams, funcName, cns)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// 通过合约名称以及版本号（默认为"latest"）解析出对应的账户地址。一个合约名可以对应多个（在注册的）合约地址，
// 通过版本号解析出对应的合约地址，但在cns平台中已注销的合约名对应的版本号无法解析出相应的账户地址。
func (cnsClient CnsClient) CnsResolve(ctx context.Context, version string) (string, error) {
	funcName := "getContractAddress"
	if version == "" {
		version = "latest"
	}
	funcParams := []string{cnsClient.name, version}

	result, err := cnsClient.contractCallWithParams(ctx, funcParams, funcName, precompile.CnsManagementAddress)
	if err != nil {
		return "", err
	}
	res := result.([]interface{})
	return res[0].(string), nil
}

// 制定cns名称对应的合约版本
func (cnsClient CnsClient) CnsRedirect(ctx context.Context, version string) (string, error) {
	funcName := "cnsRedirect"
	funcParams := []string{cnsClient.name, version}

	result, err := cnsClient.contractCallWithParams(ctx, funcParams, funcName, precompile.CnsManagementAddress)
	if err != nil {
		return "", err
	}
	res := result.([]interface{})
	return res[0].(string), nil
}

// 显示所有的cns 合约信息
func (cnsClient CnsClient) CnsQueryAll(ctx context.Context) (string, error) {
	funcName := "getRegisteredContracts"
	funcParams := []string{"0", "0"}

	result, err := cnsClient.contractCallWithParams(ctx, funcParams, funcName, precompile.CnsManagementAddress)
	if err != nil {
		return "", err
	}
	res := result.([]interface{})
	return res[0].(string), nil
}

// 通过cns 名字来查询合约信息
func (cnsClient CnsClient) CnsQueryByName(ctx context.Context) (string, error) {
	funcName := "getRegisteredContractsByName"
	funcParams := []string{cnsClient.name}

	result, err := cnsClient.contractCallWithParams(ctx, funcParams, funcName, precompile.CnsManagementAddress)
	if err != nil {
		return "", err
	}
	res := result.([]interface{})
	return res[0].(string), nil
}

// 通过cns 名字来查询合约信息, name 可以为空
func (cnsClient CnsClient) CnsQueryByAddress(ctx context.Context, ContractAddress string) (string, error) {
	funcName := "getRegisteredContractsByAddress"
	funcParams := []string{ContractAddress}

	result, err := cnsClient.contractCallWithParams(ctx, funcParams, funcName, precompile.CnsManagementAddress)
	if err != nil {
		return "", err
	}
	res := result.([]interface{})
	return res[0].(string), nil
}

// 通过账户地址来查询合约信息, name 可以为空
func (cnsClient CnsClient) CnsQueryByAccount(ctx context.Context, account string) (string, error) {
	funcName := "getRegisteredContractsByOrigin"
	funcParams := []string{account}

	result, err := cnsClient.contractCallWithParams(ctx, funcParams, funcName, precompile.CnsManagementAddress)
	if err != nil {
		return "", err
	}
	res := result.([]interface{})
	return res[0].(string), nil
}

// 根据合约地址查询该地址是否注册了cns
func (cnsClient CnsClient) CnsStateByAddress(ctx context.Context, address string) (int32, error) {
	funcName := "ifRegisteredByAddress"
	funcParams := []string{address}

	result, err := cnsClient.contractCallWithParams(ctx, funcParams, funcName, precompile.CnsManagementAddress)
	if err != nil {
		return 0, err
	}
	res := result.([]interface{})
	return res[0].(int32), nil
}

// 查询该cns 客户端的名称是否注册了cns
func (cnsClient CnsClient) CnsState(ctx context.Context) (int32, error) {
	funcName := "ifRegisteredByName"
	funcParams := []string{cnsClient.name}

	result, err := cnsClient.contractCallWithParams(ctx, funcParams, funcName, precompile.CnsManagementAddress)
	if err != nil {
		return 0, err
	}
	res := result.([]interface{})
	return res[0].(int32), nil
}

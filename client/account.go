package client

import (
	"encoding/json"

	"golang.org/x/net/context"

	common_sdk "git-c.i.wxblockchain.com/PlatONE/src/node/client-sdk-go/common"
	"git-c.i.wxblockchain.com/PlatONE/src/node/client-sdk-go/platone/common"
	precompile "git-c.i.wxblockchain.com/PlatONE/src/node/client-sdk-go/precompiled"
	"git-c.i.wxblockchain.com/PlatONE/src/node/client-sdk-go/precompiled/syscontracts"
)

type AccountClient struct {
	ContractClient
	Address common.Address
}

func (accountClient AccountClient) UserAdd(ctx context.Context, txparam common_sdk.TxParams, name, phone, email, organization string) (string, error) {
	userdescinfo := syscontracts.UserDescInfo{}
	var userinfo syscontracts.UserInfo
	funcName := "addUser"
	userdescinfo.Phone = phone
	userdescinfo.Email = email
	userdescinfo.Organization = organization
	desbytes, _ := json.Marshal(userdescinfo)
	desinfo := string(desbytes)

	userinfo.Address = accountClient.Address
	userinfo.Name = name
	userinfo.DescInfo = desinfo
	bytes, _ := json.Marshal(userinfo)
	strJson := string(bytes)
	funcParams := []string{strJson}

	result, err := accountClient.contractCallWrap(ctx, txparam, funcParams, funcName, precompile.UserManagementAddress)
	if err != nil {
		return "", err
	}
	res := result.([]interface{})
	return res[0].(string), nil
}

func (accountClient AccountClient) UserUpdate(ctx context.Context, txparam common_sdk.TxParams, phone, email, organization string) (string, error) {
	var funcParams []string
	funcParams = append(funcParams, accountClient.Address.Hex())
	funcName := "updateUserDescInfo"
	userdescinfo := syscontracts.UserDescInfo{}
	userdescinfo.Phone = phone
	userdescinfo.Email = email
	userdescinfo.Organization = organization
	desbytes, _ := json.Marshal(userdescinfo)
	funcParams = append(funcParams, string(desbytes))

	result, err := accountClient.contractCallWrap(ctx, txparam, funcParams, funcName, precompile.UserManagementAddress)
	if err != nil {
		return "", err
	}
	res := result.([]interface{})
	return res[0].(string), nil
}

// 根据账户地址或者账户名称查询用户信息，如果传入的name == ""，则为查找所有账户信息
func (accountClient AccountClient) QueryUser(ctx context.Context, txparam common_sdk.TxParams, user string) (string, error) {
	var funcName string
	var funcParams = make([]string, 0)
	if user != "" {
		isAddress, err := common_sdk.ParamParse(user, "user")
		if err != nil {
			return "", err
		}
		isAddress = isAddress.(int32)
		if isAddress == common_sdk.CnsIsAddress {
			funcName = "getUserByAddress"
		} else {
			funcName = "getUserByName"
		}
		funcParams = []string{user}
	} else {
		funcName = "getAllUsers"
	}

	result, err := accountClient.contractCallWrap(ctx, txparam, funcParams, funcName, precompile.UserManagementAddress)
	if err != nil {
		return "", err
	}
	res := result.([]interface{})
	return res[0].(string), nil
}

func (accountClient AccountClient) Lock(ctx context.Context) (bool, error) {
	funcName := "personal_lockAccount"
	funcParams := accountClient.Address.Hex()
	var res bool

	//result, err := accountClient.Client.RPCSend(ctx, funcName, funcParams)
	result, err := accountClient.RpcClient.CallContext(ctx, funcName, funcParams)

	if err != nil {
		return false, err
	}
	if err = json.Unmarshal(result, &res); err != nil {
		return false, err
	}
	return res, nil
}

func (accountClient AccountClient) UnLock(ctx context.Context) (bool, error) {
	funcName := "personal_unlockAccount"
	account := accountClient.Address.Hex()

	var res bool
	result, err := accountClient.RpcClient.CallContext(ctx, funcName, account,accountClient.Passphrase,0)
	//result, err := accountClient.Client.RPCSend(ctx, funcName, account,accountClient.Passphrase,0)
	if err != nil {
		return false, err
	}
	if err = json.Unmarshal(result, &res); err != nil {
		return false, err
	}
	return res, nil
}
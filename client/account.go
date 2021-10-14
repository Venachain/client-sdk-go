package client

import (
	"encoding/json"
	"time"

	"golang.org/x/net/context"

	common_sdk "github.com/PlatONE_Network/PlatONE-SDK-Go/common"
	"github.com/PlatONE_Network/PlatONE-SDK-Go/platone/common"
	precompile "github.com/PlatONE_Network/PlatONE-SDK-Go/precompiled"
	"github.com/PlatONE_Network/PlatONE-SDK-Go/precompiled/syscontracts"
)

type AccountClient struct {
	ContractClient
	Address common.Address
}

func (accountClient AccountClient) UserAdd(txparam common_sdk.TxParams, name, phone, email, organization string) (string, error) {
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

	result, err := accountClient.contractCallWrap(txparam, funcParams, funcName, precompile.UserManagementAddress)
	if err != nil {
		return "", err
	}
	res := result[0].([]interface{})
	return res[0].(string), nil
}

func (accountClient AccountClient) UserUpdate(txparam common_sdk.TxParams, phone, email, organization string) (string, error) {
	var funcParams []string
	funcParams = append(funcParams, accountClient.Address.Hex())
	funcName := "updateUserDescInfo"
	userdescinfo := syscontracts.UserDescInfo{}
	userdescinfo.Phone = phone
	userdescinfo.Email = email
	userdescinfo.Organization = organization
	desbytes, _ := json.Marshal(userdescinfo)
	funcParams = append(funcParams, string(desbytes))

	result, err := accountClient.contractCallWrap(txparam, funcParams, funcName, precompile.UserManagementAddress)
	if err != nil {
		return "", err
	}
	res := result[0].([]interface{})
	return res[0].(string), nil
}

// 根据账户地址或者账户名称查询用户信息，如果传入的name == ""，则为查找所有账户信息
func (accountClient AccountClient) QueryUser(txparam common_sdk.TxParams, user string) (string, error) {
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

	result, err := accountClient.contractCallWrap(txparam, funcParams, funcName, precompile.UserManagementAddress)
	if err != nil {
		return "", err
	}
	res := result[0].([]interface{})
	return res[0].(string), nil
}

func (accountClient AccountClient) Lock(txparam common_sdk.TxParams) (bool, error) {
	//var funcParams = make([]string, 0)
	funcName := "personal_lockAccount"
	funcParams := accountClient.Address.Hex()
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	var res bool
	err := accountClient.Client.RpcClient.CallContext(ctx, &res, funcName, funcParams)
	if err != nil {
		return false, err
	}
	return res, nil
}

func (accountClient AccountClient) UnLock(txparam common_sdk.TxParams) (bool, error) {
	//var funcParams = make([]string, 0)
	funcName := "personal_unlockAccount"
	funcParams := accountClient.Address.Hex()
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	var res bool
	err := accountClient.Client.RpcClient.CallContext(ctx, &res, funcName, funcParams)
	if err != nil {
		return false, err
	}
	return res, nil
}

//func (accountClient AccountClient) Transfer(txparam common_sdk.TxParams, to string) (string, error) {
//	addr := common.HexToAddress(to)
//	accountClient.Client.clientCommonV2(txparam, nil, addr, true)
//
//	result, err := accountClient.contractCallWrap(txparam, funcParams, "updateUserDescInfo", precompile.UserManagementAddress)
//	res := result[0].([]interface{})
//	return res[0].(string), err
//}

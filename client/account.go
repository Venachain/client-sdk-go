package client

import (
	"encoding/json"

	"git-c.i.wxblockchain.com/vena/src/client-sdk-go/packet"
	precompile "git-c.i.wxblockchain.com/vena/src/client-sdk-go/precompiled"
	"git-c.i.wxblockchain.com/vena/src/client-sdk-go/precompiled/syscontracts"
	common_venachain "git-c.i.wxblockchain.com/vena/src/client-sdk-go/venachain/common"
	"git-c.i.wxblockchain.com/vena/src/client-sdk-go/venachain/keystore"
	"golang.org/x/net/context"
)

type AccountClient struct {
	ContractClient
	Address common_venachain.Address
}

func NewAccountClient(ctx context.Context, url URL, keyfilePath string, passphrase string, address string) (*AccountClient, error) {
	client, err := NewContractClient(ctx, url, keyfilePath, passphrase, precompile.UserManagementAddress, "wasm")
	if err != nil {
		return nil, err
	}
	accountClient := &AccountClient{
		*client,
		common_venachain.HexToAddress(address),
	}
	return accountClient, nil
}

// 传入key 构造账户客户端
func NewAccountClientWithKey(ctx context.Context, url URL, key *keystore.Key, address string) (*AccountClient, error) {
	client, err := NewContractClientWithKey(ctx, url, key, precompile.UserManagementAddress, "wasm")
	if err != nil {
		return nil, err
	}
	accountClient := &AccountClient{
		*client,
		common_venachain.HexToAddress(address),
	}
	return accountClient, nil
}

func (accountClient AccountClient) UserAdd(ctx context.Context, name, phone, email, organization string) (string, error) {
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

	result, err := accountClient.contractCallWithParams(ctx, funcParams, funcName, precompile.UserManagementAddress)
	if err != nil {
		return "", err
	}
	res := result.([]interface{})
	return res[0].(string), nil
}

func (accountClient AccountClient) UserUpdate(ctx context.Context, phone, email, organization string) (string, error) {
	var funcParams []string
	funcParams = append(funcParams, accountClient.Address.Hex())
	funcName := "updateUserDescInfo"
	userdescinfo := syscontracts.UserDescInfo{}
	userdescinfo.Phone = phone
	userdescinfo.Email = email
	userdescinfo.Organization = organization
	desbytes, _ := json.Marshal(userdescinfo)
	funcParams = append(funcParams, string(desbytes))

	result, err := accountClient.contractCallWithParams(ctx, funcParams, funcName, precompile.UserManagementAddress)
	if err != nil {
		return "", err
	}
	res := result.([]interface{})
	return res[0].(string), nil
}

// 根据账户地址或者账户名称查询用户信息，如果传入的name == ""，则为查找所有账户信息
func (accountClient AccountClient) QueryUser(ctx context.Context, user string) (string, error) {
	var funcName string
	var funcParams = make([]string, 0)
	if user != "" {
		isAddress, err := packet.ParamParse(user, "user")
		if err != nil {
			return "", err
		}
		isAddress = isAddress.(int32)
		if isAddress == packet.CnsIsAddress {
			funcName = "getUserByAddress"
		} else {
			funcName = "getUserByName"
		}
		funcParams = []string{user}
	} else {
		funcName = "getAllUsers"
	}

	result, err := accountClient.contractCallWithParams(ctx, funcParams, funcName, precompile.UserManagementAddress)
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

func (accountClient AccountClient) UnLock(ctx context.Context, passphrase string) (bool, error) {
	funcName := "personal_unlockAccount"
	account := accountClient.Address.Hex()

	var res bool
	result, err := accountClient.RpcClient.CallContext(ctx, funcName, account, passphrase, 0)
	if err != nil {
		return false, err
	}
	if err = json.Unmarshal(result, &res); err != nil {
		return false, err
	}
	return res, nil
}

func (accountClient AccountClient) CreateAccount(ctx context.Context, passphrase string) (*common_venachain.Address, error) {
	funcName := "personal_newAccount"
	params := passphrase
	result, err := accountClient.RpcClient.CallContext(ctx, funcName, params)
	if err != nil {
		return nil, err
	}
	var res common_venachain.Address
	if err = json.Unmarshal(result, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

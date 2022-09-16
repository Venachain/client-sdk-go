package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/Venachain/client-sdk-go/packet"
	precompile "github.com/Venachain/client-sdk-go/precompiled"
	common_venachain "github.com/Venachain/client-sdk-go/venachain/common"
	"github.com/Venachain/client-sdk-go/venachain/keystore"
)

type FireWallClient struct {
	ContractClient
	ContractAddress string
}

type ExportFwStatus struct {
	AcceptedList []FwElem
	RejectedList []FwElem
}

type FwElem struct {
	Addr     common_venachain.Address
	FuncName string
}

func NewFireWallClient(ctx context.Context, url URL, keyfilePath string, passphrase string, contractAddress string) (*FireWallClient, error) {
	address := getHexAddress(contractAddress)
	client, err := NewContractClient(ctx, url, keyfilePath, passphrase, precompile.FirewallManagementAddress, "wasm")
	if err != nil {
		return nil, err
	}
	fireWallClient := &FireWallClient{
		*client,
		address,
	}
	return fireWallClient, nil
}

// 传入key 构造FireWall客户端
func NewFireWallClientWithKey(ctx context.Context, url URL, key *keystore.Key, contractAddress string) (*FireWallClient, error) {
	address := getHexAddress(contractAddress)
	client, err := NewContractClientWithKey(ctx, url, key, precompile.FirewallManagementAddress, "wasm")
	if err != nil {
		return nil, err
	}
	fireWallClient := &FireWallClient{
		*client,
		address,
	}
	return fireWallClient, nil
}

func (firewallClient FireWallClient) FwStatus(ctx context.Context) (string, error) {
	funcName := "__sys_FwStatus"
	funcParams := []string{firewallClient.ContractAddress}
	result, err := firewallClient.contractCallWithParams(ctx, funcParams, funcName, precompile.FirewallManagementAddress)
	if err != nil {
		return "", err
	}
	res := result.([]interface{})
	return res[0].(string), nil
}

func (firewallClient FireWallClient) FwStart(ctx context.Context) (string, error) {
	funcName := "__sys_FwOpen"
	funcParams := []string{firewallClient.ContractAddress}
	result, err := firewallClient.contractCallWithParams(ctx, funcParams, funcName, precompile.FirewallManagementAddress)
	if err != nil {
		return "", err
	}
	res := result.([]interface{})
	return res[0].(string), nil
}

func (firewallClient FireWallClient) FwClose(ctx context.Context) (string, error) {
	funcName := "__sys_FwClose"
	funcParams := []string{firewallClient.ContractAddress}
	result, err := firewallClient.contractCallWithParams(ctx, funcParams, funcName, precompile.FirewallManagementAddress)
	if err != nil {
		return "", err
	}
	res := result.([]interface{})
	return res[0].(string), nil
}

// 将指定合约的防火墙规则导出到指定位置的防火墙规则文件中,ture 为导出成功，false 为导出失败
func (firewallClient FireWallClient) FwExport(ctx context.Context, filePath string) (bool, error) {
	var rule = new(ExportFwStatus)
	funcName := "__sys_FwExport"
	funcParams := []string{firewallClient.ContractAddress}
	result, err := firewallClient.contractCallWithParams(ctx, funcParams, funcName, precompile.FirewallManagementAddress)
	res := result.([]interface{})

	_ = json.Unmarshal([]byte(res[0].(string)), rule)
	ruleBytes, err := json.Marshal(rule)
	if err != nil {
		return false, errors.New("marshal rule is error")
	}
	err = packet.WriteFile(ruleBytes, filePath)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (firewallClient FireWallClient) FwImport(ctx context.Context, filePath string) (string, error) {
	funcName := "__sys_FwImport"
	fileBytes, err := packet.ParseFileToBytes(filePath)
	if err != nil {
		return "", err
	}
	funcParams := []string{firewallClient.ContractAddress, string(fileBytes)}

	result, err := firewallClient.contractCallWithParams(ctx, funcParams, funcName, precompile.FirewallManagementAddress)
	if err != nil {
		return "", err
	}
	res := result.([]interface{})
	return res[0].(string), nil
}

// 新建一条或多条指定合约的防火墙规则。一条防火墙规则包含具体的防火墙操作（accept或reject操作），需要进行过滤的账户地址以及需要进行限制访问的合约接口名。
//<action>: 防火墙操作：允许accept或拒绝reject
//<account>: 指定被过滤的一个或多个用户账户地址，'*'表示防火墙规则对所有用户账户地址生效。格式["<address1>","<address2>"]，单个账户地址可省略[]。
//<api>: 指定过滤的合约接口名。格式["<funcname1>","<funcname2>"]，单个接口名可省略[]。示例--api "getName"
func (firewallClient FireWallClient) FwNew(ctx context.Context, action, targetAddr, api string) (string, error) {
	funcName := "__sys_FwAdd"
	return firewallClient.fwCommon(ctx, action, targetAddr, api, funcName)
}

// 删除一条防火墙规则
func (firewallClient FireWallClient) FwDelete(ctx context.Context, action, targetAddr, api string) (string, error) {
	funcName := "__sys_FwDel"
	return firewallClient.fwCommon(ctx, action, targetAddr, api, funcName)
}

// 重置防火墙规则
func (firewallClient FireWallClient) FwReset(ctx context.Context, action, targetAddr, api string) (string, error) {
	funcName := "__sys_FwSet"
	return firewallClient.fwCommon(ctx, action, targetAddr, api, funcName)
}

func (firewallClient FireWallClient) fwCommon(ctx context.Context, action, targetAddr, api, funcName string) (string, error) {
	packet.ParamValid(action, "action")
	packet.ParamValid(targetAddr, "fw")
	packet.ParamValid(api, "name")
	rules := packet.CombineRule(targetAddr, api) //TODO batch rules

	funcParams := packet.CombineFuncParams(firewallClient.ContractAddress, action, rules)
	result, err := firewallClient.contractCallWithParams(ctx, funcParams, funcName, precompile.FirewallManagementAddress)
	if err != nil {
		return "", err
	}
	res := result.([]interface{})
	return res[0].(string), nil
}

// 清空指定合约的防火墙的approve操作或reject操作的全部规则
// 如果action 为""，默认清除所有规则
func (firewallClient FireWallClient) FwClear(ctx context.Context, action string) (string, error) {
	if action == "" {
		result1, err := firewallClient.clearCommon(ctx, "reject")
		result2, err := firewallClient.clearCommon(ctx, "accept")
		if err != nil {
			return "", err
		} else {
			return fmt.Sprintf(result1 + "/d" + result2), nil
		}
	}
	result, err := firewallClient.clearCommon(ctx, action)
	if err != nil {
		return "", err
	}
	return result, nil

}

func (firewallClient FireWallClient) clearCommon(ctx context.Context, action string) (string, error) {
	funcName := "__sys_FwClear"
	packet.ParamValid(action, "action")
	funcParams := []string{firewallClient.ContractAddress, action}

	result, err := firewallClient.contractCallWithParams(ctx, funcParams, funcName, precompile.FirewallManagementAddress)
	if err != nil {
		return "", err
	}
	res := result.([]interface{})
	return res[0].(string), nil
}

func getHexAddress(address string) string {
	if !strings.HasPrefix(address, "0x") {
		address = "0x" + address
	}
	return address
}

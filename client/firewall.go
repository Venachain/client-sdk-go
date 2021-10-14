package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/PlatONE_Network/PlatONE-SDK-Go/packet"

	common_sdk "github.com/PlatONE_Network/PlatONE-SDK-Go/common"
	"github.com/PlatONE_Network/PlatONE-SDK-Go/platone/common"
	precompile "github.com/PlatONE_Network/PlatONE-SDK-Go/precompiled"
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
	Addr     common.Address
	FuncName string
}

func (firewallClient FireWallClient) FwStatus(ctx context.Context, txparam common_sdk.TxParams) (string, error) {
	funcName := "__sys_FwStatus"
	funcParams := []string{firewallClient.ContractAddress}
	result, err := firewallClient.contractCallWrap(ctx, txparam, funcParams, funcName, precompile.FirewallManagementAddress)
	if err != nil {
		return "", err
	}
	res := result[0].([]interface{})
	return res[0].(string), nil
}

func (firewallClient FireWallClient) FwStart(ctx context.Context, txparam common_sdk.TxParams) (string, error) {
	funcName := "__sys_FwOpen"
	funcParams := []string{firewallClient.ContractAddress}
	result, err := firewallClient.contractCallWrap(ctx, txparam, funcParams, funcName, precompile.FirewallManagementAddress)
	if err != nil {
		return "", err
	}
	res := result[0].([]interface{})
	return res[0].(string), nil
}

func (firewallClient FireWallClient) FwClose(ctx context.Context, txparam common_sdk.TxParams) (string, error) {
	funcName := "__sys_FwClose"
	funcParams := []string{firewallClient.ContractAddress}
	result, err := firewallClient.contractCallWrap(ctx, txparam, funcParams, funcName, precompile.FirewallManagementAddress)
	if err != nil {
		return "", err
	}
	res := result[0].([]interface{})
	return res[0].(string), nil
}

// 将指定合约的防火墙规则导出到指定位置的防火墙规则文件中,ture 为导出成功，false 为导出失败
func (firewallClient FireWallClient) FwExport(ctx context.Context, txparam common_sdk.TxParams, filePath string) (bool, error) {
	var rule = new(ExportFwStatus)
	funcName := "__sys_FwExport"
	funcParams := []string{firewallClient.ContractAddress}
	result, err := firewallClient.contractCallWrap(ctx, txparam, funcParams, funcName, precompile.FirewallManagementAddress)
	res := result[0].([]interface{})

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

func (firewallClient FireWallClient) FwImport(ctx context.Context, txparam common_sdk.TxParams, filePath string) (string, error) {
	funcName := "__sys_FwImport"
	fileBytes, err := packet.ParseFileToBytes(filePath)
	if err != nil {
		return "", err
	}
	funcParams := []string{firewallClient.ContractAddress, string(fileBytes)}

	result, err := firewallClient.contractCallWrap(ctx, txparam, funcParams, funcName, precompile.FirewallManagementAddress)
	if err != nil {
		return "", err
	}
	res := result[0].([]interface{})
	return res[0].(string), nil
}

// 新建一条或多条指定合约的防火墙规则。一条防火墙规则包含具体的防火墙操作（accept或reject操作），需要进行过滤的账户地址以及需要进行限制访问的合约接口名。
//<action>: 防火墙操作：允许accept或拒绝reject
//<account>: 指定被过滤的一个或多个用户账户地址，'*'表示防火墙规则对所有用户账户地址生效。格式["<address1>","<address2>"]，单个账户地址可省略[]。
//<api>: 指定过滤的合约接口名。格式["<funcname1>","<funcname2>"]，单个接口名可省略[]。示例--api "getName"
func (firewallClient FireWallClient) FwNew(ctx context.Context, txparam common_sdk.TxParams, action, targetAddr, api string) (string, error) {
	funcName := "__sys_FwAdd"
	return firewallClient.fwCommon(ctx, txparam, action, targetAddr, api, funcName)
}

// 删除一条防火墙规则
func (firewallClient FireWallClient) FwDelete(ctx context.Context, txparam common_sdk.TxParams, action, targetAddr, api string) (string, error) {
	funcName := "__sys_FwDel"
	return firewallClient.fwCommon(ctx, txparam, action, targetAddr, api, funcName)
}

// 重置防火墙规则
func (firewallClient FireWallClient) FwReset(ctx context.Context, txparam common_sdk.TxParams, action, targetAddr, api string) (string, error) {
	funcName := "__sys_FwSet"
	return firewallClient.fwCommon(ctx, txparam, action, targetAddr, api, funcName)
}

func (firewallClient FireWallClient) fwCommon(ctx context.Context, txparam common_sdk.TxParams, action, targetAddr, api, funcName string) (string, error) {
	common_sdk.ParamValid(action, "action")
	common_sdk.ParamValid(targetAddr, "fw")
	common_sdk.ParamValid(api, "name")
	rules := common_sdk.CombineRule(targetAddr, api) //TODO batch rules

	funcParams := common_sdk.CombineFuncParams(firewallClient.ContractAddress, action, rules)
	result, err := firewallClient.contractCallWrap(ctx, txparam, funcParams, funcName, precompile.FirewallManagementAddress)
	if err != nil {
		return "", err
	}
	res := result[0].([]interface{})
	return res[0].(string), nil
}

// 清空指定合约的防火墙的approve操作或reject操作的全部规则
// 如果action 为""，默认清除所有规则
func (firewallClient FireWallClient) FwClear(ctx context.Context, txparam common_sdk.TxParams, action string) (string, error) {
	if action == "" {
		result1, err := firewallClient.clearCommon(ctx, txparam, "reject")
		result2, err := firewallClient.clearCommon(ctx, txparam, "accept")
		if err != nil {
			return "", err
		} else {
			return fmt.Sprintf(result1 + "/d" + result2), nil
		}
	}
	result, err := firewallClient.clearCommon(ctx, txparam, action)
	if err != nil {
		return "", err
	}
	return result, nil

}

func (firewallClient FireWallClient) clearCommon(ctx context.Context, txparam common_sdk.TxParams, action string) (string, error) {
	funcName := "__sys_FwClear"
	common_sdk.ParamValid(action, "action")
	funcParams := []string{firewallClient.ContractAddress, action}

	result, err := firewallClient.contractCallWrap(ctx, txparam, funcParams, funcName, precompile.FirewallManagementAddress)
	if err != nil {
		return "", err
	}
	res := result[0].([]interface{})
	return res[0].(string), nil
}

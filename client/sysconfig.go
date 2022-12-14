package client

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/Venachain/client-sdk-go/log"
	"github.com/Venachain/client-sdk-go/packet"
	"github.com/Venachain/client-sdk-go/venachain/keystore"
	"github.com/Venachain/client-sdk-go/venachain/vm"
	precompile "github.com/Venachain/client-sdk-go/precompiled"
)

const (
	txUseGas    = "use-gas" // IsTxUseGas
	txNotUseGas = "not-use"

	conAudit    = "audit"
	conNotAudit = "not-audit"

	checkPerm    = "with-perm"
	notCheckPerm = "without-perm"

	prodEmp    = "allow-empty"
	notProdEmp = "notallow-empty"
)

type SysConfigClient struct {
	ContractClient
}

type SysConfigParam struct {
	Block_gaslimit                  string
	Tx_gaslimit                     string
	Tx_use_gas                      string
	IsCheckContractDeployPermission string
	IsApproveDeployedContract       string
	Empty_block                     string
	Gas_contract                    string
	VrfParams                       string
}

func NewSysConfigClient(ctx context.Context, url URL, keyfilePath string, passphrase string) (*SysConfigClient, error) {
	client, err := NewContractClient(ctx, url, keyfilePath, passphrase, precompile.ParameterManagementAddress, "wasm")
	if err != nil {
		return nil, err
	}
	sysConfigClient := &SysConfigClient{
		*client,
	}
	return sysConfigClient, nil
}

// 传入key 构造 SysConfig 客户端
func NewSysConfigClientWithKey(ctx context.Context, url URL, key *keystore.Key) (*SysConfigClient, error) {
	client, err := NewContractClientWithKey(ctx, url, key, precompile.ParameterManagementAddress, "wasm")
	if err != nil {
		return nil, err
	}
	sysConfigClient := &SysConfigClient{
		*client,
	}
	return sysConfigClient, nil
}

// 	TxGasLimitMinValue        uint64 = 12771596 * 100 // 12771596 大致相当于 0.012772s
//	TxGasLimitMaxValue        uint64 = 2e9            // 相当于 2s
//	BlockGasLimitMinValue     uint64 = 12771596 * 100 // 12771596 大致相当于 0.012772s
//	BlockGasLimitMaxValue     uint64 = 2e10           // 相当于 20s
//--block-gaslimit : the gas limit of the block
//--tx-gaslimit : the gas limit of transactions
//--tx-use-gas : if transactions use gas, 'use-gas' for transactions use gas, 'not-use' for not
//--audit-con : approve the deployed contracts, 'audit' for allowing contracts audit, 'not-audit' for not
//--check-perm : check the sender permission when deploying contracts, 'with-perm' for checking permission, 'without-perm' for not
//--empty-block : consensus produces empty block, 'allow-empty' for allowing to produce empty block, 'notallow-empty' for not
//--gas-contract : register the gas contract by contract name
func (sysConfigClient SysConfigClient) SetSysConfig(ctx context.Context, request SysConfigParam) ([]string, error) {
	var result []string
	if request.Tx_gaslimit != "" {
		res, err := setConfig(sysConfigClient, ctx, request.Tx_gaslimit, vm.TxGasLimitKey)
		if err != nil {
			return nil, err
		}
		result = append(result, res)
	}
	if request.Block_gaslimit != "" {
		res, err := setConfig(sysConfigClient, ctx, request.Block_gaslimit, vm.BlockGasLimitKey)
		if err != nil {
			return nil, err
		}
		result = append(result, res)
	}
	if request.Tx_use_gas != "" {
		res, err := setConfig(sysConfigClient, ctx, request.Tx_use_gas, vm.IsTxUseGasKey)
		if err != nil {
			return nil, err
		}
		result = append(result, res)
	}
	if request.IsApproveDeployedContract != "" {
		res, err := setConfig(sysConfigClient, ctx, request.IsApproveDeployedContract, vm.IsApproveDeployedContractKey)
		if err != nil {
			return nil, err
		}
		result = append(result, res)
	}
	if request.IsCheckContractDeployPermission != "" {
		res, err := setConfig(sysConfigClient, ctx, request.IsCheckContractDeployPermission, vm.IsCheckContractDeployPermissionKey)
		if err != nil {
			return nil, err
		}
		result = append(result, res)
	}
	if request.Empty_block != "" {
		res, err := setConfig(sysConfigClient, ctx, request.Empty_block, vm.IsProduceEmptyBlockKey)
		if err != nil {
			return nil, err
		}
		result = append(result, res)
	}
	if request.Gas_contract != "" {
		res, err := setConfig(sysConfigClient, ctx, request.Gas_contract, vm.GasContractNameKey)
		if err != nil {
			return nil, err
		}
		result = append(result, res)
	}
	if request.VrfParams != "" {
		res, err := setConfig(sysConfigClient, ctx, request.VrfParams, vm.VrfParamsKey)
		if err != nil {
			return nil, err
		}
		result = append(result, res)
	}
	return result, nil
}

// 系统参数获取 sysconfig get
func (sysConfigClient SysConfigClient) GetTxGasLimit(ctx context.Context) (uint64, error) {
	funcName := "getTxGasLimit"
	result, err := sysConfigClient.contractCallWithParams(ctx, nil, funcName, precompile.ParameterManagementAddress)
	if err != nil {
		return 0, err
	}
	res := result.([]interface{})
	return res[0].(uint64), nil
}

func (sysConfigClient SysConfigClient) GetBlockGasLimit(ctx context.Context) (uint64, error) {
	funcName := "getBlockGasLimit"
	result, err := sysConfigClient.contractCallWithParams(ctx, nil, funcName, precompile.ParameterManagementAddress)
	if err != nil {
		return 0, err
	}
	res := result.([]interface{})
	return res[0].(uint64), nil
}

func (sysConfigClient SysConfigClient) GetGasContractName(ctx context.Context) (string, error) {
	funcName := "getGasContractName"
	result, err := sysConfigClient.contractCallWithParams(ctx, nil, funcName, precompile.ParameterManagementAddress)
	if err != nil {
		return "", err
	}
	res := result.([]interface{})
	return res[0].(string), nil
}

func (sysConfigClient SysConfigClient) GetIsProduceEmptyBlock(ctx context.Context) (uint32, error) {
	funcName := "getIsProduceEmptyBlock"
	result, err := sysConfigClient.contractCallWithParams(ctx, nil, funcName, precompile.ParameterManagementAddress)
	if err != nil {
		return 0, err
	}
	res := result.([]interface{})
	return res[0].(uint32), nil
}

func (sysConfigClient SysConfigClient) GetCheckContractDeployPermission(ctx context.Context) (uint32, error) {
	funcName := "getCheckContractDeployPermission"
	result, err := sysConfigClient.contractCallWithParams(ctx, nil, funcName, precompile.ParameterManagementAddress)
	if err != nil {
		return 0, err
	}
	res := result.([]interface{})
	return res[0].(uint32), nil
}

func (sysConfigClient SysConfigClient) GetAllowAnyAccountDeployContract(ctx context.Context) (uint32, error) {
	funcName := "getAllowAnyAccountDeployContract"
	result, err := sysConfigClient.contractCallWithParams(ctx, nil, funcName, precompile.ParameterManagementAddress)
	if err != nil {
		return 0, err
	}
	res := result.([]interface{})
	return res[0].(uint32), nil
}

func (sysConfigClient SysConfigClient) GetIsApproveDeployedContract(ctx context.Context) (uint32, error) {
	funcName := "getIsApproveDeployedContract"
	result, err := sysConfigClient.contractCallWithParams(ctx, nil, funcName, precompile.ParameterManagementAddress)
	if err != nil {
		return 0, err
	}
	res := result.([]interface{})
	return res[0].(uint32), nil
}

func (sysConfigClient SysConfigClient) GetIsTxUseGas(ctx context.Context) (uint32, error) {
	funcName := "getIsTxUseGas"
	result, err := sysConfigClient.contractCallWithParams(ctx, nil, funcName, precompile.ParameterManagementAddress)
	if err != nil {
		return 0, err
	}
	res := result.([]interface{})
	return res[0].(uint32), nil
}

func (sysConfigClient SysConfigClient) GetVRFParams(ctx context.Context) (string, error) {
	funcName := "getVRFParams"
	result, err := sysConfigClient.contractCallWithParams(ctx, nil, funcName, precompile.ParameterManagementAddress)
	if err != nil {
		return "", err
	}
	res := result.([]interface{})
	return res[0].(string), nil
}

func setConfig(sysConfigClient SysConfigClient, ctx context.Context, param string, name string) (string, error) {
	// todo: optimize the code, param check, param convert
	var funcName string
	if !checkConfigParam(param, name) {
		return "", errors.New("config param is error")
	}

	newParam, err := sysConfigConvert(param, name)
	if err != nil {
		fmt.Println(err.Error())
	}
	if name == "IsCheckContractDeployPermission" {
		funcName = "setCheckContractDeployPermission"
	} else {
		funcName = "set" + name
	}
	funcParams := packet.CombineFuncParams(newParam)

	result, err := sysConfigClient.contractCallWithParams(ctx, funcParams, funcName, precompile.ParameterManagementAddress)
	if err != nil {
		return "", err
	}

	res := result.([]interface{})
	return res[0].(string), nil
}

func checkConfigParam(param string, key string) bool {

	switch key {
	case "TxGasLimit":
		// number check
		num, err := strconv.ParseUint(param, 10, 0)
		if err != nil {
			log.Error("param invalid")

			return false
		}

		// param check
		isInRange := vm.TxGasLimitMinValue <= num && vm.TxGasLimitMaxValue >= num
		if !isInRange {
			fmt.Printf("the transaction gas limit should be within (%d, %d)\n",
				vm.TxGasLimitMinValue, vm.TxGasLimitMaxValue)
			return false
		}
	case "BlockGasLimit":
		num, err := strconv.ParseUint(param, 10, 0)
		if err != nil {
			fmt.Println("param invalid")

			return false
		}
		isInRange := vm.BlockGasLimitMinValue <= num && vm.BlockGasLimitMaxValue >= num
		if !isInRange {
			fmt.Printf("the block gas limit should be within (%d, %d)\n",
				vm.BlockGasLimitMinValue, vm.BlockGasLimitMaxValue)
			return false
		}
	default:
		if param == "" {
			return false
		}
	}
	return true
}

func sysConfigConvert(param, paramName string) (string, error) {

	if paramName == vm.TxGasLimitKey || paramName == vm.BlockGasLimitKey || paramName == vm.VrfParamsKey || paramName == vm.GasContractNameKey {
		return param, nil
	}

	conv := genConfigConverter(paramName)
	result, err := conv.Convert(param)
	if err != nil {
		return "", err
	}
	return result.(string), nil
}

func genConfigConverter(paramName string) *packet.Convert {
	switch paramName {
	case vm.IsTxUseGasKey:
		return packet.NewConvert(txUseGas, txNotUseGas, "1", "0", paramName)
	case vm.IsApproveDeployedContractKey:
		return packet.NewConvert(conAudit, conNotAudit, "1", "0", paramName)
	case vm.IsCheckContractDeployPermissionKey:
		return packet.NewConvert(checkPerm, notCheckPerm, "1", "0", paramName)
	case vm.IsProduceEmptyBlockKey:
		return packet.NewConvert(prodEmp, notProdEmp, "1", "0", paramName)
	default:
		log.Error("invalid system configuration %v", paramName)
	}
	return nil
}

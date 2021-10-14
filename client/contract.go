package client

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/PlatONE_Network/PlatONE-SDK-Go/common"
	"github.com/PlatONE_Network/PlatONE-SDK-Go/packet"
)

type ContractClient struct {
	*Client
	//Address  string
	CodePath string
	AbiPath  string
	Vm       string
}

func (contractClient ContractClient) Deploy(txparam common.TxParams, consParams []string) ([]interface{}, error) {
	var abiBytes []byte
	var consArgs = make([]interface{}, 0)

	codeByte, err := common.ParamParse(contractClient.CodePath, "code")
	codeBytes := codeByte.([]byte)
	if contractClient.AbiPath != "" {
		abiByte, err := common.ParamParse(contractClient.AbiPath, "abi")
		if err != nil {
			return nil, err
		}
		abiBytes = abiByte.([]byte)
	}
	err = common.ParamValid(contractClient.Vm, "Vm")
	if err != nil {
		return nil, err
	}
	// 获取合约abi 函数
	conAbi, _ := packet.ParseAbiFromJson(abiBytes)
	constructor := conAbi.GetConstructor()
	if constructor != nil {
		consArgs, _ = constructor.StringToArgs(consParams)
	}
	dataGenerator := packet.NewDeployDataGen(conAbi)
	// set the virtual machine interpreter
	dataGenerator.SetInterpreter(contractClient.Vm, abiBytes, codeBytes, consArgs, constructor)

	return contractClient.Client.clientCommonV2(txparam, dataGenerator, nil, true)
}

// 列出该合约的所有方法
func (contractClient ContractClient) ListContractMethods() (packet.ContractAbi, error) {
	if contractClient.AbiPath == "" {
		return nil, errors.New("abi file is not found")
	}
	abiByte, err := common.ParamParse(contractClient.AbiPath, "abi")
	if err != nil {
		return nil, errors.New("parameter resolution failed")
	}
	abiBytes := abiByte.([]byte)
	return packet.ParseAbiFromJson(abiBytes)
	//abiFuncs,err := packet.ParseAbiFromJson(abiBytes)
	//if err != nil {
	//	return "", err
	//}
	//result := abiFuncs.ListAbiFuncName()
	//return result, nil
}

// execute a method in the contract(evm or wasm).
func (contractClient ContractClient) Execute(txparam common.TxParams, funcName string, funcParams []string, address string) ([]interface{}, error) {
	var res []interface{}
	isListMethods, err := contractClient.IsFuncNameInContract(funcName)
	if !isListMethods {
		return nil, err
	}
	funcName, funcParams = common.FuncParse(funcName, funcParams)

	result, err := contractClient.contractCallWrap(txparam, funcParams, funcName, address)
	for _, data := range result {
		if common.IsTypeLenLong(reflect.ValueOf(data)) {
			//fmt.Printf("result%d:\n%+v\n", i, data)
			res = append(res, data)
		} else {
			//fmt.Printf("result%d:%+v\n", i, data)
			res = append(res, data)
		}
	}
	return res, nil
}

// 通过cns 名字调用合约
func (contractClient ContractClient) CnsExecute(txparam common.TxParams, funcName string, funcParams []string, cns string) ([]interface{}, error) {
	var res []interface{}
	isListMethods, err := contractClient.IsFuncNameInContract(funcName)
	if !isListMethods {
		return nil, err
	}
	funcName, funcParams = common.FuncParse(funcName, funcParams)

	result, err := contractClient.contractCallWrap(txparam, funcParams, funcName, cns)
	for i, data := range result {
		if common.IsTypeLenLong(reflect.ValueOf(data)) {
			fmt.Printf("result%d:\n%+v\n", i, data)
			res = append(res, data)
		} else {
			fmt.Printf("result%d:%+v\n", i, data)
			res = append(res, data)
		}
	}
	return res, nil
}

// 判断该函数是否属于合约中的方法
func (contractClient ContractClient) IsFuncNameInContract(funcName string) (bool, error) {
	contracts, err := contractClient.ListContractMethods()
	if err != nil {
		return false, err
	}
	_, err = contracts.GetFuncFromAbi(funcName)
	if err != nil {
		return false, err
	}
	return true, nil
}

//// contractCall extract the common parts of the actions of contract execution
//func contractCall(c *cli.Context, funcParams []string, funcName, contract string) interface{} {
//	result := contractCallWrap(c, funcParams, funcName, contract)
//	return result[0]
//}

func (contractClient ContractClient) contractCallWrap(txparam common.TxParams, funcParams []string, funcName, contract string) ([]interface{}, error) {
	common.ParamValid(contractClient.Vm, "Vm")

	// get the abi bytes of the contracts
	abiPath := contractClient.AbiPath
	funcAbi := common.AbiParse(abiPath, contract)
	// abi bytes parsing
	contractAbi, _ := packet.ParseAbiFromJson(funcAbi)
	// find the method in abi obj.
	methodAbi, err := contractAbi.GetFuncFromAbi(funcName)
	if err != nil {
		return nil, err
	}
	// convert user input string to args in Golang
	funcArgs, _ := methodAbi.StringToArgs(funcParams)

	// judge whether the input string is contract Address or contract name
	cns, to, err := common.CnsParse(contract)
	if err != nil {
		return nil, err
	}

	/// dataGenerator := packet.NewContractDataGenWrap(funcName, funcParams, funcAbi, *cns, Vm)
	data := packet.NewData(funcArgs, methodAbi)
	dataGenerator := packet.NewContractDataGen(data, contractAbi, cns.TxType)
	dataGenerator.SetInterpreter(contractClient.Vm, cns.Name, cns.TxType)
	return contractClient.Client.clientCommonV2(txparam, dataGenerator, &to, true)
}

func (contractClient ContractClient) GetReceipt(txhash string) (*packet.Receipt, error) {
	var response interface{}
	_ = contractClient.RpcClient.Call(&response, "eth_getTransactionReceipt", txhash)
	if response == nil {
		return nil, nil
	}

	// parse the rpc response
	receipt, err := packet.ParseTxReceipt(response)
	if err != nil {
		return nil, err
	}

	return receipt, nil
}

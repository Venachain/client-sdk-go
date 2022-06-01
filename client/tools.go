package client

import (
	"git-c.i.wxblockchain.com/vena/src/client-sdk-go/packet"
	"git-c.i.wxblockchain.com/vena/src/client-sdk-go/venachain/keystore"
)

type ContractType int32

const (
	UnknownType ContractType = 0
	WasmVmType  ContractType = 1
	EvmVmType   ContractType = 2
)

type ExecuteContract struct {
	// 合约地址
	Address string `json:"address"`
	// 合约类型, 1 wasm, 2 evm
	Type ContractType `json:"type"`
	// 合约方法
	Func string `json:"func"`
	// 合约参数
	FuncParams []string `json:"funcParams"`
}

func GetContractCallParamsSigned(executeContract ExecuteContract, contractAbiPath string, key *keystore.Key) (string, error) {
	contractContent, err := GenContractContent(contractAbiPath)
	if err != nil {
		return "", err
	}
	cns, to, err := packet.CnsParse(executeContract.Address)
	if err != nil {
		return "", err
	}
	methodAbi, err := contractContent.GetFuncFromAbi(executeContract.Func)
	if err != nil {
		return "", err
	}
	funcArgs, err := methodAbi.StringToArgs(executeContract.FuncParams)
	if err != nil {
		return "", err
	}
	data := packet.NewData(funcArgs, methodAbi)
	dataGenerator := packet.NewContractDataGen(data, contractContent, cns.TxType)
	vmType := "wasm"
	if executeContract.Type == EvmVmType {
		vmType = "evm"
	}
	dataGenerator.SetInterpreter(vmType, cns.Name, cns.TxType)
	dataGenerator.To = to
	txparam, err := dataGenerator.MakeTxparamForContract(&key.Address, &dataGenerator.To)
	if err != nil {
		return "", err
	}
	return txparam.GetSignedTx(key)
}

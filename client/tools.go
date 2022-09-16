package client

import (
	"git-c.i.wxblockchain.com/vena/src/client-sdk-go/common"
	"git-c.i.wxblockchain.com/vena/src/client-sdk-go/packet"
	vena_common "git-c.i.wxblockchain.com/vena/src/client-sdk-go/venachain/common"
	"git-c.i.wxblockchain.com/vena/src/client-sdk-go/venachain/keystore"
)

type ExecuteContract struct {
	// 合约地址
	Address string `json:"address"`
	// 合约类型, 1 wasm, 2 evm
	Type string `json:"type"`
	// 合约方法
	Func string `json:"func"`
	// 合约参数
	FuncParams []string `json:"funcParams"`
}

func NewExecuteContract(contractAddress string, vmType string, funcName string, funcParams []string) *ExecuteContract {
	return &ExecuteContract{
		contractAddress,
		vmType,
		funcName,
		funcParams,
	}
}

func GetTxSigned(executeContract ExecuteContract, contractAbiPath string, key *keystore.Key) (string, error) {
	txparam, err := GetTxparam(executeContract, contractAbiPath, key.Address.Hex())
	if err != nil {
		return "", err
	}
	return txparam.GetSignedTx(key)
}

func GetTxparam(executeContract ExecuteContract, contractAbiPath string, account string) (*common.TxParams, error) {
	contractContent, err := GenContractContent(contractAbiPath)
	if err != nil {
		return nil, err
	}
	cns, to, err := packet.CnsParse(executeContract.Address)
	if err != nil {
		return nil, err
	}
	methodAbi, err := contractContent.GetFuncFromAbi(executeContract.Func)
	if err != nil {
		return nil, err
	}
	funcArgs, err := methodAbi.StringToArgs(executeContract.FuncParams)
	if err != nil {
		return nil, err
	}
	data := packet.NewData(funcArgs, methodAbi)
	dataGenerator := packet.NewContractDataGen(data, contractContent, cns.TxType)
	vmType := "wasm"
	if executeContract.Type == "evm" {
		vmType = "evm"
	}
	dataGenerator.SetInterpreter(vmType, cns.Name, cns.TxType)
	dataGenerator.To = to
	userAccount := vena_common.HexToAddress(account)
	txparam, err := dataGenerator.MakeTxparamForContract(&userAccount, &dataGenerator.To)
	if err != nil {
		return nil, err
	}
	return txparam, nil
}

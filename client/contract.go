package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"git-c.i.wxblockchain.com/vena/src/client-sdk-go/common"
	"git-c.i.wxblockchain.com/vena/src/client-sdk-go/log"
	"git-c.i.wxblockchain.com/vena/src/client-sdk-go/packet"
	common_plaone "git-c.i.wxblockchain.com/vena/src/client-sdk-go/venachain/common"
	"git-c.i.wxblockchain.com/vena/src/client-sdk-go/venachain/keystore"
)

type ContractClient struct {
	*Client
	ContractContent *packet.ContractContent
	VmType          string
}

// contract：合约abi 文件的位置或合约地址
// vm： 虚拟机 evm 或 wasm，不传默认为wasm
func NewContractClient(ctx context.Context, url URL, keyfilePath, passphrase, contract, vmType string) (*ContractClient, error) {
	err := packet.ParamValid(vmType, "VmType")
	if err != nil {
		return nil, err
	}
	client, err := NewClient(ctx, url, keyfilePath, passphrase)
	if err != nil {
		return nil, err
	}
	contractContent, err := GenContractContent(contract)
	if err != nil {
		return nil, err
	}
	contractClient := &ContractClient{
		client,
		&contractContent,
		vmType,
	}
	return contractClient, nil
}

// 传入key 构造合约客户端
func NewContractClientWithKey(ctx context.Context, url URL, key *keystore.Key, contract, vmType string) (*ContractClient, error) {
	err := packet.ParamValid(vmType, "VmType")
	if err != nil {
		return nil, err
	}
	client, err := NewClientWithKey(ctx, url, key)
	if err != nil {
		return nil, err
	}
	contractContent, err := GenContractContent(contract)
	if err != nil {
		return nil, err
	}
	contractClient := &ContractClient{
		client,
		&contractContent,
		vmType,
	}
	return contractClient, nil
}

// execute a method in the contract(evm or wasm)
// contract 可以为合约地址或cns 名字
func (contractClient ContractClient) Execute(ctx context.Context, funcName string, funcParams []string, contract string, sync bool) (interface{}, error) {
	funcName, funcParams = packet.FuncParse(funcName, funcParams)
	// 构造dataGenerator
	dataGenerator, err := contractClient.MakeContractGenerator(contract, funcParams, funcName)
	if err != nil {
		return nil, err
	}
	result, err := contractClient.contractCall(ctx, dataGenerator, sync)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// consParams 为solidyty 合约中constructor的相关参数
func (contractClient ContractClient) Deploy(ctx context.Context, abipath string, codepath string, consParams []string, sync bool) (interface{}, error) {
	// 构造dataGenerator
	dataGenerator, err := contractClient.MakeDeployGenerator(abipath, codepath, consParams)
	if err != nil {
		return nil, err
	}
	txParams, err := MakeTxparamForDeploy(dataGenerator, &contractClient.Key.Address)
	if err != nil {
		return nil, err
	}
	result, err := contractClient.MessageCall(ctx, dataGenerator, *txParams, contractClient.Key, sync)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// 列出该合约的所有方法
func (contractClient ContractClient) ListContractMethods() (packet.ContractContent, error) {
	if contractClient.ContractContent == nil {
		return nil, errors.New("get contract content is nil")
	}
	return *contractClient.ContractContent, nil
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

func (contractClient ContractClient) GetReceipt(txhash string) (*packet.Receipt, error) {
	var res interface{}
	response, err := contractClient.RpcClient.Call(context.Background(), "eth_getTransactionReceipt", txhash)
	if err != nil {
		return nil, err
	}
	if response == nil {
		return nil, nil
	}
	if err := json.Unmarshal(response, &res); err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	// parse the rpc response
	receipt, err := packet.ParseTxReceipt(res)
	if err != nil {
		return nil, err
	}
	return receipt, nil
}

func (contractClient ContractClient) MakeContractGenerator(contract string, funcParams []string, funcName string) (*packet.ContractDataGen, error) {
	// judge whether the input string is contract Address or contract name
	cns, to, err := packet.CnsParse(contract)
	if err != nil {
		return nil, err
	}
	contractContent := contractClient.ContractContent
	// find the method in abi obj.
	methodAbi, err := contractContent.GetFuncFromAbi(funcName)
	if err != nil {
		return nil, err
	}
	// convert user input string to args in Golang
	funcArgs, _ := methodAbi.StringToArgs(funcParams)
	//if err != nil {
	//	return nil, err
	//}
	data := packet.NewData(funcArgs, methodAbi)
	dataGenerator := packet.NewContractDataGen(data, *contractContent, cns.TxType)
	dataGenerator.SetInterpreter(contractClient.VmType, cns.Name, cns.TxType)
	dataGenerator.To = to
	return dataGenerator, nil
}

func (contractClient ContractClient) MakeDeployGenerator(abipath string, codepath string, consParams []string) (*packet.DeployDataGen, error) {
	var consArgs = make([]interface{}, 0)
	if codepath == "" || abipath == "" {
		return nil, errors.New("code path or abi path is nil")
	}
	abiBytes, err := getAbiBytes(abipath)
	if err != nil {
		return nil, err
	}
	codeBytes, err := getCodeBytes(codepath)
	if err != nil {
		return nil, err
	}
	conAbi, err := GetContractByAbiPath(abipath)
	if err != nil {
		return nil, err
	}
	constructor := conAbi.GetConstructor()
	if constructor != nil {
		consArgs, _ = constructor.StringToArgs(consParams)
	}
	dataGenerator := packet.NewDeployDataGen(conAbi)
	// set the virtual machine interpreter
	dataGenerator.SetInterpreter(contractClient.VmType, abiBytes, codeBytes, consArgs, constructor)
	return dataGenerator, nil
}

func (contractClient ContractClient) SendTxparam(ctx context.Context, txparam *common.TxParams) (interface{}, error) {
	res, err := contractClient.Send(ctx, txparam, contractClient.Key)
	if err != nil {
		return nil, err
	}
	polRes, err := contractClient.GetReceiptByPolling(res)
	if err != nil {
		log.Error("error:%s", err)
		return res, nil
	}
	receiptBytes, err := json.MarshalIndent(polRes, "", "\t")
	if err != nil {
		return nil, err
	}
	log.Info(string(receiptBytes))
	return string(receiptBytes), nil
}

func (contractClient ContractClient) CallTxparam(ctx context.Context, txparam *common.TxParams) (interface{}, error) {
	var params = make([]interface{}, 0)
	params = append(params, txparam)
	params = append(params, "latest")
	action := "eth_call"
	// send the RPC calls
	var resp string
	result, err := contractClient.RpcClient.Call(ctx, action, params...)
	if err != nil {
		return nil, errors.New("send Transaction through http error")
	}
	err = json.Unmarshal(result, &resp)
	if err != nil {
		return nil, err
	}
	fmt.Println(resp)
	return resp, nil
}

// 封装合约的方法,同步获取receipt
func (contractClient ContractClient) contractCall(ctx context.Context, dataGenerator *packet.ContractDataGen, sync bool) (interface{}, error) {
	// 构造txparam
	txparam, err := dataGenerator.MakeTxparamForContract(&contractClient.Key.Address, &dataGenerator.To)
	if err != nil {
		return nil, err
	}
	result, err := contractClient.MessageCall(ctx, dataGenerator, *txparam, contractClient.Key, sync)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// 封装合约的方法,同步获取receipt
func (contractClient ContractClient) contractCallWithParams(ctx context.Context, funcParams []string, funcName string, contract string) (interface{}, error) {
	dataGenerator, err := contractClient.MakeContractGenerator(contract, funcParams, funcName)
	if err != nil {
		return nil, err
	}
	// 构造txparam
	txparam, err := dataGenerator.MakeTxparamForContract(&contractClient.Key.Address, &dataGenerator.To)
	if err != nil {
		return nil, err
	}
	result, err := contractClient.MessageCall(ctx, dataGenerator, *txparam, contractClient.Key, true)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func getAbiBytes(abipath string) ([]byte, error) {
	var abiBytes []byte
	abiByte, err := packet.ParamParse(abipath, "abi")
	if err != nil {
		return nil, err
	}
	abiBytes = abiByte.([]byte)
	return abiBytes, nil
}

func getCodeBytes(codepath string) ([]byte, error) {
	codeByte, err := packet.ParamParse(codepath, "code")
	if err != nil {
		return nil, err
	}
	codeBytes := codeByte.([]byte)
	return codeBytes, nil
}

func MakeTxparamForDeploy(dataGenerator *packet.DeployDataGen, from *common_plaone.Address) (*common.TxParams, error) {
	txparam := common.TxParams{}
	var err error
	if from != nil {
		txparam.From = *from
	} else {
		return nil, errors.New("the from of the transaction is empty")
	}
	txparam.Data, err = dataGenerator.CombineData()
	if err != nil {
		return nil, errors.New("packet data err: %s\n")
	}
	return &txparam, nil
}
func GenContractContent(contract string) (packet.ContractContent, error) {
	var contractContent packet.ContractContent
	var err error
	if contract == "" {
		contractContent = nil
	} else if packet.IsMatch(contract, "address") { // 生成预编译合约内容
		contractContent, err = GetContractByContractAddress(contract)
		if err != nil {
			return nil, err
		}
	} else if strings.Contains(contract, ".json") || strings.Contains(contract, ".abi") {
		contractContent, err = GetContractByAbiPath(contract)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("abipaht or precompiled contract address is wrong")
	}
	return contractContent, nil
}

// 通过abiPath 或者预编译合约的地址获得合约内容
func GetContractByAbiPath(abiPath string) (packet.ContractContent, error) {
	return getContractContent(abiPath, "")
}

// 通过预编译合约的地址获得合约内容
func GetContractByContractAddress(precompiledContractAddress string) (packet.ContractContent, error) {
	return getContractContent("", precompiledContractAddress)
}

// get the abi bytes of the contracts
func getContractContent(abiPath string, contract string) (packet.ContractContent, error) {
	funcAbi := packet.AbiParse(abiPath, contract)
	// abi bytes parsing
	contractContent, err := packet.ParseAbiFromJson(funcAbi)
	if err != nil {
		return nil, err
	}
	return contractContent, nil
}

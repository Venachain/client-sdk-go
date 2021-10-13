package packet

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PlatONE_Network/PlatONE-SDK-Go/platone/abi"
	"github.com/PlatONE_Network/PlatONE-SDK-Go/platone/crypto"
	"strings"
)

//
//================================ABI================================

// FuncDesc, the object of the contract abi files
type FuncDesc struct {
	Name            string                   `json:"name"`
	Inputs          []abi.ArgumentMarshaling `json:"inputs"`
	Outputs         []abi.ArgumentMarshaling `json:"outputs"`
	Constant        interface{}              `json:"constant"`
	Type            string                   `json:"type"`
	StateMutability string                   `json:"stateMutability,omitempty"` // tag for solidity ver > 0.6.0
}

type ContractAbi []*FuncDesc

// ParseAbiFromJson parses the application binary interface(abi) files to []FuncDesc object array
func ParseAbiFromJson(abiBytes []byte) (ContractAbi, error) {
	var a ContractAbi

	if abiBytes == nil {
		return nil, errors.New("abiBytes are null")
	}

	if err := json.Unmarshal(abiBytes, &a); err != nil {
		return nil, errors.New("unmarshal %s bytes error")
	}

	return a, nil
}

// ParseFuncFromAbi searches the function (or event) names in the []FuncDesc object array
func (funcs ContractAbi) GetFuncFromAbi(name string) (*FuncDesc, error) {
	for _, value := range funcs {
		if strings.EqualFold(value.Name, name) {
			return value, nil
		}
	}

	if name == "" {
		name = "null"
	}
	funcList := funcs.ListAbiFuncName()

	return nil, fmt.Errorf("function/event %s is not found in\n%s", name, funcList)
}

func (funcs ContractAbi) GetConstructor() *FuncDesc {
	for _, value := range funcs {
		if strings.EqualFold(value.Type, "constructor") {
			return value
		}
	}

	return nil
}

func (funcs ContractAbi) GetEvents() []*FuncDesc {
	var events = make([]*FuncDesc, 0)

	for _, value := range funcs {
		if strings.EqualFold(value.Type, "event") {
			events = append(events, value)
		}
	}

	return events
}

// 把abi文件中函数的列表转换为string 格式
func (abiFuncs ContractAbi) ListAbiFuncName() string {
	var result string
	for _, function := range abiFuncs {
		strInput := []string{}
		strOutput := []string{}
		for _, param := range function.Inputs {
			strInput = append(strInput, param.Name+" "+param.Type)
		}
		for _, param := range function.Outputs {
			strOutput = append(strOutput, param.Name+" "+param.Type)
		}
		result += fmt.Sprintf("%s: ", function.Type)
		result += fmt.Sprintf("%s(%s)%s\n", function.Name, strings.Join(strInput, ","), strings.Join(strOutput, ","))
	}

	return result
}

func (abiFunc *FuncDesc) StringToArgs(funcParams []string) ([]interface{}, error) {
	var arguments abi.Arguments
	var argument abi.Argument

	var args = make([]interface{}, 0)

	var err error

	// Judging whether the number of inputs matches
	if len(abiFunc.Inputs) != len(funcParams) {
		return nil, fmt.Errorf("param check error, required %d inputs, recieved %d.\n", len(abiFunc.Inputs), len(funcParams))
	}

	for i, v := range abiFunc.Inputs {
		if argument.Type, err = abi.NewTypeV2(v.Type, v.InternalType, v.Components); err != nil {
			return nil, err
		}
		arguments = append(arguments, argument)

		/// arg, err := abi.SolInputTypeConversion(input.Type, v)
		arg, err := argument.Type.StringConvert(funcParams[i])
		if err != nil {
			return nil, err
		}

		args = append(args, arg)
	}

	return args, nil
}

func (abiFunc *FuncDesc) getParamType() []string {
	var paramTypes = make([]string, 0)

	for _, v := range abiFunc.Inputs {
		paramTypes = append(paramTypes, GenFuncSig(v))
	}

	return paramTypes
}

//// ======================== Precompiled Contract Return ========================
//
//// ContractReturn, system contract return object
//type ContractReturn struct {
//	Code int         `json:"code"`
//	Msg  string      `json:"msg"`
//	Data interface{} `json:"data"`
//}
//
//// ParseSysContractResult parsed the result to ContractReturn object
//func ParseSysContractResult(result []byte) (*ContractReturn, error) {
//	a := ContractReturn{} //删除？
//	err := json.Unmarshal(result, &a)
//	if err != nil {
//		// packet.Fatalf(utl.ErrUnmarshalBytesFormat, "contract return", err.Error())
//		errStr := fmt.Sprintf(packet.ErrUnmarshalBytesFormat, "contract return", err.Error())
//		return nil, errors.New(errStr)
//	}
//
//	return &a, nil
//}
//
////=========================Transaction=================================
//
//// TxParamsDemo, the object of the eth_call, eth_sendTransaction
//type TxParams struct {
//	From     common.Address  `json:"from"` // the address used to send the transaction
//	To       *common.Address `json:"to"`   // the address receives the transactions
//	Gas      string          `json:"gas"`
//	GasPrice string          `json:"gasPrice"`
//	Value    string          `json:"value"`
//	Data     string          `json:"data"`
//}
//
//// NewTxParams news a TxParams object
//func NewTxParams(from common.Address, to *common.Address, value, gas, gasPrice, data string) *TxParams {
//
//	tx := &TxParams{
//		From:     from,
//		To:       to,
//		GasPrice: gasPrice,
//		Gas:      gas,
//		Value:    value,
//		Data:     data,
//	}
//
//	return tx
//}
//
//// SendMode selects the rpc calls (eth_call, eth_sendTransaction, and eth_sendRawTransaction)
//func (tx *TxParams) SendMode(isWrite bool, keyfile *packet.Keyfile) ([]interface{}, string) {
//	var action string
//	var params = make([]interface{}, 0)
//
//	switch {
//	case !isWrite:
//		params = append(params, tx)
//		params = append(params, "latest")
//		action = "eth_call"
//	case keyfile.Json != nil:
//		signedTx, _ := tx.GetSignedTx(keyfile)
//		params = append(params, signedTx)
//		action = "eth_sendRawTransaction"
//	default:
//		params = append(params, tx)
//		action = "eth_sendTransaction"
//	}
//
//	return params, action
//}
//

//// GetSignedTx gets the signed transaction
//func (tx *TxParams) GetSignedTx(keyfile *packet.Keyfile) (string, error) {
//
//	var txSign *types.Transaction
//
//	// convert the TxParams object to types.Transaction object
//	nonce := getNonceRand()
//	value, _ := hexutil.DecodeBig(tx.Value)
//	gas, _ := hexutil.DecodeUint64(tx.Gas)
//	gasPrice, _ := hexutil.DecodeBig(tx.GasPrice)
//	data, _ := hexutil.Decode(tx.Data)
//
//	if tx.To == nil {
//		txSign = types.NewContractCreation(nonce, value, gas, gasPrice, data)
//	} else {
//		txSign = types.NewTransaction(nonce, *tx.To, value, gas, gasPrice, data)
//	}
//
//	// extract pk from keystore file and sign the transaction
//	// deprecated: move to the outter method
//	/*
//		priv, err := keyfile.GetPrivateKey()
//		if err != nil {
//			return "", err
//		}*/
//
//	// todo: choose the correct signer
//	txSign, _ = types.SignTx(txSign, types.HomesteadSigner{}, keyfile.GetPrivateKey())
//	/// txSign, _ = types.SignTx(txSign, types.NewEIP155Signer(big.NewInt(300)), priv)
//	/// utl.Logger.Printf("the signed transaction is %v\n", txSign)
//
//	str, err := rlpEncode(txSign)
//	if err != nil {
//		return "", err
//	}
//
//	return str, nil
//}
//
//// getNonceRand generate a random nonce
//// Warning: if the design of the nonce mechanism is modified
//// this part should be modified as well
//func getNonceRand() uint64 {
//	rand.Seed(time.Now().Unix())
//	return rand.Uint64()
//}
//
//// todo: move to packet.go ?
var errorSig = crypto.Keccak256([]byte("Error(string)"))[:4]

func UnpackError(res []byte) (string, error) {
	var revStr string

	if !bytes.Equal(res[:4], errorSig) {
		return "<not revert string>", errors.New("not a revert string")
	}

	typ, _ := abi.NewTypeV2("string", "", nil)
	err := abi.Arguments{{Type: typ}}.UnpackV2(&revStr, res[4:])
	if err != nil {
		return "<invalid revert string>", err
	}

	return revStr, nil
}

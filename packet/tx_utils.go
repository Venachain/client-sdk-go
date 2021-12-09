package packet

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"git-c.i.wxblockchain.com/PlatONE/src/node/client-sdk-go/platone/abi"
	"git-c.i.wxblockchain.com/PlatONE/src/node/client-sdk-go/platone/crypto"
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

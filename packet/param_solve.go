package packet

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"git-c.i.wxblockchain.com/vena/src/client-sdk-go/log"
	precompile "git-c.i.wxblockchain.com/vena/src/client-sdk-go/precompiled"
	"git-c.i.wxblockchain.com/vena/src/client-sdk-go/venachain/abi"
	common_venachain "git-c.i.wxblockchain.com/vena/src/client-sdk-go/venachain/common"
	"git-c.i.wxblockchain.com/vena/src/client-sdk-go/venachain/vm"
)

const defaultTxType = 2

// ParamValid check if the input is valid
func ParamValid(param, paramName string) error {

	valid := ParamValidWrap(param, paramName)

	if !valid {
		return errors.New("param is invalid")
	}
	return nil
}

//===============================Abi Parsing========================================
// AbiParse gets the abi bytes by the input parameters provided
// The abi file can be obtained through following ways:
// 1. user provide the abi file path
// 2. abiBytes of precompiled contracts (see precompiled/bindata.go)
// (currently, the following features are not enabled)
// a. get the abi files from default abi file locations
// b. get the abi bytes on chain (wasm contract only).
func AbiParse(abiFilePath, address string) []byte {
	var err error
	var abiBytes []byte

	if abiFilePath == "" {
		if p := precompile.List[address]; p != "" {
			abibyte, err_byte := precompile.GetContractByte(p)
			if err_byte != nil {
				return nil
			}
			return abibyte
		}
	}

	abiBytes, err = ParseFileToBytes(abiFilePath)
	if err != nil {
		log.Error("abiParse fail")
	}

	return abiBytes
}

//===============================User Input Convert=======================================

// convert, convert user input from key to value
type Convert struct {
	key1      string      // user input 1
	key2      string      // user input 2
	value1    interface{} // the convert value of user input 1
	value2    interface{} // the convert value of user input 2
	paramName string
}

// Some of the contract function inputs are numbers,
// these numbers are hard for users to remember the meanings behind them,
// Thus, to simplify the user input, we convert the meaningful strings to number automatically
// For example, if user input: "valid", the converter will convert the string to 1
func NewConvert(key1, key2 string, value1, value2 interface{}, paramName string) *Convert {
	return &Convert{
		key1:      key1,
		key2:      key2,
		value1:    value1,
		value2:    value2,
		paramName: paramName,
	}
}

func ConvertSelect(param, paramName string) (interface{}, error) {
	var conv *Convert

	switch paramName {
	case "operation": // registration operation
		conv = NewConvert("approve", "reject", "2", "3", paramName)
	case "status": // node status
		conv = NewConvert("valid", "invalid", 1, 2, paramName)
	case "type": // node type
		conv = NewConvert("consensus", "observer", 1, 0, paramName)
	default:
		return nil, errors.New("no suitable param, conver param fail")
	}

	return conv.Convert(param)
}

func (conv *Convert) Convert(param string) (interface{}, error) {
	key1NotEqual := !strings.EqualFold(param, conv.key1)
	key2NotEqual := !strings.EqualFold(param, conv.key2)

	if key1NotEqual && key2NotEqual {
		return nil, fmt.Errorf("the %s should be either \"%s\" or \"%s\"", conv.paramName, conv.key1, conv.key2)
	}

	if key2NotEqual {
		return conv.value1, nil
	} else {
		return conv.value2, nil
	}
}

func (conv *Convert) Parse(param interface{}) string {

	value1NotEqual := param != conv.value1
	value2NotEqual := param != conv.value2

	if value1NotEqual && value2NotEqual {
		panic("not match")
	}

	if value2NotEqual {
		return conv.key1
	} else {
		return conv.key2
	}
}

// ========================== Param Convert ========================================

// ParamParse convert the user inputs to the value needed
func ParamParse(param, paramName string) (interface{}, error) {
	var err error
	var i interface{}

	switch paramName {
	case "contract", "user":
		i = IsNameOrAddress(param)
		if i == CnsIsUndefined {
			err = fmt.Errorf(ErrParamInValidSyntax, "name or contract address")
		}
	case "delayNum", "p2pPort", "rpcPort":
		if IsInRange(param, 65535) {
			i, err = strconv.ParseInt(param, 10, 0)
		} else {
			err = errors.New("value out of range")
		}
	case "operation", "status", "type":
		i, err = ConvertSelect(param, paramName)
	case "code", "abi":
		i, err = ParseFileToBytes(param)
	default:
		i, err = param, nil
	}
	if err != nil {
		return nil, errors.New("paramParse fail")
	}

	return i, nil
}

//================================CNS=================================
const (
	CnsIsName int32 = iota
	CnsIsAddress
	CnsIsUndefined
)

type Cns struct {
	/// To     common.Address
	Name   string // the cns name of contract
	TxType uint64 // the transaction type of the contract execution (EXECUTE_CONTRACT or CNS_TX_TYPE)
}

func NewCns(name string, txType uint64) *Cns {
	return &Cns{
		/// To:     common.HexToAddress(to),
		Name:   name,
		TxType: txType,
	}
}

// CnsParse judge whether the input string is contract address or contract name
// and return the corresponding infos
func CnsParse(contract string) (*Cns, common_venachain.Address, error) {
	isAddress := IsNameOrAddress(contract)

	switch isAddress {
	case CnsIsAddress:
		return NewCns("", defaultTxType), common_venachain.HexToAddress(contract), nil
	case CnsIsName:
		return NewCns(contract, defaultTxType), common_venachain.HexToAddress(precompile.CnsInvokeAddress), nil
	default:
		return nil, common_venachain.Address{}, fmt.Errorf(ErrParamInValidSyntax, "contract address")
	}
}

// IsNameOrAddress Judge whether the input string is an address or a name
func IsNameOrAddress(str string) int32 {
	var valid int32

	switch {
	case IsMatch(str, "address"):
		valid = CnsIsAddress
	case IsMatch(str, "name") &&
		!strings.HasPrefix(strings.ToLower(str), "0x"):
		valid = CnsIsName
	default:
		valid = CnsIsUndefined
	}

	return valid
}

//===============================User Input Verification=======================================

func ParamValidWrap(param, paramName string) bool {
	var valid = true

	switch paramName {
	case "fw":
		if param != "*" {
			valid = IsMatch(param, "address")
		}
	case "to":
		valid = param == "" || IsMatch(param, "address")
	case "contract":
		valid = IsMatch(param, "address") || IsMatch(param, "name")
	case "action":
		valid = strings.EqualFold(param, "accept") || strings.EqualFold(param, "reject")
	case "vm":
		valid = param == "" || strings.EqualFold(param, "evm") || strings.EqualFold(param, "wasm")
	case "ipAddress":
		valid = IsUrl(param)
	case "externalIP", "internalIP":
		valid = IsUrl(param + ":0")
	//case "version":
	//	valid = utl.IsVersion(param)
	case "role":
		valid = IsRoleMatch(param)
	case "roles":
		valid = IsValidRoles(param)
	case "email", "mobile", "version", "num":
		valid = IsMatch(param, paramName)

	// newly added for restful server
	// todo; fix the toLower problem
	case "orgin", "address":
		valid = IsMatch(param, "address")
	case "contractname", "name":
		valid = IsMatch(param, "name")
	case "sysparam":
		valid = strings.EqualFold(param, "0") || strings.EqualFold(param, "1")
	case "blockgaslimit":
		num, err := strconv.ParseUint(param, 10, 0)
		if err != nil {
			return false
		}
		valid = vm.BlockGasLimitMinValue <= num && vm.BlockGasLimitMaxValue >= num
	case "txgaslimit":
		num, err := strconv.ParseUint(param, 10, 0)
		if err != nil {
			return false
		}
		valid = vm.TxGasLimitMinValue <= num && vm.TxGasLimitMaxValue >= num
	default:
		/// Logger.Printf("param valid function used but not validate the <%s> param\n", paramName)
	}

	return valid
}

// FuncParse wraps the GetFuncNameAndParams
// it separates the function method name and the parameters
func FuncParse(funcName string, funcParams []string) (string, []string) {
	var funcParamsNew []string

	if funcName == "" {
		return "", nil
	}

	funcName, funcParamsNew = GetFuncNameAndParams(funcName)
	if len(funcParamsNew) != 0 && len(funcParams) != 0 {
		fmt.Errorf("function parameters error")
	}
	funcParams = append(funcParams, funcParamsNew...)

	/// Logger.Printf("after function parse, the function is %s, %s", funcName, funcParams)
	return funcName, funcParams
}

// GetFuncNameAndParams parse the function params from the input string
func GetFuncNameAndParams(funcAndParams string) (string, []string) {
	// eliminate space
	f := TrimSpace(funcAndParams)

	hasBracket := strings.Contains(f, "(") && strings.Contains(f, ")")
	if !hasBracket {
		return f, nil
	}

	funcName := f[0:strings.Index(f, "(")]

	paramString := f[strings.Index(f, "(")+1 : strings.LastIndex(f, ")")]
	params := abi.GetFuncParams(paramString)

	return funcName, params
}

// TrimSpace trims all the space in the string
func TrimSpace(str string) string {
	strNoSpace := strings.Split(str, " ")
	return strings.Join(strNoSpace, "")
}

func IsTypeLenLong(v reflect.Value) bool {
	k := v.Kind()
	switch k {
	case reflect.Array, reflect.String, reflect.Slice, reflect.Map, reflect.Chan:
		return v.Len() > 20
	default:
		return false
	}
}

// CombineRule combines firewall rules
func CombineRule(addr, api string) string {
	return addr + ":" + api
}

// CombineFuncParams combines the function parameters
func CombineFuncParams(args ...string) []string {
	strArray := append([]string{}, args...)
	return strArray
}

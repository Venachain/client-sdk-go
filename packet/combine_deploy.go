package packet

import (
	"encoding/json"
	"errors"
	precompile "git-c.i.wxblockchain.com/PlatONE/src/node/client-sdk-go/precompiled"

	"git-c.i.wxblockchain.com/PlatONE/src/node/client-sdk-go/platone/abi"
)

// DeployCall, used for combining the data of contract deployment
type DeployDataGen struct {
	/// codeBytes         []byte
	/// abiBytes          []byte
	/// ConstructorParams []string
	Interpreter deployInter

	conAbi ContractContent
}

// NewDeployCall new a DeployCall object
func NewDeployDataGen(conAbi ContractContent) *DeployDataGen {
	var dataGen = new(DeployDataGen)
	dataGen.conAbi = conAbi

	return dataGen
}

func parseAbiConstructor(abiBytes []byte, funcParams []string) ([]byte, error) {
	var abiFunc *FuncDesc

	funcs, err := ParseAbiFromJson(abiBytes)
	if err != nil {
		return nil, err
	}

	for _, value := range funcs {
		if value.Type == "constructor" {
			abiFunc = value
			break
		}
	}

	// todo: better solution?
	if abiFunc == nil {
		return nil, nil
	}

	conBytes, _, err := EvmStringToEncodeByte(abiFunc, funcParams)
	return conBytes, err
}

// SetInterpreter set the interpreter of DeployCall object
func (dataGen *DeployDataGen) SetInterpreter(vmType string, abiBytes, codeBytes []byte, consParams []interface{}, methodAbi *FuncDesc) {

	if IsWasmContract(codeBytes) {
		// packet.Fatalf("the input  is not evm byte code")
		// return errors.New("the input is not evm byte code")
		vmType = "wasm"
	}

	switch vmType {
	case "evm":
		dataGen.Interpreter = &EvmDeployInterpreter{
			codeBytes:        codeBytes,
			constructorInput: consParams,
			constructorAbi:   methodAbi,
		}
	// the default interpreter is wasm
	default:
		dataGen.Interpreter = &WasmDeployInterpreter{
			codeBytes: codeBytes,
			abiBytes:  abiBytes,
		}
	}
}

func (dataGen *DeployDataGen) ReceiptParsing(receipt *Receipt) *ReceiptParsingReturn {
	return dataGen.Interpreter.ReceiptParsingV2(receipt, dataGen.conAbi)
}

// CombineData of DeployCall data struct is used for packeting the data of wasm or evm contracts deployment
// Implement the MessageCallDemo interface
func (dataGen DeployDataGen) CombineData() (string, error) {
	if dataGen.Interpreter == nil {
		return "", errors.New("interpreter is not provided")
	}

	return dataGen.Interpreter.combineData()
}

func (dataGen *DeployDataGen) GetIsWrite() bool {
	return true
}

func (dataGen *DeployDataGen) GetContractDataDen() *ContractDataGen {
	return nil
}

func (dataGen *DeployDataGen) ParseNonConstantResponse(respStr string, outputType []abi.ArgumentMarshaling) []interface{} {
	return nil
}

func getReceiptByte(receipt *Receipt) ([]byte, error) {
	return json.MarshalIndent(receipt, "", "\t")
}

func ReceiptParsing(receipt *Receipt, conAbi ContractContent) *ReceiptParsingReturn {

	var fn = WasmEventParsingPerLogV2
	var sysEvents = []string{precompile.PermDeniedEvent, precompile.CnsInitRegEvent}

	events := GetSysEvents(sysEvents)
	events = append(events, conAbi.GetEvents()...)

	return receipt.ParsingWrap(events, fn)
}

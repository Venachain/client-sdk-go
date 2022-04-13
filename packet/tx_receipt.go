package packet

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"git-c.i.wxblockchain.com/vena/src/client-sdk-go/types"

	precompile "git-c.i.wxblockchain.com/vena/src/client-sdk-go/precompiled"
	"git-c.i.wxblockchain.com/vena/src/client-sdk-go/venachain/abi"
	"git-c.i.wxblockchain.com/vena/src/client-sdk-go/venachain/common"
	"git-c.i.wxblockchain.com/vena/src/client-sdk-go/venachain/common/byteutil"
	"git-c.i.wxblockchain.com/vena/src/client-sdk-go/venachain/common/hexutil"
	"git-c.i.wxblockchain.com/vena/src/client-sdk-go/venachain/crypto"
	"git-c.i.wxblockchain.com/vena/src/client-sdk-go/venachain/rlp"
)

const (
	// ReceiptStatusFailed is the status code of a transaction if execution failed.
	ReceiptStatusFailed = uint64(0)

	// ReceiptStatusSuccessful is the status code of a transaction if execution succeeded.
	ReceiptStatusSuccessful = uint64(1)
)

var (
	txReceiptSuccessCode = hexutil.EncodeUint64(ReceiptStatusSuccessful)
	txReceiptFailureCode = hexutil.EncodeUint64(ReceiptStatusFailed)
)

const (
	TxReceiptSuccessMsg = "Operation Succeeded"
	TxReceiptFailureMsg = "Operation Failed"
)

type ReceiptParsingReturn struct {
	Status          string   `json:"status"`
	ContractAddress string   `json:"contractAddress,omitempty"`
	Logs            []string `json:"logs,omitempty"`
	BlockNumber     uint64   `json:"blockNumber"`
	GasUsed         uint64
	From            string
	To              string
	TxHash          string
	Err             string `json:"err,omitempty"`
}

func (r *ReceiptParsingReturn) String() string {
	/// rBytes, _ := json.Marshal(r)
	rBytes, _ := json.MarshalIndent(r, "", "\t")
	if rBytes == nil {
		return ""
	}

	return string(rBytes)
}

// Receipt, eth_getTransactionReceipt return data struct
type Receipt struct {
	BlockHash         string    `json:"blockHash"`         // hash of the block
	BlockNumber       string    `json:"blockNumber"`       // height of the block
	ContractAddress   string    `json:"contractAddress"`   // contract address of the contract deployment. otherwise null
	CumulativeGasUsed string    `json:"cumulativeGasUsed"` //
	From              string    `json:"from"`              // the account address used to send the transaction
	GasUsed           string    `json:"gasUsed"`           // gas used by executing the transaction
	Root              string    `json:"root"`
	To                string    `json:"to"`               // the address the transaction is sent to
	TransactionHash   string    `json:"transactionHash"`  // the hash of the transaction
	TransactionIndex  string    `json:"transactionIndex"` // the index of the transaction
	Logs              RecptLogs `json:"logs"`
	Status            string    `json:"status"` // the execution status of the transaction, "0x1" for success
}

type Log struct {
	Address string   `json:"address"`
	Topics  []string `json:"topics"`
	Data    string   `json:"data"`
}

type RecptLogs []*Log

// ParseSysContractResult parsed the rpc response to Receipt object
func ParseTxReceipt(response interface{}) (*Receipt, error) {
	var receipt = &Receipt{}
	if response == nil {
		return nil, errors.New("get receipt response is nil")
	}
	temp, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(temp, receipt)
	if err != nil {
		return nil, errors.New("error in transaction receipt")
	}

	return receipt, nil
}

func (receipt *Receipt) ParsingWrap(events []*FuncDesc, fn eventParsingFuncV2) *ReceiptParsingReturn {
	receiptParse := receipt.Parsing()
	receiptParse.Logs = EventParsingV2(receipt.Logs, events, fn)

	return receiptParse
}

func (receipt *Receipt) Parsing() *ReceiptParsingReturn {
	var recpParsing = new(ReceiptParsingReturn)

	recpParsing.Status = receiptStatusReturn(receipt.Status)
	recpParsing.BlockNumber, _ = hexutil.DecodeUint64(receipt.BlockNumber)
	recpParsing.ContractAddress = receipt.ContractAddress
	recpParsing.From = receipt.From
	recpParsing.To = receipt.To
	recpParsing.GasUsed, _ = hexutil.DecodeUint64(receipt.GasUsed)
	recpParsing.TxHash = receipt.TransactionHash

	/// recpParsing.Logs = EventParsingV2(receipt.Logs, events, fn)

	return recpParsing
}

func receiptStatusReturn(status string) (result string) {

	switch status {
	case txReceiptSuccessCode:
		result = TxReceiptSuccessMsg
	case txReceiptFailureCode:
		result = TxReceiptFailureMsg
	default:
		result = "undefined status. Something wrong"
	}

	return
}

// ======================== Receipt log Parsing =========================

type eventParsingFuncV2 func(*Log, []*FuncDesc) string

func GetSysEvents(SysEventList []string) []*FuncDesc {
	var events = make([]*FuncDesc, 0)

	for _, data := range SysEventList {
		p := precompile.List[data]
		abiBytes, _ := precompile.Asset(p)
		abiFunc, _ := ParseAbiFromJson(abiBytes)
		events = append(events, abiFunc.GetEvents()...)
	}

	return events
}

func EventParsingV2(logs RecptLogs, events []*FuncDesc, fn eventParsingFuncV2) []string {
	var res = make([]string, 0)

	for _, logData := range logs {
		result := fn(logData, events)
		if result != "" {
			res = append(res, result)
			//break
		}
	}

	return res
}

// ------------------------------ EVM --------------------------------------
func EvmEventParsingPerLogV2(eLog *Log, events []*FuncDesc) string {
	eventName, arguments := findEvmLogTopicV2(eLog.Topics[0], events)
	if arguments == nil {
		return ""
	}

	result := fmt.Sprintf("Event %s: ", eventName)
	rlpList := arguments.ReturnBytesUnpack(eLog.Data)

	for _, data := range rlpList {
		if data != nil && !reflect.ValueOf(data).IsZero() {
			result += fmt.Sprintf("%v ", data)
		}
	}

	return result
}

func EvmEventParsingLog(eLog *Log, events []*FuncDesc) []string {
	var result []string
	_, arguments := findEvmLogTopicV2(eLog.Topics[0], events)
	if arguments == nil {
		return nil
	}

	//result = append(result, fmt.Sprintf("Event%s: ", eventName))
	rlpList := arguments.ReturnBytesUnpack(eLog.Data)

	for _, data := range rlpList {
		if data != nil && !reflect.ValueOf(data).IsZero() {
			result = append(result, fmt.Sprintf("%v", data))
		}
	}

	return result
}

func findEvmLogTopicV2(topic string, events []*FuncDesc) (string, abi.Arguments) {

	for _, data := range events {

		if strings.EqualFold(evmLogTopicEncode(data), topic) {
			name := data.Name
			arguments := GenUnpackArgs(data.Inputs)
			return name, arguments
		}
	}

	return "", nil
}

// todo: similar to function selector???
// todo: optimization
func evmLogTopicEncode(data *FuncDesc) string {
	var strArray = make([]string, 0)

	for _, event := range data.Inputs {
		strArray = append(strArray, event.Type)
	}

	topic := data.Name + "(" + strings.Join(strArray, ",") + ")"
	return common.BytesToHash(crypto.Keccak256([]byte(topic))).String()
}

func GenUnpackArgs(data []abi.ArgumentMarshaling) (arguments abi.Arguments) {
	var argument abi.Argument

	for _, v := range data {
		argument.Type, _ = abi.NewTypeV2(v.Type, v.InternalType, v.Components)
		argument.Name = v.Name
		argument.Indexed = v.Indexed

		arguments = append(arguments, argument)
	}

	return
}

// --------------------------- WASM ------------------------------------
func WasmEventParsingPerLogV2(eLog *Log, events []*FuncDesc) string {
	var rlpList []interface{}

	eventName, topicTypes := findWasmLogTopicV2(eLog.Topics[0], events)

	if len(topicTypes) == 0 {
		return ""
	}

	dataBytes, _ := hexutil.Decode(eLog.Data)
	err := rlp.DecodeBytes(dataBytes, &rlpList)
	if err != nil {
		// todo: error handle
		fmt.Printf("the error is %v\n", err)
	}

	result := fmt.Sprintf("Event %s: ", eventName)
	result += parseReceiptLogData(rlpList, topicTypes)

	return result
}

func findWasmLogTopicV2(topic string, abiFunc []*FuncDesc) (string, []string) {

	for _, data := range abiFunc {
		if strings.EqualFold(wasmLogTopicEncode(data.Name), topic) {
			topicTypes := make([]string, 0)
			name := data.Name
			for _, v := range data.Inputs {
				topicTypes = append(topicTypes, v.Type)
			}
			return name, topicTypes
		}
	}

	return "", nil
}

func parseReceiptLogData(data []interface{}, types []string) string {
	var str string

	for i, v := range data {
		result := ConvertRlpBytesTo(v.([]uint8), types[i])
		str += fmt.Sprintf("%v ", result)
	}

	return str
}

func wasmLogTopicEncode(name string) string {
	return common.BytesToHash(crypto.Keccak256([]byte(name))).String()
}

func ConvertRlpBytesTo(input []byte, targetType string) interface{} {
	v, ok := Bytes2X_CMD[targetType]
	if !ok {
		panic("unsupported type")
	}

	return reflect.ValueOf(v).Call([]reflect.Value{reflect.ValueOf(input)})[0].Interface()
}

var Bytes2X_CMD = map[string]interface{}{
	"string": byteutil.BytesToString,

	// "uint8":  RlpBytesToUint,
	"uint16": RlpBytesToUint16,
	"uint32": RlpBytesToUint32,
	"uint64": RlpBytesToUint64,

	// "uint8":  RlpBytesToUint,
	"int16": RlpBytesToUint16,
	"int32": RlpBytesToUint32,
	"int64": RlpBytesToUint64,

	"bool": RlpBytesToBool,
}

func RlpBytesToUint16(b []byte) uint16 {
	b = common.LeftPadBytes(b, 32)
	result := common.CallResAsUint32(b)
	return uint16(result)
}

func RlpBytesToUint32(b []byte) uint32 {
	b = common.LeftPadBytes(b, 32)
	return common.CallResAsUint32(b)
}

func RlpBytesToUint64(b []byte) uint64 {
	b = common.LeftPadBytes(b, 32)
	return common.CallResAsUint64(b)
}

func RlpBytesToBool(b []byte) bool {
	if bytes.Compare(b, []byte{1}) == 0 {
		return true
	}
	return false
}

type Proof struct {
	ChainID        string              // 链ID
	Txhash         string              // 跨链交易的哈希
	TxStatus       string              // 跨链交易的执行状态，表示该交易是否执行成功
	ContractStatus string              // 业务合约的执行状态
	BlockNumber    uint64              // 该交易所在的区块号
	TxIndex        int                 // 交易所在区块的下标
	TxKey          string              // 交易维护在跨链状态表中的下标
	Signature      []map[string][]byte // bridge 对proof的签名
	ProofSig       map[string][]byte   // 主bridge收集好的签名
}

type DealProof interface {
	LisenBlock() ([]byte, error)                                   //监听区块
	GetBlock(BlockHeader types.Header) (*types.Block, error)       //获取区块信息
	VerifyTx(Block types.Block) (bool, error)                      //交易验证
	GenNodeProof(block types.Block, Txhash string) ([]byte, error) // 节点生成 Proof
	SignProof(proof []byte) ([]byte, error)                        // Proof签名
	BroadcastProof(proof []byte) error                             //广播proof
	Collectsignatures(proof []byte, nodeId string) ([]byte, error) //签名收集
	GenProof(proof []byte, threshold int) ([]byte, error)          // 封装proof
	SendProof(proof []byte) error                                  // 将Proof发送到链上
	SendTx(tx types.Transaction) error                             //交易路由
	//VerifyProof(pub *PublicKey, hash []byte, r, s *big.Int)
	//(bool, error) //验证签名
}

type NodeTimer interface {
	StartTimer() error                                   // 启动计时器
	StopTimer() error                                    // 停止计时器
	ResetTimer() error                                   // 重置计时器
	StartRaft(msg string, NodeID string)                 // 触发raft重新选主
	RequestProof(NodeId string, msg string) error        //计时器到时，触发重新向节点请求Proof，并收集Proof
	ReSendProof(proof []byte, chainId string) error      // 重新发送Proof
	ReSendTx(Tx types.Transaction, chainID string) error // 重新发送Tx
}

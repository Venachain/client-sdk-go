package common

import (
	"fmt"
	"strconv"

	"github.com/PlatONE_Network/PlatONE-SDK-Go/platone/common/hexutil"
)

const (
	JsonrpcVersion = "2.0"

	SuccessStatus = "1"
)

type Response struct {
	Jsonrpc string    `json:"jsonrpc"`
	Result  string    `json:"result"`
	Id      int       `json:"id"`
	Error   JsonError `json:"error"`
}

type JsonError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (err *JsonError) Error() string {
	if err.Message == "" {
		return fmt.Sprintf("json-rpc error %d", err.Code)
	}
	return err.Message
}

func (err *JsonError) ErrorCode() int {
	return err.Code
}

func ErrorResponse(msg string) *Response {
	return &Response{
		Jsonrpc: JsonrpcVersion,
		Error:   JsonError{Message: msg},
	}
}

func ErrorResponseWithResult(result string, msg string) *Response {
	return &Response{
		Result:  result,
		Jsonrpc: JsonrpcVersion,
		Error:   JsonError{Message: msg},
	}
}

func SuccessResponse(result string) *Response {
	return &Response{
		Jsonrpc: JsonrpcVersion,
		Result:  result,
	}
}

type Receipt struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      int    `json:"id"`
	Result  struct {
		BlockHash         string `json:"blockHash"`
		BlockNumber       string `json:"blockNumber"`
		ContractAddress   string `json:"contractAddress"`
		CumulativeGasUsed string `json:"cumulativeGasUsed"`
		From              string `json:"from"`
		GasUsed           string `json:"gasUsed"`
		Root              string `json:"root"`
		To                string `json:"to"`
		TransactionHash   string `json:"transactionHash"`
		TransactionIndex  string `json:"transactionIndex"`
		Status            string `json:"status"`
	} `json:"result"`
}

func (r *Receipt) ConvertFormatter() {
	if len(r.Result.Status) != 0 {
		i, _ := hexutil.DecodeUint64(r.Result.Status)
		r.Result.Status = strconv.Itoa(int(i))
	}
	if len(r.Result.BlockNumber) != 0 {
		i, _ := hexutil.DecodeUint64(r.Result.BlockNumber)
		r.Result.BlockNumber = strconv.FormatInt(int64(i), 10)
	}
	if len(r.Result.TransactionIndex) != 0 {
		i, _ := hexutil.DecodeUint64(r.Result.TransactionIndex)
		r.Result.TransactionIndex = strconv.FormatInt(int64(i), 10)
	}
	if len(r.Result.GasUsed) != 0 {
		i, _ := hexutil.DecodeUint64(r.Result.GasUsed)
		r.Result.GasUsed = strconv.FormatInt(int64(i), 10)
	}
	if len(r.Result.CumulativeGasUsed) != 0 {
		i, _ := hexutil.DecodeUint64(r.Result.CumulativeGasUsed)
		r.Result.CumulativeGasUsed = strconv.FormatInt(int64(i), 10)
	}
}

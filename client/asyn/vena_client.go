package asyn

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"git-c.i.wxblockchain.com/vena/src/client-sdk-go/client"
	"git-c.i.wxblockchain.com/vena/src/client-sdk-go/log"
	"git-c.i.wxblockchain.com/vena/src/client-sdk-go/packet"
	"git-c.i.wxblockchain.com/vena/src/client-sdk-go/types"
	"github.com/gorilla/websocket"
)

type WsResponseResult struct {
	Jsonrpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  wsResParams `json:"params"`
}

type wsResParams struct {
	Subscription string       `json:"subscription"`
	Result       types.Header `json:"result"`
}

type VenaClient struct {
	RpcClient *client.Client
	WsClient  *WsClient
	Result    chan interface{}
}

// Websocket Client 客户端
type WsClient struct {
	Socket  *websocket.Conn
	Message chan []byte // 监听区块头返回的消息
	//Response chan string // 存储交易hash
}

func NewVenaClient(ctx context.Context, ip string, wsPort uint64, rpcPort uint64, keyfilePath string, passphrase string, buffSize int) (*VenaClient, error) {
	wsClient, err := NewWsClient(ctx, ip, wsPort, buffSize)
	if err != nil {
		log.Error("websocket dial err: ", err)
		return nil, err
	}
	url := client.URL{
		IP:      ip,
		RPCPort: rpcPort,
	}
	rpcClient, err := client.NewClient(ctx, url, keyfilePath, passphrase)
	if err != nil {
		return nil, err
	}
	venaClient := &VenaClient{
		RpcClient: rpcClient,
		WsClient:  wsClient,
	}
	return venaClient, nil
}

func NewVenaClientWithClient(ctx context.Context, wsPort uint64, buffSize int, rpcClient *client.Client) (*VenaClient, error) {
	wsClient, err := NewWsClient(ctx, rpcClient.URL.IP, wsPort, buffSize)
	if err != nil {
		log.Error("websocket dial err: ", err)
		return nil, err
	}
	venaClient := &VenaClient{
		RpcClient: rpcClient,
		WsClient:  wsClient,
		Result:    make(chan interface{}, 1),
	}
	return venaClient, nil
}

func NewWsClient(ctx context.Context, ip string, wsPort uint64, buffSize int) (*WsClient, error) {
	conn, err := DialWS(ctx, ip, wsPort)
	if err != nil {
		log.Error("websocket dial err: ", err)
		return nil, err
	}
	client := &WsClient{
		Socket:  conn,
		Message: make(chan []byte, buffSize),
	}
	return client, nil
}

func DialWS(ctx context.Context, ip string, wsPort uint64) (*websocket.Conn, error) {
	host := fmt.Sprintf("%s:%v", ip, wsPort)
	uri := url.URL{
		Scheme: "ws",
		Host:   host,
	}
	conn, resp, err := websocket.DefaultDialer.DialContext(ctx, uri.String(), nil)
	if err != nil {
		return nil, err
	}
	log.Debug("websocket dial success, response: %+v", resp)
	return conn, nil
}

// 订阅区块头和读取区块头的消息
func (venaClient *VenaClient) SubNewHeads() {
	// 订阅区块头
	message := []byte("{\"jsonrpc\":\"2.0\",\"method\":\"eth_subscribe\", \"params\": [\"newHeads\"],\"id\":\"subscription\"}")
	err := venaClient.WsClient.Socket.WriteMessage(websocket.BinaryMessage, message)
	if err != nil {
		log.Error("error is ", err)
	}
	// 监听订阅到的消息
	venaClient.WsReadMsg()
}

// 读取从ws订阅到到消息
func (venaClient *VenaClient) WsReadMsg() {
	for {
		messageType, message, err := venaClient.WsClient.Socket.ReadMessage()
		if err != nil || messageType == websocket.CloseMessage {
			log.Error("error is ", err)
			break
		}
		log.Info("sub message is ", message)
		venaClient.WsClient.Message <- message
	}
}

func GetBlockHash(blockMessage []byte) (string, error) {
	blockHeader, err := GetBlockHeader(blockMessage)
	if err != nil {
		return "", err
	}
	blockHash := blockHeader.Params.Result.Hash
	return blockHash, nil
}

func GetBlockHeader(msg []byte) (*WsResponseResult, error) {
	var blockHeaderRes WsResponseResult
	if err := json.Unmarshal(msg, &blockHeaderRes); err != nil {
		log.Error("err: ", err)
	}
	return &blockHeaderRes, nil
}

// 通过blockHash查所有的交易
func (venaClient *VenaClient) queryIsTxInBlock(transaction string, blockHash string) (bool, error) {
	if blockHash == "" || transaction == "" {
		return false, nil
	}
	block, err := venaClient.RpcClient.GetBlockByHash(blockHash)
	if err != nil {
		return false, err
	}
	result := queryIsTxInTxs(transaction, block.Transactions)
	return result, nil
}

// 通过blockHash查所有的交易
func (venaClient *VenaClient) queryTxsInBlock(transactions []string, blockHash string) ([]packet.Receipt, error) {
	var result []packet.Receipt
	if blockHash == "" || transactions == nil {
		return nil, nil
	}
	block, err := venaClient.RpcClient.GetBlockByHash(blockHash)
	if err != nil {
		return nil, err
	}
	for _, transaction := range transactions {
		exit := queryIsTxInTxs(transaction, block.Transactions)
		if exit {
			res, err := venaClient.getTxReceipt(transaction)
			if err != nil {
				log.Error("error ", err)
				return nil, err
			}
			result = append(result, *res)
		}
	}
	return result, nil
}

// 通过blockHash查所有的交易
func (venaClient *VenaClient) queryTxInBlock(transaction string, blockHash string) (*packet.Receipt, error) {
	var res *packet.Receipt
	if blockHash == "" || transaction == "" {
		return nil, nil
	}
	block, err := venaClient.RpcClient.GetBlockByHash(blockHash)
	if err != nil {
		return nil, err
	}
	exit := queryIsTxInTxs(transaction, block.Transactions)
	if exit {
		res, err = venaClient.getTxReceipt(transaction)
		if err != nil {
			log.Error("error ", err)
			return nil, err
		}
	}
	return res, nil
}

// 查询交易是否在区块中
func queryIsTxInTxs(transaction string, transactions []string) bool {
	setq := make(map[string]struct{})
	for _, transaction := range transactions {
		setq[transaction] = struct{}{}
	}
	if _, ok := setq[transaction]; ok {
		return true
	} else {
		return false
	}
}

// 获取交易回执
func (venaClient *VenaClient) getTxReceipt(txhash string) (*packet.Receipt, error) {
	return venaClient.RpcClient.GetTransactionReceipt(txhash)
}

// 判断获取到的结果是否是交易的hash
func (venaClient *VenaClient) isTxHash(value interface{}) (bool, string) {
	tmp := value.([]interface{})
	txhash := tmp[0].(string)
	if strings.HasPrefix(txhash, "0x") && len(txhash) == 66 {
		return true, txhash
	}
	return false, ""
}

// 获取receipt，可执行相关函数
func (venaClient *VenaClient) GetResultWithesChan() {
	for {
		select {
		case result := <-venaClient.Result:
			log.Info("receipt is: ", result)
			return
		}
	}
}

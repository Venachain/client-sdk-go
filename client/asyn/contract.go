package asyn

import (
	"context"
	"time"

	rpcClient "git-c.i.wxblockchain.com/vena/src/client-sdk-go/client"
	"git-c.i.wxblockchain.com/vena/src/client-sdk-go/common"
	"git-c.i.wxblockchain.com/vena/src/client-sdk-go/log"
	"git-c.i.wxblockchain.com/vena/src/client-sdk-go/packet"
	"git-c.i.wxblockchain.com/vena/src/client-sdk-go/venachain/keystore"
	"github.com/gorilla/websocket"
)

type AsynContractClient struct {
	RpcContractClient *rpcClient.ContractClient
	WsClient          *WsClient
	Result            chan interface{}
	Txhash            chan string // 存储交易hash
}

// buffSize ：接收交易hash缓存池的大小
func NewAsynContractClient(ctx context.Context, ip string, rpcPort uint64, wsPort uint64, keyfilePath, passphrase, contract, vmType string, buffSize int) (*AsynContractClient, error) {
	err := packet.ParamValid(vmType, "VmType")
	if err != nil {
		return nil, err
	}
	url := rpcClient.NewURL(ip, rpcPort)
	rpcContractClient, err := rpcClient.NewContractClient(ctx, url, keyfilePath, passphrase, contract, vmType)
	if err != nil {
		return nil, err
	}
	wsClient, err := NewWsClient(ctx, ip, wsPort, buffSize)
	if err != nil {
		return nil, err
	}
	asynContractClient := &AsynContractClient{
		rpcContractClient,
		wsClient,
		make(chan interface{}, buffSize),
		make(chan string, 1),
	}
	return asynContractClient, nil
}

// buffSize ：接收交易hash缓存池的大小
func NewAsynContractClientWithKey(ctx context.Context, ip string, rpcPort uint64, wsPort uint64, key *keystore.Key, contract, vmType string, buffSize int) (*AsynContractClient, error) {
	err := packet.ParamValid(vmType, "VmType")
	if err != nil {
		return nil, err
	}
	url := rpcClient.NewURL(ip, rpcPort)
	rpcContractClient, err := rpcClient.NewContractClientWithKey(ctx, url, key, contract, vmType)
	if err != nil {
		return nil, err
	}
	wsClient, err := NewWsClient(ctx, ip, wsPort, buffSize)
	if err != nil {
		return nil, err
	}
	asynContractClient := &AsynContractClient{
		rpcContractClient,
		wsClient,
		make(chan interface{}, buffSize),
		make(chan string, buffSize),
	}
	return asynContractClient, nil
}

// 订阅区块头和读取区块头的消息
func (asynContractClient AsynContractClient) SubNewHeads() {
	// 订阅区块头
	message := []byte("{\"jsonrpc\":\"2.0\",\"method\":\"eth_subscribe\", \"params\": [\"newHeads\"],\"id\":\"subscription\"}")
	err := asynContractClient.WsClient.Socket.WriteMessage(websocket.BinaryMessage, message)
	if err != nil {
		log.Error("error is", err)
	}
	// 监听订阅到的消息
	asynContractClient.WsReadMsg()
}

// 定时去查看sockct 的状态，订阅区块头和读取区块头的消息
func (asynContractClient AsynContractClient) SubNewHeadsWithPing(tryGetClientInterval int64) {
	// 订阅区块头
	message := []byte("{\"jsonrpc\":\"2.0\",\"method\":\"eth_subscribe\", \"params\": [\"newHeads\"],\"id\":\"subscription\"}")
	err := asynContractClient.WsClient.Socket.WriteMessage(websocket.BinaryMessage, message)
	if err != nil {
		log.Error("error is", err)
	}
	ticker := time.NewTicker(time.Duration(tryGetClientInterval) * time.Millisecond)
	for range ticker.C {
		if err := asynContractClient.WsClient.Socket.WriteMessage(websocket.PingMessage, nil); err != nil {
			msg := []byte("CloseMessage")
			asynContractClient.WsClient.Message <- msg
			ticker.Stop()
			return
		}
	}
}

// 读取从ws订阅到到消息
func (asynContractClient AsynContractClient) WsReadMsg() {
	for {
		messageType, message, err := asynContractClient.WsClient.Socket.ReadMessage()
		if err != nil || messageType == websocket.CloseMessage {
			msg := []byte("CloseMessage")
			asynContractClient.WsClient.Message <- msg
			log.Error("error is", err)
			break
		}
		log.Info("sub message is ", string(message))
		asynContractClient.WsClient.Message <- message
	}
}

// 异步获取部署合约的结果,返回的第一个参数是调用call 方法，返回第二个参数为调用send 方法
func (asynContractClient AsynContractClient) DeployAsyncGetReceipt(ctx context.Context, abipath string, codepath string, consParams []string) error {
	// 构造dataGenerator
	dataGenerator, err := asynContractClient.RpcContractClient.MakeDeployGenerator(abipath, codepath, consParams)
	if err != nil {
		return err
	}
	txParams, err := rpcClient.MakeTxparamForDeploy(dataGenerator, &asynContractClient.RpcContractClient.Key.Address)
	if err != nil {
		return err
	}
	return asynContractClient.MessageCallWithAsync(ctx, dataGenerator, *txParams, asynContractClient.RpcContractClient.Key)
}

// execute a method in the contract(evm or wasm)
// contract 可以为合约地址或cns 名字
func (asynContractClient AsynContractClient) ExecuteAsyncGetReceipt(ctx context.Context, funcName string, funcParams []string, contract string) error {
	funcName, funcParams = packet.FuncParse(funcName, funcParams)
	// 构造dataGenerator
	dataGenerator, err := asynContractClient.RpcContractClient.MakeContractGenerator(contract, funcParams, funcName)
	if err != nil {
		return err
	}
	err = asynContractClient.contractCall(ctx, dataGenerator)
	if err != nil {
		return err
	}
	return nil
}

// 封装合约的方法,同步获取receipt
func (asynContractClient AsynContractClient) contractCall(ctx context.Context, dataGenerator *packet.ContractDataGen) error {
	// 构造txparam
	txparam, err := dataGenerator.MakeTxparamForContract(&asynContractClient.RpcContractClient.Key.Address, &dataGenerator.To)
	if err != nil {
		return err
	}
	err = asynContractClient.MessageCallWithAsync(ctx, dataGenerator, *txparam, asynContractClient.RpcContractClient.Key)
	if err != nil {
		return err
	}
	return nil
}

func (asynContractClient AsynContractClient) MessageCallWithAsync(ctx context.Context, dataGen packet.MsgDataGen, tx common.TxParams, key *keystore.Key) error {
	var result = make([]interface{}, 1)
	var err error
	// constant == false 或部署合约的情况
	if dataGen.GetIsWrite() {
		res, err := asynContractClient.RpcContractClient.Send(ctx, &tx, key)
		if err != nil {
			return err
		}
		// 把交易的哈希存到 WsClient的Response中
		//log.Info("sendTransaction hash success,hash is %v",res)
		asynContractClient.Txhash <- res
	} else {
		result, err = asynContractClient.RpcContractClient.Call(dataGen.GetContractDataDen(), &tx)
		if err != nil {
			return err
		}
		asynContractClient.Result <- result
	}
	return nil
}

// 处理sendTransaction 的消息
func (asynContractClient AsynContractClient) GetTxsReceipt() {
	for {
		select {
		case txhash := <-asynContractClient.Txhash:
		Loop:
			for {
				select {
				case block := <-asynContractClient.WsClient.Message:
					blockHash, err := GetBlockHash(block)
					//log.Debug("block hash %v", blockHash)
					if err != nil {
						log.Error("error: ", err)
						return
					}
					//log.Debug("txhash %v", txhash)
					txReceipt, err := asynContractClient.queryTxInBlock(txhash, blockHash)
					if err != nil {
						log.Error("error:", err)
						return
					}
					if txReceipt != nil {
						res := packet.ReceiptParsing(txReceipt, *asynContractClient.RpcContractClient.ContractContent)
						asynContractClient.Result <- res
						break Loop
					}
				}
			}
		}
	}
}

// 通过blockHash查所有的交易
func (asynContractClient AsynContractClient) queryTxInBlock(transaction string, blockHash string) (*packet.Receipt, error) {
	var res *packet.Receipt
	if blockHash == "" || transaction == "" {
		return nil, nil
	}
	block, err := asynContractClient.RpcContractClient.GetBlockByHash(blockHash)
	if err != nil {
		return nil, err
	}
	exit := queryIsTxInTxs(transaction, block.Transactions)
	if exit {
		res, err = asynContractClient.RpcContractClient.GetReceipt(transaction)
		if err != nil {
			log.Error("error ", err)
			return nil, err
		}
	}
	return res, nil
}

// 获取receipt，可执行相关函数
func (asynContractClient AsynContractClient) GetResultWithChan() {
	for {
		select {
		case result := <-asynContractClient.Result:
			log.Info("receipt is: %v", result)
		}
	}
}

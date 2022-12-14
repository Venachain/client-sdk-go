package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Venachain/client-sdk-go/common"
	"github.com/Venachain/client-sdk-go/log"
	"github.com/Venachain/client-sdk-go/packet"
	"github.com/Venachain/client-sdk-go/venachain/common/hexutil"
	"github.com/Venachain/client-sdk-go/venachain/keystore"
)

const (
	sleepTime = 1000000000 // 1 seconds
)

// syn从：true 时会返回交易的receipt，false 时只返回交易hash
func (pc Client) MessageCall(ctx context.Context, dataGen packet.MsgDataGen, tx common.TxParams, key *keystore.Key, sync bool) ([]interface{}, error) {
	var result = make([]interface{}, 1)
	var err error
	if dataGen.GetIsWrite() {
		res, err := pc.Send(ctx, &tx, key)
		if err != nil {
			return nil, err
		}
		if sync {
			polRes, err := pc.GetReceiptByPolling(res)
			if err != nil {
				log.Error("error:%s,you can try get receipt again", err)
				return result, nil
			}
			receiptBytes, err := json.MarshalIndent(polRes, "", "\t")
			if err != nil {
				return nil, err
			}
			log.Info(string(receiptBytes))

			recpt := dataGen.ReceiptParsing(polRes)
			result[0] = recpt.String()
		} else {
			result[0] = res
		}
	} else {
		result, err = pc.Call(dataGen.GetContractDataDen(), &tx)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (pc *Client) Call(dataGen *packet.ContractDataGen, tx *common.TxParams) ([]interface{}, error) {
	var params = make([]interface{}, 0)

	params = append(params, tx)
	params = append(params, "latest")
	action := "eth_call"
	// send the RPC calls
	var resp string
	result, err := pc.RpcClient.Call(context.Background(), action, params...)
	if err != nil {
		return nil, errors.New("send Transaction through http error")
	}
	err = json.Unmarshal(result, &resp)
	if err != nil {
		return nil, err
	}

	outputType := dataGen.GetMethodAbi().Outputs
	return dataGen.ParseNonConstantResponse(resp, outputType), nil
}

func (pc *Client) Send(context context.Context, tx *common.TxParams, key *keystore.Key) (string, error) {
	params, action, err := tx.SendMode(key)
	if err != nil {
		return "", err
	}
	// send the RPC calls
	var resp string
	result, err := pc.RpcClient.Call(context, action, params...)
	if err != nil {
		return "", err
	}

	if err = json.Unmarshal(result, &resp); err != nil {
		return "", err
	}
	if err != nil {
		return "", errors.New("send Transaction through http error")
	}

	return resp, nil
}

func (pc *Client) GetReceiptByPolling(txHash string) (*packet.Receipt, error) {
	ch := make(chan interface{}, 1)
	go pc.getReceiptByPolling(txHash, ch)

	select {
	case receipt := <-ch:
		return receipt.(*packet.Receipt), nil

	case <-time.After(time.Second * 30):
		errStr := fmt.Sprintf("get contract receipt timeout...more than %d second.", 30)
		return nil, errors.New(errStr)
	}
}

// todo: end goroutine?
func (client *Client) getReceiptByPolling(txHash string, ch chan interface{}) {
	var receipt *packet.Receipt
	for {
		var err error
		receipt, err = client.GetTransactionReceipt(txHash)
		// limit the times of the polling
		if err != nil {
			fmt.Println(err.Error())
			time.Sleep(2 * sleepTime)
			continue
		}

		if receipt == nil {
			time.Sleep(2 * sleepTime)
			continue
		}
		ch <- receipt
		break
	}

	//ch <- receipt
	//fmt.Println("***receipt:", receipt)
}

// ============================ Tx Receipt ===================================

func (p *Client) GetTransactionReceipt(txHash string) (*packet.Receipt, error) {

	//var response interface{}
	response, err := p.RpcClient.Call(context.Background(), "eth_getTransactionReceipt", txHash)
	if err != nil {
		return nil, err
	}
	if response == nil {
		return nil, nil
	}
	var resp interface{}
	json.Unmarshal(response, &resp)
	// parse the rpc response
	receipt, err := packet.ParseTxReceipt(resp)
	if err != nil {
		return nil, err
	}

	return receipt, nil
}

// ========================== Sol require/ =============================

func (p *Client) GetRevertMsg(msg *common.TxParams, blockNum uint64) ([]byte, error) {

	var hex = new(hexutil.Bytes)
	res, err := p.RpcClient.Call(context.Background(), "eth_call", msg, hexutil.EncodeUint64(blockNum))
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(res, &hex); err != nil {
		return nil, err
	}
	return *hex, nil
}

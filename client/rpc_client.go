package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/PlatONE_Network/PlatONE-SDK-Go/common"
	"github.com/PlatONE_Network/PlatONE-SDK-Go/packet"
	"github.com/PlatONE_Network/PlatONE-SDK-Go/platone/common/hexutil"
	"github.com/PlatONE_Network/PlatONE-SDK-Go/platone/keystore"
)

const (
	sleepTime = 1000000000 // 1 seconds
)

func (pc Client) MessageCallV2(ctx context.Context, dataGen packet.MsgDataGen, tx common.TxParams, key *keystore.Key, isSync bool) ([]interface{}, error) {
	var result = make([]interface{}, 1)
	var err error

	// combine the data based on the types of the calls (contract call, inner call or deploy call)
	tx.Data, err = dataGen.CombineData()
	if err != nil {
		return nil, errors.New("packet data err: %s\n")
	}

	// 部署合约
	if dataGen.GetIsWrite() {
		res, err := pc.Send(ctx, &tx, key)
		if err != nil {
			return nil, err
		}
		result[0] = res

		if isSync {
			polRes, err := pc.GetReceiptByPolling(res)
			if err != nil {
				return result, nil
			}

			receiptBytes, _ := json.MarshalIndent(polRes, "", "\t")
			fmt.Println(string(receiptBytes))

			recpt := dataGen.ReceiptParsing(polRes)
			// recpt := polRes.Parsing()
			if recpt.Status != packet.TxReceiptSuccessMsg {
				result, _ := pc.GetRevertMsg(&tx, recpt.BlockNumber)
				if len(result) >= 4 {
					recpt.Err, _ = packet.UnpackError(result)
				}
			}

			result[0] = recpt.String()
		}
	} else {
		result, _ = pc.Call(dataGen.GetContractDataDen(), &tx)
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
	err := pc.RpcClient.Call(context.Background(), &resp, action, params...)
	if err != nil {
		return nil, errors.New("send Transaction through http error")
	}

	outputType := dataGen.GetMethodAbi().Outputs
	return dataGen.ParseNonConstantResponse(resp, outputType), nil
}

func (pc *Client) Send(context context.Context, tx *common.TxParams, key *keystore.Key) (string, error) {
	params, action, err := tx.SendModeV2(key)
	if err != nil {
		return "", err
	}

	// send the RPC calls
	var resp string
	err = pc.RpcClient.Call(context, &resp, action, params...)
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

	case <-time.After(time.Second * 10):
		// temp := fmt.Sprintf("\nget contract receipt timeout...more than %d second.\n", 10)
		// return temp + txHash

		errStr := fmt.Sprintf("get contract receipt timeout...more than %d second.", 10)
		return nil, errors.New(errStr)
	}
}

// todo: end goroutine?
func (client *Client) getReceiptByPolling(txHash string, ch chan interface{}) {

	for {
		receipt, err := client.GetTransactionReceipt(txHash)

		// limit the times of the polling
		if err != nil {
			fmt.Println(err.Error())
			fmt.Printf("try again 5s later...")
			time.Sleep(5 * sleepTime)
			fmt.Printf("try again...\n")
			continue
		}

		if receipt == nil {
			time.Sleep(1 * sleepTime)
			continue
		}

		ch <- receipt
	}
}

// ============================ Tx Receipt ===================================

func (p *Client) GetTransactionReceipt(txHash string) (*packet.Receipt, error) {

	var response interface{}
	_ = p.RpcClient.Call(context.Background(), &response, "eth_getTransactionReceipt", txHash)
	if response == nil {
		return nil, nil
	}

	// parse the rpc response
	receipt, err := packet.ParseTxReceipt(response)
	if err != nil {
		return nil, err
	}

	return receipt, nil
}

// ========================== Sol require/ =============================

func (p *Client) GetRevertMsg(msg *common.TxParams, blockNum uint64) ([]byte, error) {

	var hex = new(hexutil.Bytes)
	err := p.RpcClient.Call(context.Background(), hex, "eth_call", msg, hexutil.EncodeUint64(blockNum))
	if err != nil {
		return nil, err
	}

	return *hex, nil
}

package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"git-c.i.wxblockchain.com/PlatONE/src/node/client-sdk-go/platone/common"
	"git-c.i.wxblockchain.com/PlatONE/src/node/client-sdk-go/platone/common/hexutil"
	"git-c.i.wxblockchain.com/PlatONE/src/node/client-sdk-go/platone/keystore"
	"git-c.i.wxblockchain.com/PlatONE/src/node/client-sdk-go/platone/rlp"
	"git-c.i.wxblockchain.com/PlatONE/src/node/client-sdk-go/types"
)

type JsonParam struct {
	Jsonrpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	Id      int         `json:"id"`
}

type TxParams struct {
	From     common.Address  `json:"from"` // the address used to send the transaction
	To       *common.Address `json:"to"`   // the address receives the transactions
	Gas      string          `json:"gas"`
	GasPrice string          `json:"gasPrice"`
	Value    string          `json:"value"`
	Data     string          `json:"data"`
}

func Send(params interface{}, action string, url string) (string, error) {
	param := JsonParam{
		Jsonrpc: "2.0",
		Method:  action,
		Params:  params,
		Id:      1,
	}
	if !(strings.HasPrefix(url, "http")) {
		url = "http://" + url
	}
	resp, err := HttpPost(param, url)
	if err != nil {
		panic(fmt.Sprintf("send http post error .\n %s" + err.Error()))
	}

	return resp, err
}

func HttpPost(param JsonParam, url string) (string, error) {

	client := &http.Client{}
	req, _ := json.Marshal(param)
	reqNew := bytes.NewBuffer(req)

	request, err := http.NewRequest("POST", url, reqNew)
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-type", "application/json")
	response, err := client.Do(request)

	if err != nil {
		return "", err
	}

	if response.StatusCode == 200 {
		body, _ := ioutil.ReadAll(response.Body)
		return string(body), nil
	} else {
		return "", fmt.Errorf("http response status :%s", response.Status)
	}
	return "", err
}

func ParseResponse(r string) (*Response, error) {
	var resp = Response{}
	fmt.Println(r)
	err := json.Unmarshal([]byte(r), &resp)

	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (tx *TxParams) SendModeV2(key *keystore.Key) ([]interface{}, string, error) {
	var action string
	var params = make([]interface{}, 0)

	if key.PrivateKey != nil {
		signedTx, err := tx.GetSignedTx(key)
		if err != nil {
			return nil, "", err
		}

		params = append(params, signedTx)
		action = "eth_sendRawTransaction"
	} else {
		params = append(params, tx)
		action = "eth_sendTransaction"
	}

	return params, action, nil
}

// GetSignedTx gets the signed transaction
func (tx *TxParams) GetSignedTx(keyfile *keystore.Key) (string, error) {

	var txSign *types.Transaction

	// convert the TxParams object to types.Transaction object
	nonce := getNonceRand()
	value, _ := hexutil.DecodeBig(tx.Value)
	gas, _ := hexutil.DecodeUint64(tx.Gas)
	gasPrice, _ := hexutil.DecodeBig(tx.GasPrice)
	data, _ := hexutil.Decode(tx.Data)

	if tx.To == nil {
		txSign = types.NewContractCreation(nonce, value, gas, gasPrice, data)
	} else {
		txSign = types.NewTransaction(nonce, *tx.To, value, gas, gasPrice, data)
	}

	// todo: choose the correct signer
	txSign, _ = types.SignTx(txSign, types.HomesteadSigner{}, keyfile.PrivateKey)
	/// txSign, _ = types.SignTx(txSign, types.NewEIP155Signer(big.NewInt(300)), priv)
	/// utl.Logger.Printf("the signed transaction is %v\n", txSign)

	str, err := rlpEncodeSignedTx(txSign)
	if err != nil {
		return "", err
	}

	return str, nil
}

// RlpEncode encode the input value by RLP and convert the output bytes to hex string
func rlpEncodeSignedTx(val interface{}) (string, error) {

	dataRlp, err := rlp.EncodeToBytes(val)
	if err != nil {
		return "", errors.New("errRlpEncodeFormat")
	}

	return hexutil.Encode(dataRlp), nil

}

// getNonceRand generate a random nonce
// Warning: if the design of the nonce mechanism is modified
// this part should be modified as well
func getNonceRand() uint64 {
	rand.Seed(time.Now().Unix())
	return rand.Uint64()
}

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"

	"github.com/Venachain/client-sdk-go/common"
	"github.com/Venachain/client-sdk-go/types"
	"github.com/Venachain/client-sdk-go/venachain/abi"
	"github.com/Venachain/client-sdk-go/venachain/keystore"
	"github.com/Venachain/client-sdk-go/venachain/rpc"
)

// Client 链 RPC 连接客户端
type Client struct {
	RpcClient *rpc.Client
	Key       *keystore.Key
	URL       *URL
}

type URL struct {
	IP      string
	RPCPort uint64
}

// 传入URL 和keyfil.json 文件路径，密码构建客户端
func NewClient(ctx context.Context, url URL, keyfilePath string, passphrase string) (*Client, error) {
	endpoint := url.GetEndpoint()
	rpcClient, err := rpc.DialContext(ctx, endpoint)
	if err != nil {
		return nil, err
	}
	key, err := NewKey(keyfilePath, passphrase)
	if err != nil {
		return nil, err
	}
	client := &Client{
		RpcClient: rpcClient,
		Key:       key,
		URL:       &url,
	}
	return client, nil
}

// 通过URL和Key 构建客户端
func NewClientWithKey(ctx context.Context, url URL, key *keystore.Key) (*Client, error) {
	endpoint := url.GetEndpoint()
	rpcClient, err := rpc.DialContext(ctx, endpoint)
	if err != nil {
		return nil, err
	}
	client := &Client{
		RpcClient: rpcClient,
		Key:       key,
		URL:       &url,
	}
	return client, nil
}

func NewURL(ip string, port uint64) URL {
	return URL{
		IP:      ip,
		RPCPort: port,
	}
}

func NewKey(KeyfilePath, Passphrase string) (*keystore.Key, error) {
	keyjson, err := ioutil.ReadFile(KeyfilePath)
	if err != nil {
		return nil, err
	}
	return keystore.DecryptKey(keyjson, Passphrase)
}

// rpc 调用通用接口，funcName 为函数的名字，funcParam 为函数参数
// 结果为json.RawMessage格式，需要根据结果的类型使用json.unmarshal() 进行相应的数据类型转换
func (client Client) RpcCall(ctx context.Context, funcName string, funcParam interface{}) (json.RawMessage, error) {
	result, err := client.RpcClient.CallContext(ctx, funcName, funcParam)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func getRpcResult(value json.RawMessage, resultType string) interface{} {
	return abi.BytesConverter(value, resultType)
}

func (c *URL) GetEndpoint() string {
	return fmt.Sprintf("http://%v:%v", c.IP, c.RPCPort)
}

func (c *URL) GetEndpointAddr() string {
	return fmt.Sprintf("%v:%v", c.IP, c.RPCPort)
}

func (client Client) GetBlockByHash(hash string) (*types.GetBlockResponse, error) {
	funcName := types.GetBlockByHash
	result, err := client.RpcClient.Call(context.Background(), funcName, hash, false)
	if err != nil {
		return nil, err
	}
	var res types.GetBlockResponse
	err = json.Unmarshal(result, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (client Client) GetBlockAllByHash(hash string) (*types.Block, error) {
	funcName := types.GetBlockByHash
	raw, err := client.RpcClient.Call(context.Background(), funcName, hash, true)
	if err != nil {
		return nil, err
	}

	var head *types.Header
	var body *types.Body
	if err := json.Unmarshal(raw, &head); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(raw, &body); err != nil {
		return nil, err
	}
	block := types.Block{
		Header:       head,
		Transactions: body.Transactions,
	}
	return &block, nil
}

func (client Client) GetBlockByNumber(blockNumber *big.Int) (*types.GetBlockResponse, error) {
	funcName := types.GetBlockByNumber
	result, err := client.RpcClient.Call(context.Background(), funcName, blockNumber, false)
	if err != nil {
		return nil, err
	}

	var res types.GetBlockResponse
	err = json.Unmarshal(result, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (client Client) GetBlockAllByNumber(blockNumber string) (*types.Block, error) {
	funcName := types.GetBlockByNumber
	raw, err := client.RpcClient.Call(context.Background(), funcName, blockNumber, true)
	if err != nil {
		return nil, err
	}
	var head *types.Header
	var body *types.Body
	if err := json.Unmarshal(raw, &head); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(raw, &body); err != nil {
		return nil, err
	}
	block := types.Block{
		Header:       head,
		Transactions: body.Transactions,
	}
	return &block, nil
}

func (client Client) GetLatestBlock() (*types.Block, error) {
	funcName := types.GetBlockByNumber
	raw, err := client.RpcClient.Call(context.Background(), funcName, "latest", true)
	if err != nil {
		return nil, err
	}
	var head *types.Header
	if err := json.Unmarshal(raw, &head); err != nil {
		return nil, err
	}
	block := types.Block{
		Header: head,
	}
	return &block, nil
}

func (client Client) GetTransactionByHash(txhash string) (*common.TxResponse, error) {
	funcName := types.GetTransactionByHash
	result, err := client.RpcClient.Call(context.Background(), funcName, txhash)
	if err != nil {
		return nil, err
	}
	var res common.TxResponse
	err = json.Unmarshal(result, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (client Client) GetFirstAccount(ctx context.Context) (string, error) {
	funcName := "personal_listAccounts"
	raw, err := client.RpcCall(ctx, funcName, nil)
	if err != nil {
		return "", err
	}
	var addresses []string
	err = json.Unmarshal(raw, &addresses)
	if err != nil {
		return "", err
	}
	return addresses[0], nil
}

func (client Client) UnlockAccount(ctx context.Context, account string, passphrase string) (bool, error) {
	funcName := "personal_unlockAccount"
	var res bool
	result, err := client.RpcClient.CallContext(ctx, funcName, account, passphrase, 0)
	if err != nil {
		return false, err
	}
	if err = json.Unmarshal(result, &res); err != nil {
		return false, err
	}
	return res, nil
}

package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	common_sdk "git-c.i.wxblockchain.com/PlatONE/src/node/client-sdk-go/common"
	"git-c.i.wxblockchain.com/PlatONE/src/node/client-sdk-go/packet"
	"git-c.i.wxblockchain.com/PlatONE/src/node/client-sdk-go/platone/common"
	"git-c.i.wxblockchain.com/PlatONE/src/node/client-sdk-go/platone/keystore"
	"git-c.i.wxblockchain.com/PlatONE/src/node/client-sdk-go/platone/rpc"
	"git-c.i.wxblockchain.com/PlatONE/src/node/client-sdk-go/types"
	"github.com/sirupsen/logrus"
)

// Client 链 RPC 连接客户端
type Client struct {
	RpcClient   *rpc.Client
	Passphrase  string
	KeyfilePath string
	URL         *URL
}

type URL struct {
	IP      string
	RPCPort uint64
}

func NewURL(ip string, port uint64) URL {
	return URL{
		IP:      ip,
		RPCPort: port,
	}
}

func (c *URL) GetEndpoint() string {
	return fmt.Sprintf("http://%v:%v", c.IP, c.RPCPort)
}

func (c *URL) GetEndpointAddr() string {
	return fmt.Sprintf("%v:%v", c.IP, c.RPCPort)
}

// 当Client 的Passphrase 和KeyfilePath 为空时，使用默认客户端的账户和密码进行加密
func (c *Client) ClientSend(action string, params []interface{}) (interface{}, error) {
	var response types.Response
	res, err := common_sdk.Send(params, action, c.URL.GetEndpoint())
	if err != nil {
		logrus.Error("")
		return nil, errors.New("send request is error")
	}
	b := []byte(res)
	err = json.Unmarshal(b, &response)
	if err != nil {
		return nil, errors.New("unmarshal response is error")
	}
	return response, nil
}

// NewClient 创建一个可以操作链的 RPC 客户端
// url 链连接地址
// passphrase 链账户密码
// keyfilePath keyfile 存放的相对地址，默认为 "./keystore"
func NewClient(ctx context.Context, url URL, passphrase, keyfilePath string) (*Client, error) {
	endpoint := url.GetEndpoint()
	rpcClient, err := rpc.DialContext(ctx, endpoint)
	if err != nil {
		return nil, err
	}
	client := &Client{
		RpcClient:   rpcClient,
		Passphrase:  passphrase,
		KeyfilePath: keyfilePath,
		URL:         &url,
	}
	return client, nil
}

// rpcClient 继承了PlatONE RPC 客户端的方法
func (client *Client) GetRpcClient() *rpc.Client {
	return client.RpcClient
}

func (c *Client) RPCSend(ctx context.Context, method string, args ...interface{}) (json.RawMessage, error) {
	return c.RpcClient.CallContext(ctx, method, args)
}

func (pc *Client) clientCommonV2(ctx context.Context, txparam common_sdk.TxParams, dataGen packet.MsgDataGen, to *common.Address, isSync bool) (interface{}, error) {
	// get the client global parameters
	keyjson, err := ioutil.ReadFile(pc.KeyfilePath)
	if err != nil {
		return nil, err
	}
	key, err := keystore.DecryptKey(keyjson, pc.Passphrase)
	if err != nil {
		return nil, err
	}
	rpcClient, err := NewClient(ctx, *pc.URL, pc.Passphrase, pc.KeyfilePath)
	pc.RpcClient = rpcClient.RpcClient
	tx := txparam
	tx.To = to
	if key.Address.String() != "" {
		tx.From = key.Address
	}
	// dataGen == nil 为普通发送交易的逻辑
	if dataGen == nil {
		res, err := pc.Send(ctx, &tx, key)
		if err != nil {
			return nil, errors.New("send transaction is error")
		}
		return res, nil
	}
	// dataGen ！= nil，以下为部署合约的逻辑
	result, err := pc.MessageCallV2(ctx, dataGen, tx, key, isSync)
	if err != nil {
		return nil, err
	}
	return result, nil
}

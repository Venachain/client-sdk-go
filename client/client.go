package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	common_sdk "github.com/PlatONE_Network/PlatONE-SDK-Go/common"
	"github.com/PlatONE_Network/PlatONE-SDK-Go/packet"
	common "github.com/PlatONE_Network/PlatONE-SDK-Go/platone/common"
	"github.com/PlatONE_Network/PlatONE-SDK-Go/platone/keystore"
	"github.com/PlatONE_Network/PlatONE-SDK-Go/platone/rpc"
	"github.com/PlatONE_Network/PlatONE-SDK-Go/types"
	"github.com/sirupsen/logrus"
)

var (
	cliContainer map[string]*Client
	lock         sync.Mutex
)

// Client 链 RPC 连接客户端
type Client struct {
	RpcClient   *rpc.Client
	Passphrase  string
	KeyfilePath string
	URL         *URL
}

// 通过 Ethclient 可以使用以太坊通用接口的方法
type URL struct {
	IP      string
	RPCPort uint64
}

func NewEthClient(ip string, port uint64) URL {
	return URL{
		IP:      ip,
		RPCPort: port,
	}
}

func (c *URL) GetEndpoint() string {
	return fmt.Sprintf("http://%v:%v", c.IP, c.RPCPort)
}

func (c *URL) EthSend(action string, params []interface{}) interface{} {
	var response types.Response
	res, err := common_sdk.Send(params, action, c.GetEndpoint())
	if err != nil {
		logrus.Error("")
		return nil
	}
	b := []byte(res)
	err = json.Unmarshal(b, &response)
	if err != nil {
		return nil
	}
	return response
}

// NewClient 创建一个可以操作链的 RPC 客户端
// url 链连接地址
// passphrase 链账户密码
// keyfilePath keyfile 存放的相对地址，默认为 "./keystore"
func NewClient(ctx context.Context, url string, passphrase string, keyfilePath string) (*Client, error) {
	lock.Lock()
	defer lock.Unlock()
	//if cli, ok := cliContainer[url]; ok {
	//	return cli, nil
	//}
	//ethCli, err := ethclient.DialContext(ctx, url)
	//if err != nil {
	//	return nil, err
	//}
	rpcCli, err := rpc.DialContext(ctx, url)
	if err != nil {
		return nil, err
	}
	client := &Client{
		//ethClient:   ethCli,
		RpcClient:   rpcCli,
		Passphrase:  passphrase,
		KeyfilePath: keyfilePath,
	}
	//cliContainer[url] = client
	return client, nil
}

func (client *Client) GetRpcClient() *rpc.Client {
	return client.RpcClient
}

func (c *Client) RpcSend(result interface{}, method string, args ...interface{}) (interface{}, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	err := c.RpcClient.CallContext(ctx, result, method, args)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func EthClient(action string, params []interface{}, url string) {
	r, err := common_sdk.Send(params, action, url)
	if err != nil {
		fmt.Println("err")
	}
	fmt.Println(r)
}

func (pc *Client) clientCommonV2(txparam common_sdk.TxParams, dataGen packet.MsgDataGen, to *common.Address, isSync bool) ([]interface{}, error) {
	var result = make([]interface{}, 1)
	// get the client global parameters
	keyjson, err := ioutil.ReadFile(pc.KeyfilePath)
	if err != nil {
		return nil, err
	}
	key, err := keystore.DecryptKey(keyjson, pc.Passphrase)
	if err != nil {
		return nil, err
	}
	endpoint := fmt.Sprintf("http://%v:%v", pc.URL.IP, pc.URL.RPCPort)
	rpcClient, err := NewClient(context.Background(), endpoint, pc.Passphrase, pc.KeyfilePath)
	pc.RpcClient = rpcClient.RpcClient
	tx := txparam
	tx.To = to
	if key.Address.String() != "" {
		tx.From = key.Address
	}
	// dataGen == nil 为普通发送交易的逻辑
	if dataGen == nil {
		res, err := pc.Send(&tx, key)
		if err != nil {
			return nil, errors.New("send transaction is error")
		}
		result[0] = res
		return result, nil
	}
	// dataGen ！= nil，以下为部署合约的逻辑
	result[0], err = pc.MessageCallV2(dataGen, tx, key, isSync)
	if err != nil {
		return nil, err
	}
	return result, nil
}

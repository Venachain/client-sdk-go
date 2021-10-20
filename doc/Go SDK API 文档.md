# Go SDK API 文档



```go
// Go SDK Client 接口
type Client struct {
   RpcClient   *rpc.Client
   Passphrase  string
   KeyfilePath string
   URL         *URL
}
```

通用的Client 结构体包括RpcClient，可以继承 PlatONE 中RpcClient 的方法，Passphrase 和 KeyfilePath 为必须参数，URL 为要使用的 PlatONE 的相关 IP 和 RPCPort。 其中提供了NewURL 函数来创建URL。

```go
NewURL(ip string, port uint64) URL
```

```go
type URL struct {
   IP      string
   RPCPort uint64
}
```

可以使用NewClient 方法初始化一个Client。
```go
NewClient(ctx context.Context, url URL, passphrase, keyfilePath string) (*Client, error)
```



## 以太坊通用接口

Client 实现的接口如下：

```go
type IClient interface {
	GetRpcClient() *rpc.Client
	RPCSend(ctx context.Context, result interface{}, method string, args ...interface{}) error
}
```

 其中  ```RPCSend ```  可以使用[以太坊RPC API手册](http://cw.hubwiz.com/card/c/parity-rpc-api/)中的方法，result 定义返回值的类型，只需要构造方法和参数即可。例如以下使用方法：

函数名字为"personal_unlockAccount"， 参数为账户地址，调用解锁账户返回的是一个bool 类型的参数，表示该账户的状态。

```go
func (accountClient AccountClient) UnLock(ctx context.Context) (bool, error) {
   funcName := "personal_unlockAccount"
   funcParams := accountClient.Address.Hex()
   var res bool
   err := accountClient.Client.RPCSend(ctx, &res, funcName, funcParams)
   if err != nil {
      return false, err
   }
   return res, nil
}

```

或者可以参考以下test 的使用
```go
// 通过client rpc 调用rpcsend 方法
func TestRpcSend(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	var addresses []string
	url := URL{
		IP:      "127.0.0.1",
		RPCPort: 6791,
	}
	client, _ := NewClient(ctx, url, "0", "./keystore")
	client.RPCSend(ctx, &addresses, "personal_listAccounts")
	fmt.Println(addresses)
	assert.True(t, addresses != nil)
}
```

如果不确定返回参数的类型，也可以使用EthSend 发送交易。

```go
func TestNewClient(t *testing.T) {
   url := NewURL("127.0.0.1", 6791)
   //client, _ := NewClient(context.Background(), url, "0", "./keystore")
   var param []interface{}
   res, _ := url.EthSend("eth_blockNumber", param)
   fmt.Println(res)
   assert.True(t, res != nil)
}
```



## WebSocket 接口

Websocket 包连接了Client 前端和PlatONE 的websocket 接口。为前端Client 提供订阅区块头和订阅事件的功能。

其相应的代码在`ws`中。

```go
// ws/type.go:21
// Client 单个 websocket 信息
type Client struct {
   Id, Group  string
   LocalAddr  string
   RemoteAddr string
   Path       string
   Socket     *websocket.Conn
   IsAlive    bool
   IsDial     bool
   RetryCnt   int64
   Message    chan []byte
}
```

`ws_subscriber.go `中的 `wsSubscriber` 负责向PlatONE 发送订阅消息。`sub_msg_processor.go` 中的 `SubMsgProcessor` 负责向前端推送消息。

### 使用指南：

`example.go `中使用`gin`框架为后端实现写了一个例子。同时提供前端测试页 `ws_sub_test.html`。首先需要基于PlatONE 的运行环境创建一个 `wsSubscriber` ：

```go
// ws/ws_subscriber.go:31
func newWSSubscriber() *wsSubscriber {
	return &wsSubscriber{
		wsManager: DefaultWebsocketManager,
		ip:        "127.0.0.1",
		port:      26791,
		group: "platone",
	}
}
```

其中 PlatONE-SDK-Go 提供了两种订阅功能：订阅区块头和订阅log。
```go
type Subscription interface {
	SubHeadForChain() error
	SubLogForChain(address, topic string) error
}
```

`SubLogForChain` 需要传递订阅合约的地址和订阅的topic。例如，如果要订阅防火墙打开的事件，可以传入以下参数：

```go
// ws/example.go:113
address := "0x1000000000000000000000000000000000000005"
topic := "0x8cd284134f0437457b5542cb3a7da283d0c38208c497c5b4b005df47719f98a1"
```

**以下使用订阅区块头为例：**

外部 `main` 包中直接运行mian 函数。前端访问 `ws_sub_test.html`测试页 http://127.0.0.1:8888/api/ws/ws_sub_test.html。 

```go
// main.go:9
func main() {
   gin.SetMode(gin.DebugMode)
   gracesRouter := ws.InitRouter()
   err := gracesRouter.Run("127.0.0.1:8888")
   if err != nil {
      logrus.Errorf("test start err: %v", err)
      return
   }
}
```

可在 输入栏中输入 `ping` 查看当前的连接是否成功。如果返回 `pong`，则表示当前连接成功。此时在PlatONE 中发送交易，该订阅的结果会返回到前端页面。

说明：

测试页中 `mygroup1` 为前端的group，可使用不同的group。

```javascript
var host = "ws://127.0.0.1:8888/api/ws/head/mygroup1"
```

如果不使用测试页，可以使用例如 http://coolaf.com/tool/chattest 这样的websocket 测试网页，与 `ws://127.0.0.1:8888/api/ws/head/mygroup1  `建立连接，也可查看到订阅的结果。



## 预编译合约接口

预编译合约不需要传入abi 文件，可以设置为“ ”。

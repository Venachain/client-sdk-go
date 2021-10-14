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





## 预编译合约接口

预编译合约不需要传入abi 文件，可以设置为“ ”。

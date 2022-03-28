# Go SDK API 文档

```go
// Client 链 RPC 连接客户端
type Client struct {
	RpcClient *rpc.Client
	Key       *keystore.Key  //初始化后的key
	URL       *URL
}
```

通用的Client 结构体包括RpcClient，可以继承 PlatONE 中RpcClient 的方法，Passphrase 和 Key 为必须参数，URL 为要使用的 PlatONE 的相关 IP 和 RPCPort。 其中提供了NewURL 函数来创建URL。

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
NewClient(ctx context.Context, url URL, keyfilePath string, passphrase string)
```

也可以通过传入key 来初始化Client。

```go
// 初始化keystore.Key
NewKey(KeyfilePath, Passphrase string) (*keystore.Key, error)
```

```go
// 通过keystore.Key构建Client
NewClientWithKey(ctx context.Context, url URL, key *keystore.Key) (*Client, error) 
```

## 使用方法

由于初始化keystore.Key 时会消耗一部分比较大的内存，因此在使用时建议使用单例模式，通过定义一个全局的DefaultContractClient变量去调用Client 的方法。或者可以调用`NewKey(KeyfilePath, Passphrase string)` 将key定义为全局变量，然后使用`NewClientWithKey`去构建Client。具体的使用方法如下：

### 一. 合约客户端 ContractClient

```go
# 合约客户端数据结构
type ContractClient struct {
	*Client  
	ContractContent *packet.ContractContent
	VmType          string
}
```
初始化合约客户端

```go
// 入参：
// url: 链的IP 和RPC 地址
// keyfilePath: keyfile.json 文件的路径
// passphrase: keyfile 文件对应的密码
// contract: 调用合约的abi文件或预编译合约地址
// vmType: 虚拟机类型
func NewContractClient(ctx context.Context, url URL, keyfilePath, passphrase, contract, vmType string)(*ContractClient, error) {...}
```

1. **定义默认的合约客户端**

```go
var DefaultContractClient *ContractClient
```

2. **初始化默认合约客户端**

使用NewContractClient()初始化合约客户端。

**示例：**

```go
// initContractClient 函数可放到main()函数中调用
// 调用成功后步骤1中的DefaultContractClient会被初始化
func initContractClient() {
	var err error
	contract := "0x0000000000000000000000000000000000000099" //存证合约
	keyfile := "/Users/cxh/go/src/github.com/PlatONE_Network/PlatONE-Go/release/linux/conf/keyfile.json"
	PassPhrase := "0"
	vm := "wasm"
	url := URL{
		IP:      "127.0.0.1",
		RPCPort: 6791,
	}
	DefaultContractClient, err = NewContractClient(context.Background(), url, keyfile, PassPhrase, contract, vm)
	if err != nil {
		log.Error("%s", err)
    return
	}
}
```

#### 1. 调用合约 Execute

调用合约能够实现所有预编译合约合约调用的功能。

使用`DefaultContractClient.Execute()` 调用合约。

```go
// 入参：
// funcName: 合约调用的函数名
// funcParams： 合约调用的参数
// contract：合约地址或cns名字
// 输出：上链的receipt结果
func Execute(ctx context.Context, funcName string, funcParams []string, contract string) (interface{}, error){...}
```

**SaveEvidence 示例**

```go
func SaveEvidence(key string,value string) {
	funcname := "saveEvidence"
	funcparam := []string{}
	funcparam = append(funcparam, key)
	funcparam = append(funcparam, value)
	result, err := DefaultContractClient.Execute(context.Background(), funcname, funcparam, "0x0000000000000000000000000000000000000099")
	if err != nil {
		log.Error("%s", err)
		return
	}
	log.Info("result:%v", result)
	DefaultContractClient.Client.RpcClient.Close()
}
```

**GetEvidence 示例**

```go
func GetEvidence(key string) {
	funcname := "getEvidence"
	funcparam := []string{}
	funcparam = append(funcparam, key)
	result, err := DefaultContractClient.Execute(context.Background(), funcname, funcparam, "0x0000000000000000000000000000000000000099")
	if err != nil {
		log.Error("%s", err)
		return
	}
	log.Info("result:%v", result)
	DefaultContractClient.Client.RpcClient.Close()
}
```
**VerifyProofByRange 示例**

```go
func VerifyProofByRange(userid, proof, pid, rang string) {
	funcname := "verifyProofByRange"
	funcparam := []string{}
	funcparam = append(funcparam, userid)
	funcparam = append(funcparam, proof)
	funcparam = append(funcparam, pid)
	funcparam = append(funcparam, rang)
	result, err := DefaultContractClient.Execute(context.Background(), funcname, funcparam, "0x0000000000000000000000000000000000000100")
	if err != nil {
		log.Error("%s", err)
		return
	}
	log.Info("result:%v", result)
	DefaultContractClient.Client.RpcClient.Close()
}
```

**BpGetResult 示例**

```go
func VerifyProofByRange(pid string) {
	funcname := "getResult"
	funcparam := []string{}
	funcparam = append(funcparam, pid)
	result, err := DefaultContractClient.Execute(context.Background(), funcname, funcparam, "0x0000000000000000000000000000000000000100")
	if err != nil {
		log.Error("%s", err)
		return
	}
	log.Info("result:%v", result)
	DefaultContractClient.Client.RpcClient.Close()
}
```

#### 2. 部署合约 Deploy

##### 2.1 同步获取reciept

使用`DefaultContractClient.Deploy()` 调用合约。
```go
// 入参：
// abipath: 部署合约的abi 文件路径
// codepath： 部署合约的code 文件路径，比如wasm合约以.wasm结尾的文件的绝对路径
// consParams：solidity 合约的contractor参数
// 输出：上链的receipt结果
func Deploy(ctx context.Context, abipath string, codepath string, consParams []string) (interface{}, error){...}
```

部署合约需要使用codePath 和 ContractClient中的AbiPath (此时初始化AbiPath 时不能设置为空)。 consParams 为合约的某个参数。执行该测试函数后，成功部署了一个合约。得到返回的事件和合约地址：

```go
// 同步过去部署结果的示例
func TestContractClient_Deploy(t *testing.T) {
	codePath := "/Users/cxh/Downloads/example/example.wasm"
	abiPath := "/Users/cxh/Downloads/example/example.cpp.abi.json"
	contract, err := InitContractClient("")
	if err != nil {
		log.Error("error is:%v", err)
	}
	var consParams []string
	result, err := contract.Deploy(context.Background(), abiPath, codePath, consParams)
	if err != nil {
		log.Error("error:%v", err)
	}
	log.Info("result:%v", result)
	assert.True(t, result != nil)
}
```

```json
// 部署合约成功的结果
{
	"status": "Operation Succeeded",
	"contractAddress": "0x35853e5643104cd96bd4590f5d4466c577786cfe",
	"logs": [
		"Event init: init success... "
	],
	"blockNumber": 198,
	"GasUsed": 1911684,
	"From": "0x3fcaa0a86dfbbe105c7ed73ca505c7a59c579667",
	"To": "",
	"TxHash": "0xacdc551f2539068eab227f112ebaeade286d75852aa27e63ecb53176489bee3f"
}
```

##### 2.1 异步获取reciept





## 以太坊通用接口

Client 实现的接口如下：

```go
// client/types.go:16
type IClient interface {
	GetRpcClient() *rpc.Client
	RPCSend(ctx context.Context, result interface{}, method string, args ...interface{}) error
}
```

 其中  ```RPCSend ```  可以使用[以太坊RPC API手册](http://cw.hubwiz.com/card/c/parity-rpc-api/)中的方法。例如以下使用方法：

该方法表示解锁一个账户，函数名字为"personal_unlockAccount"， 参数为账户地址，调用解锁账户返回的是一个bool 类型的参数，表示该账户的状态。

```go
// client/account.go:105
func (accountClient AccountClient) UnLock(ctx context.Context) (bool, error) {
	funcName := "personal_unlockAccount"
	funcParams := accountClient.Address.Hex()
	var res bool
	result, err := accountClient.Client.RPCSend(ctx, funcName, funcParams)
	if err != nil {
		return false, err
	}
	if err = json.Unmarshal(result, &res); err != nil {
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

如果不确定返回参数的类型，也可以使用ClientSend 发送交易。

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
// 相关的ws/ws_subscriber.go:39
func NewWSSubscriber(ip string, port int64, group string) *WsSubscriber {
	return &WsSubscriber{
		wsManager: DefaultWebsocketManager,
		ip:        ip,
		port:      port,
		group:     group,
	}
}
```
使用方法，使用NewWSSubscriber方法定义一个全局DefaultWSSubscriber变量，后续调用DefaultWSSubscriber 
```go
// 定义全局WS订阅对象
var (
	// DefaultWSSubscriber 默认的 websocket 订阅器
	DefaultWSSubscriber *ws.WsSubscriber
)

// 使用NewWSSubscriber创建WS全局对象
DefaultWSSubscriber = ws.NewWSSubscriber("127.0.0.1",26791,"platone")

// 使用 DefaultWebsocketManager 订阅Log事件
DefaultWebsocketManager.WsClientForLog

// 使用 DefaultWebsocketManager 订阅NewHeads事件
DefaultWebsocketManager.WsClientForNewHeads
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

## 合约部署和合约调用接口

合约客户端定义如下：

```go
// client/contract.go:13
type ContractClient struct {
   *Client
   AbiPath string
   VmType      string
}
```

实现的接口如下：

```go
// client/types.go:22
type IContract interface {
   Deploy(ctx context.Context, txparam common.TxParams, codepath string, consParams []string) ([]interface{}, error)
   ListContractMethods() (packet.ContractContent, error)
   Execute(ctx context.Context, txparam common.TxParams, funcName string, funcParams []string, address string) ([]interface{}, error)
   IsFuncNameInContract(funcName string) (bool, error)
   GetReceipt(txhash string) (*packet.Receipt, error)
}
```

### 合约部署

在` client/contract_test.go `中的`initContractClient()` 展示了初始化合约客户端的例子：

以下测试用例展示了如何部署合约：

```
func TestContractClient_Deploy(t *testing.T) {
   codePath := "/Users/cxh/Downloads/example/example.wasm"
   txparam, contract := InitContractClient()
   var consParams []string
   result, _ := contract.Deploy(context.Background(), txparam, codePath, consParams)
   fmt.Println(result)
   assert.True(t, result != nil)
}
```



### 展示合约方法

 ListContractMethods()  可以展示合约的所有方法。

```go
// client/contract_test.go:48
func TestContractClient_ListContractMethods(t *testing.T) {
	_, contract := InitContractClient()
	result, _ := contract.ListContractMethods()
	fmt.Println(result.ListAbiFuncName())
	assert.True(t, result != nil)
}
```

显示该合约的所有方法为：

```shell
function: init()
function: setEvidence(key string,msg string)
function: deleteEvidence(key string)
function: getEvidence(key string) string
event: setName( string)
event: init( string)
```

### 合约调用

上传合约数据：

```go
// client/contract_test.go:55
func TestContractClient_Execute(t *testing.T) {
   txparam, contract := InitContractClient()
   funcname := "setEvidence"
   funcparam := []string{"1", "data"}
   addr := "0x35853e5643104cd96bd4590f5d4466c577786cfe"
   result, _ := contract.Execute(context.Background(), txparam, funcname, funcparam, addr)
   assert.True(t, result != nil)
}
```

存证合约传入的参数是`"data"`，此时可以得到该交易的receipt。合约调用的核心是需要知道所要调用合约的函数，根据函数需要的input 类型构造函数的传入参数。

此外，还可以通过cns 调用合约，此时Execute 传入的最后一个参数为cns 名字。例如：

```go
// client/contract_test.go:74
func TestContractClient_CnsExecute(t *testing.T) {
   txparam, contract := InitContractClient()
   funcname := "setEvidence"
   funcparam := []string{"1", "23"}
   cns := "wxbc1"
   result, _ := contract.Execute(context.Background(), txparam, funcname, funcparam, cns)
   fmt.Println(result)
   assert.True(t, result != nil)
}
```

### 查询合约方法是否属于该合约

```go
// client/contract_test.go:84
func TestContractClient_IsFuncNameInContract(t *testing.T) {
   _, contract := InitContractClient()
   funcname := "setEvidence"
   result, _ := contract.IsFuncNameInContract(funcname)
   fmt.Println(result)
   assert.True(t, result != false)
}
```

返回 true 则表示该合约中存在该方法。

### 根据交易hash 查询交易的receipt

```go
GetReceipt(txhash string) (*packet.Receipt, error)
```

```go
// client/contract_test.go:64 
func TestContractClient_GetReceipt(t *testing.T) {
   txhash := "0x35972a847e8c29148976e8a1884665732c862706c71bbaaf573e8cbd432ba921"
   _, contractClient := InitContractClient()
   result, _ := contractClient.GetReceipt(txhash)
   if result != nil {
      resultBytes, _ := json.MarshalIndent(result, "", "\t")
      fmt.Printf("result:\n%s\n", resultBytes)
   }
   assert.True(t, result != nil)
}
```

## 预编译合约接口

因为预编译合约需要调用合约，因此各个预编译合约的结构体都需要包含`ContractClient`。 并且预编译合约不需要传入合约的 abi 文件，因此在对各个合约初始化时，可以将contract 的AbiPath 设置为空。

### 账户合约 Account

账户客户端的数据结构如下，其中需要包括账户地址Address。

```go
// client/account.go:14
type AccountClient struct {
   ContractClient
   Address common.Address
}
```

账户的接口分别包括新增用户，更新用户，用户查询                                                                                                                                                                 

```go
// client/types.go:30
type IAccount interface {
   UserAdd(ctx context.Context, txparam common_sdk.TxParams, name, phone, email, organization string) (string, error)
   UserUpdate(ctx context.Context, txparam common_sdk.TxParams, phone, email, organization string) (string, error)
   QueryUser(ctx context.Context, txparam common_sdk.TxParams, user string) (string, error)
   Lock(ctx context.Context) (bool, error)
   UnLock(ctx context.Context) (bool, error)
}
```

#### 新增账户

以下为账户合约新增账户的示例：

```go
// client/account_test.go:15
// 如果没有abipath 和codepath 的话，可以设置为空
func InitAccountClient() (common_sdk.TxParams, AccountClient) {
   txparam, contract := InitContractClient()
   contract.AbiPath = ""
   client := AccountClient{
      ContractClient: contract,
      Address:        common.HexToAddress("3fcaa0a86dfbbe105c7ed73ca505c7a59c579667"),
   }
   return txparam, client
}
```

首先需要初始化账户客户端，以下展示添加账户的例子，需要的参数分别代表name, phone, email, organization。如果不需要可以设置为“ ”。

```go
// client/account_test.go:25
func TestAccountClient_UserAdd(t *testing.T) {
   txparam, client := InitAccountClient()
   result, _ := client.UserAdd(context.Background(), txparam, "Alice", "110", "", "")
   fmt.Println(result)
   assert.True(t, result != "")
}
```

#### 更新账户

传入需要更新的账户phone, email, organization。注意只能更新这些信息，账户的名字是不能更改的。

```go
// client/account_test.go:32
func TestAccountClient_UserUpdate(t *testing.T) {
   txparam, client := InitAccountClient()
   result, _ := client.UserUpdate(context.Background(), txparam, "13556672653", "test@163.com", "wxbc2")
   fmt.Println(result)
   assert.True(t, result != "")
}
```

#### 查询账户信息

传入账户的名字即可查询到该账户的相关信息

```go
// func TestAccountClient_QueryUser(t *testing.T) {
func TestAccountClient_QueryUser(t *testing.T) {
   txparam, client := InitAccountClient()
   result, _ := client.QueryUser(context.Background(), txparam, "Alice")
   fmt.Println(result)
   assert.True(t, result != "")
}
```

#### 账户锁定和解锁

在账户客户端，sdk 还提供了账户锁定和解锁的功能。调用以下函数即可：

```go
client.UnLock(context.Background())
client.Lock(context.Background())
```

### 剩下的预编译合约

#### CnsClient

 需要指定cns 的名字

```go
type CnsClient struct {
   ContractClient
   name string
}
```
#### FireWallClient

防火墙客户端管理合约的防火墙，需要指定合约地址。

```go
type FireWallClient struct {
   ContractClient
   ContractAddress string
}
```

#### NodeClient

节点客户端需要指定节点的名字

```go
type NodeClient struct {
   ContractClient
   NodeName string
}
```

SysConfigClient 和 RoleClient 则使用合约的方法。

剩下的预编译合约相关接口如下，相关的test 已在client 包中。可查看使用：

```go
type ICns interface {
   CnsExecute(ctx context.Context, txparam common.TxParams, funcName string, funcParams []string, cns string) ([]interface{}, error)
   CnsRegister(ctx context.Context, txparam common_sdk.TxParams, version, address string) (string, error)
   CnsResolve(ctx context.Context, txparam common_sdk.TxParams, version string) (string, error)
   CnsRedirect(ctx context.Context, txparam common_sdk.TxParams, version string) (string, error)
   CnsQueryAll(ctx context.Context, txparam common_sdk.TxParams) (string, error)
   CnsQueryByName(ctx context.Context, txparam common_sdk.TxParams) (string, error)
   CnsQueryByAddress(ctx context.Context, txparam common_sdk.TxParams, address string) (string, error)
   CnsQueryByAccount(ctx context.Context, txparam common_sdk.TxParams, account string) (string, error)
   CnsStateByAddress(ctx context.Context, txparam common_sdk.TxParams, address string) (int32, error)
   CnsState(ctx context.Context, txparam common_sdk.TxParams) (int32, error)
}

type IFireWall interface {
   FwStatus(ctx context.Context, txparam common_sdk.TxParams) (string, error)
   FwStart(ctx context.Context, txparam common_sdk.TxParams) (string, error)
   FwClose(ctx context.Context, txparam common_sdk.TxParams) (string, error)
   FwExport(ctx context.Context, txparam common_sdk.TxParams, filePath string) (bool, error)
   FwImport(ctx context.Context, txparam common_sdk.TxParams, filePath string) (string, error)
   FwNew(ctx context.Context, txparam common_sdk.TxParams, action, targetAddr, api string) (string, error)
   FwDelete(ctx context.Context, txparam common_sdk.TxParams, action, targetAddr, api string) (string, error)
   FwReset(ctx context.Context, txparam common_sdk.TxParams, action, targetAddr, api string) (string, error)
   FwClear(ctx context.Context, txparam common_sdk.TxParams, action string) (string, error)
}

type INode interface {
   NodeAdd(ctx context.Context, txparam common_sdk.TxParams, requestNodeInfo syscontracts.NodeInfo) (string, error)
   NodeDelete(ctx context.Context, txparam common_sdk.TxParams) (string, error)
   NodeUpdate(ctx context.Context, txparam common_sdk.TxParams, request syscontracts.NodeUpdateInfo) (string, error)
   NodeQuery(ctx context.Context, txparam common_sdk.TxParams, request *syscontracts.NodeQueryInfo) (string, error)
   NodeStat(ctx context.Context, txparam common_sdk.TxParams, request *syscontracts.NodeStatInfo) (int32, error)
}

type IRole interface {
   SetSuperAdmin(ctx context.Context, txparam common_sdk.TxParams) (string, error)
   TransferSuperAdmin(ctx context.Context, txparam common_sdk.TxParams, address string) (string, error)
   AddChainAdmin(ctx context.Context, txparam common_sdk.TxParams, address string) (string, error)
   DelChainAdmin(ctx context.Context, txparam common_sdk.TxParams, address string) (string, error)
   AddGroupAdmin(ctx context.Context, txparam common_sdk.TxParams, address string) (string, error)
   DelGroupAdmin(ctx context.Context, txparam common_sdk.TxParams, address string) (string, error)
   AddNodeAdmin(ctx context.Context, txparam common_sdk.TxParams, address string) (string, error)
   DelNodeAdmin(ctx context.Context, txparam common_sdk.TxParams, address string) (string, error)
   AddContractAdmin(ctx context.Context, txparam common_sdk.TxParams, address string) (string, error)
   DelContractAdmin(ctx context.Context, txparam common_sdk.TxParams, address string) (string, error)
   AddContractDeployer(ctx context.Context, txparam common_sdk.TxParams, address string) (string, error)
   DelContractDeployer(ctx context.Context, txparam common_sdk.TxParams, address string) (string, error)
   GetAddrListOfRole(ctx context.Context, txparam common_sdk.TxParams, role string) (string, error)
   GetRoles(ctx context.Context, txparam common_sdk.TxParams, address string) (string, error)
   HasRole(ctx context.Context, txparam common_sdk.TxParams, address, role string) (int32, error)
}

type ISysconfig interface {
   SetSysConfig(ctx context.Context, txparam common_sdk.TxParams, request SysConfigParam) ([]string, error)
   GetTxGasLimit(ctx context.Context, txparam common_sdk.TxParams) (uint64, error)
   GetBlockGasLimit(ctx context.Context, txparam common_sdk.TxParams) (uint64, error)
   GetGasContractName(ctx context.Context, txparam common_sdk.TxParams) (string, error)
   GetIsProduceEmptyBlock(ctx context.Context, txparam common_sdk.TxParams) (uint32, error)
   GetCheckContractDeployPermission(ctx context.Context, txparam common_sdk.TxParams) (uint32, error)
   GetAllowAnyAccountDeployContract(ctx context.Context, txparam common_sdk.TxParams) (uint32, error)
   GetIsApproveDeployedContract(ctx context.Context, txparam common_sdk.TxParams) (uint32, error)
   GetIsTxUseGas(ctx context.Context, txparam common_sdk.TxParams) (uint32, error)
   GetVRFParams(ctx context.Context, txparam common_sdk.TxParams) (string, error)
}
```

相关预编译合约的使用可以查看以下文档：

https://git-c.i.wxblockchain.com/PlatONE/doc/Dev/blob/develop/PlatONE%20%E6%96%87%E6%A1%A3%E9%9B%86%E5%90%88/platoneCli/%E9%93%BE%E4%BA%A4%E4%BA%92%E5%B7%A5%E5%85%B7%20platonecli%20%E6%93%8D%E4%BD%9C%E6%8C%87%E5%8D%97.md


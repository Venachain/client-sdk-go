package asyn

import (
	"context"
	"fmt"
	"testing"
	"time"

	"git-c.i.wxblockchain.com/PlatONE/src/node/client-sdk-go/client"
	"git-c.i.wxblockchain.com/PlatONE/src/node/client-sdk-go/log"
)

var (
	DefaultVenaContractClient *AsynContractClient
)

func init() {
	var err error
	wsPort := uint64(26791)
	keyfile := "/Users/cxh/go/src/VenaChain/venachain/release/linux/conf/keyfile.json"
	PassPhrase := "0"
	key, _ := client.NewKey(keyfile, PassPhrase)
	DefaultVenaContractClient, err = NewAsynContractClientWithKey(context.Background(), "127.0.0.1", uint64(6791), wsPort, key, "0x0000000000000000000000000000000000000100", "wasm", 100)
	if err != nil {
		log.Error("err%v", err)
	}
}

func TestExampleClientSubscription(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	keyfile := "/Users/cxh/go/src/VenaChain/venachain/release/linux/conf/keyfile.json"
	PassPhrase := "0"
	url := client.URL{
		IP:      "127.0.0.1",
		RPCPort: 6791,
	}
	rpcClient, err := client.NewClient(ctx, url, keyfile, PassPhrase)
	if err != nil {
		fmt.Println(err)
	}
	venaClient, err := NewVenaClientWithClient(ctx, 26791, 100, rpcClient)
	if err != nil {
		fmt.Println(err)
	}
	venaClient.SubNewHeads()
	go venaClient.WsReadMsg()
	time.Sleep(300 * time.Second)
}

func TestContractClient_DeployAsyncGetReceipt(t *testing.T) {
	codePath := "/Users/cxh/Downloads/example/example.wasm"
	abiPath := "/Users/cxh/Downloads/example/example.cpp.abi.json"
	// 在wsClient 的Message 中存储收到的区块头，监听区块
	go DefaultVenaContractClient.SubNewHeads()
	// 程序执行后关闭socket 消息
	defer DefaultVenaContractClient.WsClient.Socket.Close()
	// 处理receipt,获取结果
	go DefaultVenaContractClient.GetTxsReceipt()
	// receipt 结果放在result channel中，以下为输出result的示例
	go DefaultVenaContractClient.GetResultWithChan()
	var consParams []string
	err := DefaultVenaContractClient.DeployAsyncGetReceipt(context.Background(), abiPath, codePath, consParams)
	if err != nil {
		log.Error("error:%v", err)
	}
	time.Sleep(time.Second)

	err = DefaultVenaContractClient.DeployAsyncGetReceipt(context.Background(), abiPath, codePath, consParams)
	if err != nil {
		log.Error("error:%v", err)
	}
	time.Sleep(300 * time.Second)
}

func TestContractClient_bpAsynGetResult(t *testing.T) {
	funcname := "verifyProofByRange"
	funcparam := []string{}
	funcparam = append(funcparam, "cx1h")
	funcparam = append(funcparam, "0xf904bbb88025688c60d9d76c78f49250f23e737d80c2837c47550a94bee565b9"+
		"52385edb441904f32b67f0af1e7142c2ed5a3eab789315d6f204d6900a86edb01e7ab3efde0d8e955fcfab1396f034db0491c5"+
		"a3c0697020c664af95f930f9bbaf283d095c204ace0000efb16f22243a17c3f19bd8232b8f98e2847ba01a9a6ceb54131eceb840"+
		"1da555a486e651f3fba4e44a0ff5a258ac6cc2369db60e5227e194838dcf0c5803487277abd6d8743ea38835ee8a49da62949d3e"+
		"f9bf798519b5ee8a4d6f9f99b8400047a403177126bc2ebe7265605b4e5f83896bfca9381799a77866d054fe9bcc242b2268934ba"+
		"453070998b4306b81588d70bbf8fd925c8251ce1fb4bc5c094cb8402120f085d1942469db49c01a88a6b6e88ec28bb4af1b7177cdaf"+
		"d3ac425507c0132c3f7666aacc69823d04005e8157ff443784a6ef5df6ecf70f912243172464b84028387ff6ec87cbcf832b5e06e35d"+
		"38524648975aea037b98600571620ef1d36b06c4f05fee29ab1f4f3558e620e8cbdab315e74d2a8e071db247053395d318eda002a725b"+
		"67c2cbafe59d65f17c9075adf64f9c08e013ef8f44c3006227c061567a010a3b4503968041aadd2d374ab3985194c48d0793537dfe5503"+
		"502294209dd2fa0075cac3db9e205d37c92e68f73d943fa1c6fc7d387716e7ac02ac08409c48fadb902cbf902c8a00f1d852e0698d25ec"+
		"1b5e559d8294b69978314a9d1734993dcf283c0ac7629efa018e347046628a0b44ad4dc00b65dc142eb9af7b5b69335dcd0323f2527c714"+
		"28b9014010456b0ec52c26f06b4da2b0b1c5c083e3f1bbe4dd0dfd59a28f3a5e8023ea2d24ed00cc8da561cb1afa97806768a4e3e1a1f0415"+
		"792433899ebe8c14936dfe12476eadebf3e96ff3b0ba57fdc2d72ff8a21d64c4d9c700b00da4e8616d7ea1f0ab43c5db055b9d70bc07167f13f"+
		"254927a6ac1ccb88d32df8967db35f1c056b107ecc5fd51698d099a1e4cb34286026a0bebcc844f23d3dc1fc7de9c5872c2a09418435610daf6dcd"+
		"605c6f02785c73327863b1ac04a7b4fe7845cf9e58a1412387c3d89e9b05a6c3af3214b44e1a0585c288c49f3269fde8622f0b9bfd29ee1f27d63e51"+
		"06657a3ce7f135d8c006eec562711d44703fdb82ce2091e85d537413c18fc98174b3b4c161455d25d2c75b5f724d01fefa1ca9587d939b49b6275c2"+
		"666124c70ebefa9b7d8592efcbcae7883d957a709200f7dfd470f41dbe58adbb901400f7cd8a1e22193876a3e795715a683d142b5ef9f0e38b97bfde"+
		"43fc2422d1fb210b5fa73a9b92dd8cf64da005fd4d211d9dd7e80bb1eb1e20002a8819b1f346915430694b34e4462cf9c72cbb58e4f0131f4f9fd1dc6"+
		"2ff1248b4ce06c953937272067cdf5143efd05db89ce34999ccbc91644266ea00939f6d04be1fdb638b21cb5cb6a0ca613178fad7a3692c62476129d7"+
		"0d68df659b86dffccb6a24ac95230322142e5c218296ffc0e367a70250c141347adff292b7e2f210000f2aa78e22e3ba8a81482fb3f14a9763f04c269"+
		"487201d394d270376b898ef6a0cd82f90a0e6a4695c44bdf8b8f30ea9c041be6e9293a8d43d72c005dad7b3ea9c6d02db00ce227ab022cdb072a23880df"+
		"fa64816517dfbb3b8bc2fc8b160f38a57e2a06d0859e369ec6f41fbea2556ef885eb1aebc5d960d415e328dfb0d330f6f9afdd6")
	funcparam = append(funcparam, "121")
	funcparam = append(funcparam, "test")
	go DefaultVenaContractClient.SubNewHeads()
	// 程序执行后关闭socket 消息
	defer DefaultVenaContractClient.WsClient.Socket.Close()
	// 处理receipt,获取结果,结果存储在AsynContractClient的result中
	go DefaultVenaContractClient.GetTxsReceipt()
	// receipt 结果放在result channel中，以下为输出result的示例
	go DefaultVenaContractClient.GetResultWithChan()
	err := DefaultVenaContractClient.ExecuteAsyncGetReceipt(context.Background(), funcname, funcparam, "0x0000000000000000000000000000000000000100")
	if err != nil {
		log.Error("error:%v", err)
	}
	time.Sleep(time.Second)
	err = DefaultVenaContractClient.ExecuteAsyncGetReceipt(context.Background(), funcname, funcparam, "0x0000000000000000000000000000000000000100")
	if err != nil {
		log.Error("error:%v", err)
	}
	time.Sleep(300 * time.Second)

}

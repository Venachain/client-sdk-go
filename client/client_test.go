package client

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

//// 通过ethclient 调用，拿到的是json response
//func TestNewClient(t *testing.T) {
//	url := NewURL("127.0.0.1", 6791)
//	//client, _ := NewClient(context.Background(), url, "0", "./keystore")
//	var param []interface{}
//	res, _ := url.EthSend("eth_blockNumber", param)
//	fmt.Println(res)
//	assert.True(t, res != nil)
//}

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

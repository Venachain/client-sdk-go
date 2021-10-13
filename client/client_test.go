package client

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 通过ethclient 调用，拿到的是json response
func TestNewClient(t *testing.T) {
	c := NewEthClient("127.0.0.1", 6791)
	var param []interface{}
	res := c.EthSend("eth_blockNumber", param)
	fmt.Println(res)
	assert.True(t, res != nil)
}

// 通过client rpc 调用rpcsend 方法
func TestRpcSend(t *testing.T) {
	ctx := context.Background()
	var addresses []string
	client, _ := NewClient(ctx, "http://127.0.0.1:6791", "0", "./keystore")
	client.RpcSend(&addresses, "personal_listAccounts")
	fmt.Println(addresses)
	assert.True(t, addresses != nil)
}

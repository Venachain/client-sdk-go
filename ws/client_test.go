package ws

import (
	"testing"
	"time"
)

func TestManager_WsClient(t *testing.T) {
	InitWsSubscriber("127.0.0.1", 26791, "venachain")
	DefaultWSSubscriber.SubHeadForChain()
	time.Sleep(time.Second * 100)
}

func TestManager_SubLogForChain(t *testing.T) {
	address := "0x1000000000000000000000000000000000000005"
	topic := "0x8cd284134f0437457b5542cb3a7da283d0c38208c497c5b4b005df47719f98a1"
	go DefaultWSSubscriber.SubLogForChain(address, topic)
	time.Sleep(time.Second * 10000)
}

package common

import (
	"encoding/json"
	"testing"
)

var (
	req = Request{
		Contract:  "0xb8682ed9127cdf25baa9f8c5df36f7c4dc37ba40",
		ChannelID: "001",
		Method:    "set",
		Params:    []string{"1"},
	}

	reqJson = "{\"channelID\":\"001\",\"contract\":\"0xb8682ed9127cdf25baa9f8c5df36f7c4dc37ba40\",\"method\":\"set\",\"params\":[\"1\"]}"

	xaReq = XARequest{
		XAID:     "xxxx01",
		Requests: []*Request{&req},
	}

	xaReqJson = "{\"xaID\":\"xxxx01\",\"requests\":[{\"channelID\":\"001\",\"contract\":\"0xb8682ed9127cdf25baa9f8c5df36f7c4dc37ba40\",\"method\":\"set\",\"params\":[\"1\"]}]}"
)

func TestRequestJson(t *testing.T) {

	data, err := json.Marshal(req)
	if err != nil {
		t.Error(err)
	}
	if string(data) != reqJson {
		t.Error("request marshal json err")
	}
}

func TestXARequestJson(t *testing.T) {

	data, err := json.Marshal(xaReq)
	if err != nil {
		t.Error(err)
	}
	println(string(data))
	if string(data) != xaReqJson {
		t.Error("request marshal json err")
	}

}

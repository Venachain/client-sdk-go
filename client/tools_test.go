package client

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetContractCallParamsSigned(t *testing.T) {
	address := "0x0000000000000000000000000000000000000099"
	vmType := "wasm"
	funcname := "saveEvidence"
	funcparam := []string{}
	funcparam = append(funcparam, "key2")
	funcparam = append(funcparam, "value")
	contractContent := NewExecuteContract(address, vmType, funcname, funcparam)

	abiPath := "/Users/cxh/go/src/VenaChain/sdk/client-sdk-go/precompiled/syscontracts/evidenceManager.cpp.abi.json"
	keyfile := "/Users/cxh/go/src/VenaChain/venachain/release/linux/conf/keyfile.json"
	key, err := NewKey(keyfile, "0")
	if err != nil {
		fmt.Println(err)
	}

	txSigned, err := GetTxSigned(*contractContent, abiPath, key)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(txSigned)
	assert.True(t, txSigned != "")
}

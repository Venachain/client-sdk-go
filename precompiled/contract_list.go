package precompile

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	"git-c.i.wxblockchain.com/vena/src/client-sdk-go/log"
	"git-c.i.wxblockchain.com/vena/src/client-sdk-go/precompiled/syscontracts"
)

var (
	UserManagementAddress        = syscontracts.UserManagementAddress.String()        // The Venachain Precompiled contract addr for user management
	NodeManagementAddress        = syscontracts.NodeManagementAddress.String()        // The Venachain Precompiled contract addr for node management
	CnsManagementAddress         = syscontracts.CnsManagementAddress.String()         // The Venachain Precompiled contract addr for CNS
	ParameterManagementAddress   = syscontracts.ParameterManagementAddress.String()   // The Venachain Precompiled contract addr for parameter management
	FirewallManagementAddress    = syscontracts.FirewallManagementAddress.String()    // The Venachain Precompiled contract addr for fire wall management
	GroupManagementAddress       = syscontracts.GroupManagementAddress.String()       // The Venachain Precompiled contract addr for group management
	ContractDataProcessorAddress = syscontracts.ContractDataProcessorAddress.String() // The Venachain Precompiled contract addr for group management
	CnsInvokeAddress             = syscontracts.CnsInvokeAddress.String()             // The Venachain Precompiled contract addr for group management
	NFTContractAddress           = syscontracts.NFTContractAddress.String()
	EvidenceManagementAddress    = syscontracts.EvidenceManagementAddress.String() // The Venachain Precompiled contract addr for evidence management
	BulletProofAddress           = syscontracts.BulletProofAddress.String()        // The Venachain Precompiled contract addr for Bullet proof
)

const (
	PermDeniedEvent = "the contract deployment is denied"
	CnsInvokeEvent  = "the event generated by cns Invoke"
	CnsInitRegEvent = "register the contract to cns from init()"
)

// link the precompiled contract addresses with abi file bytes
var List = map[string]string{
	UserManagementAddress:        "syscontracts/userManager.cpp.abi.json",
	NodeManagementAddress:        "syscontracts/nodeManager.cpp.abi.json",
	CnsManagementAddress:         "syscontracts/cnsManager.cpp.abi.json",
	ParameterManagementAddress:   "syscontracts/paramManager.cpp.abi.json",
	FirewallManagementAddress:    "syscontracts/fireWall.abi.json",
	GroupManagementAddress:       "syscontracts/groupManager.cpp.abi.json",
	ContractDataProcessorAddress: "syscontracts/contractData.cpp.abi.json",
	NFTContractAddress:           "syscontracts/nft.abi.json",
	EvidenceManagementAddress:    "syscontracts/evidenceManager.cpp.abi.json",
	BulletProofAddress:           "syscontracts/RangeProof.cpp.abi.json",

	CnsInitRegEvent: "syscontracts/cnsInitRegEvent.json",
	CnsInvokeEvent:  "syscontracts/cnsInvokeEvent.json",
	PermDeniedEvent: "syscontracts/permissionDeniedEvent.json",
}

func isWindowsSystem() bool {
	if runtime.GOOS == "windows" {
		return true
	}
	return false
}

func getCurrentFilePath() string {
	var separator string
	if isWindowsSystem() {
		separator = "\\"
	} else {
		separator = "/"
	}
	_, fileStr, _, _ := runtime.Caller(1)
	split_value := strings.Split(fileStr, "/")
	split_value = split_value[:len(split_value)-1]
	var result string
	for _, val := range split_value {
		result = result + val + separator
	}
	return result
}

func GetContractByte(jsonName string) ([]byte, error) {
	parentFilePath := getCurrentFilePath()
	objectPath := parentFilePath + jsonName
	file, err := os.Open(objectPath)
	if err != nil {
		log.Error("error", err)
		return nil, err
	}
	defer file.Close()
	buf := bytes.NewBuffer(nil)
	if _, err = io.Copy(buf, file); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return buf.Bytes(), nil
}

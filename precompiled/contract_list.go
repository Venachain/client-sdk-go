package precompile

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	"git-c.i.wxblockchain.com/PlatONE/src/node/client-sdk-go/precompiled/syscontracts"
)

var (
	UserManagementAddress        = syscontracts.UserManagementAddress.String()        // The PlatONE Precompiled contract addr for user management
	NodeManagementAddress        = syscontracts.NodeManagementAddress.String()        // The PlatONE Precompiled contract addr for node management
	CnsManagementAddress         = syscontracts.CnsManagementAddress.String()         // The PlatONE Precompiled contract addr for CNS
	ParameterManagementAddress   = syscontracts.ParameterManagementAddress.String()   // The PlatONE Precompiled contract addr for parameter management
	FirewallManagementAddress    = syscontracts.FirewallManagementAddress.String()    // The PlatONE Precompiled contract addr for fire wall management
	GroupManagementAddress       = syscontracts.GroupManagementAddress.String()       // The PlatONE Precompiled contract addr for group management
	ContractDataProcessorAddress = syscontracts.ContractDataProcessorAddress.String() // The PlatONE Precompiled contract addr for group management
	CnsInvokeAddress             = syscontracts.CnsInvokeAddress.String()             // The PlatONE Precompiled contract addr for group management
	NFTContractAddress           = syscontracts.NFTContractAddress.String()
	EvidenceManagementAddress    = syscontracts.EvidenceManagementAddress.String() // The PlatONE Precompiled contract addr for evidence management
	BulletProofAddress           = syscontracts.BulletProofAddress.String()        // The PlatONE Precompiled contract addr for Bullet proof
)

const (
	PermDeniedEvent = "the contract deployment is denied"
	CnsInvokeEvent  = "the event generated by cns Invoke"
	CnsInitRegEvent = "register the contract to cns from init()"
)

// link the precompiled contract addresses with abi file bytes
var List = map[string]string{
	UserManagementAddress:        "../precompiled/syscontracts/userManager.cpp.abi.json",
	NodeManagementAddress:        "../precompiled/syscontracts/nodeManager.cpp.abi.json",
	CnsManagementAddress:         "../precompiled/syscontracts/cnsManager.cpp.abi.json",
	ParameterManagementAddress:   "../precompiled/syscontracts/paramManager.cpp.abi.json",
	FirewallManagementAddress:    "../precompiled/syscontracts/fireWall.abi.json",
	GroupManagementAddress:       "../precompiled/syscontracts/groupManager.cpp.abi.json",
	ContractDataProcessorAddress: "../precompiled/syscontracts/contractData.cpp.abi.json",
	NFTContractAddress:           "../precompiled/syscontracts/nft.abi.json",
	EvidenceManagementAddress:    "evidenceManager.cpp.abi.json",
	BulletProofAddress:           "RangeProof.cpp.abi.json",

	CnsInitRegEvent: "../precompiled/syscontracts/cnsInitRegEvent.json",
	CnsInvokeEvent:  "../precompiled/syscontracts/cnsInvokeEvent.json",
	PermDeniedEvent: "../precompiled/syscontracts/permissionDeniedEvent.json",
}

func isWindowsSystem() bool {
	if runtime.GOOS == "windows" {
		return true
	}
	return false
}

func getCurrentFilePath() string {
	var separator = "/"
	if isWindowsSystem() {
		separator = "\\"
	}
	_, fileStr, _, _ := runtime.Caller(1)
	split_value := strings.Split(fileStr, separator)
	split_value = split_value[:len(split_value)-1]
	var result string
	for _, val := range split_value {
		result = result + val + separator
	}
	return result
}

func GetContractByte(jsonName string) ([]byte, error) {
	parentFilePath := getCurrentFilePath()
	objectPath := parentFilePath + "syscontracts/" + jsonName
	file, err := os.Open(objectPath)
	if err != nil {
		fmt.Println(err)
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

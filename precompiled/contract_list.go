package precompile

import "git-c.i.wxblockchain.com/PlatONE/src/node/client-sdk-go/precompiled/syscontracts"

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

	CnsInitRegEvent: "../precompiled/syscontracts/cnsInitRegEvent.json",
	CnsInvokeEvent:  "../precompiled/syscontracts/cnsInvokeEvent.json",
	PermDeniedEvent: "../precompiled/syscontracts/permissionDeniedEvent.json",
}

package client

import (
	"context"
	"encoding/json"
	"git-c.i.wxblockchain.com/PlatONE/src/node/client-sdk-go/types"

	"git-c.i.wxblockchain.com/PlatONE/src/node/client-sdk-go/common"
	"git-c.i.wxblockchain.com/PlatONE/src/node/client-sdk-go/packet"
	common_platone "git-c.i.wxblockchain.com/PlatONE/src/node/client-sdk-go/platone/common"
	"git-c.i.wxblockchain.com/PlatONE/src/node/client-sdk-go/precompiled/syscontracts"
)

type IClient interface {
	RpcCall(ctx context.Context, funcName string, funcParam interface{}) (json.RawMessage, error)
	GetBlockByHash(hash string) (*types.GetBlockResponse, error)
}

type IContract interface {
	Deploy(ctx context.Context, txparam common.TxParams, codepath string, consParams []string) (interface{}, error)
	ListContractMethods() (packet.ContractContent, error)
	Execute(ctx context.Context, txparam common.TxParams, funcName string, funcParams []string, address string) (interface{}, error)
	IsFuncNameInContract(funcName string) (bool, error)
	GetReceipt(txhash string) (*packet.Receipt, error)
}

type IAccount interface {
	UserAdd(ctx context.Context, name, phone, email, organization string) (string, error)
	UserUpdate(ctx context.Context, phone, email, organization string) (string, error)
	QueryUser(ctx context.Context, user string) (string, error)
	Lock(ctx context.Context) (bool, error)
	UnLock(ctx context.Context, passphrase string) (bool, error)
	CreateAccount(ctx context.Context, passphrase string) (*common_platone.Address, error)
}

type ICns interface {
	CnsExecute(ctx context.Context, funcName string, funcParams []string, cns string) (interface{}, error)
	CnsResolve(ctx context.Context, version string) (string, error)
	CnsRedirect(ctx context.Context, version string) (string, error)
	CnsQueryAll(ctx context.Context) (string, error)
	CnsQueryByName(ctx context.Context) (string, error)
	CnsQueryByAddress(ctx context.Context, address string) (string, error)
	CnsQueryByAccount(ctx context.Context, account string) (string, error)
	CnsStateByAddress(ctx context.Context, address string) (int32, error)
	CnsState(ctx context.Context) (int32, error)
}

type IFireWall interface {
	FwStatus(ctx context.Context) (string, error)
	FwStart(ctx context.Context) (string, error)
	FwClose(ctx context.Context) (string, error)
	FwExport(ctx context.Context, filePath string) (bool, error)
	FwImport(ctx context.Context, filePath string) (string, error)
	FwNew(ctx context.Context, action, targetAddr, api string) (string, error)
	FwDelete(ctx context.Context, action, targetAddr, api string) (string, error)
	FwReset(ctx context.Context, action, targetAddr, api string) (string, error)
	FwClear(ctx context.Context, action string) (string, error)
}

type INode interface {
	NodeAdd(ctx context.Context, requestNodeInfo syscontracts.NodeInfo) (string, error)
	NodeDelete(ctx context.Context) (string, error)
	NodeUpdate(ctx context.Context, request syscontracts.NodeUpdateInfo) (string, error)
	NodeQuery(ctx context.Context, request *syscontracts.NodeQueryInfo) (string, error)
	NodeStat(ctx context.Context, request *syscontracts.NodeStatInfo) (int32, error)
}

type IRole interface {
	SetSuperAdmin(ctx context.Context) (string, error)
	TransferSuperAdmin(ctx context.Context, address string) (string, error)
	AddChainAdmin(ctx context.Context, address string) (string, error)
	DelChainAdmin(ctx context.Context, address string) (string, error)
	AddGroupAdmin(ctx context.Context, address string) (string, error)
	DelGroupAdmin(ctx context.Context, address string) (string, error)
	AddNodeAdmin(ctx context.Context, address string) (string, error)
	DelNodeAdmin(ctx context.Context, address string) (string, error)
	AddContractAdmin(ctx context.Context, address string) (string, error)
	DelContractAdmin(ctx context.Context, address string) (string, error)
	AddContractDeployer(ctx context.Context, address string) (string, error)
	DelContractDeployer(ctx context.Context, address string) (string, error)
	GetAddrListOfRole(ctx context.Context, role string) (string, error)
	GetRoles(ctx context.Context, address string) (string, error)
	HasRole(ctx context.Context, address, role string) (int32, error)
}

type ISysconfig interface {
	SetSysConfig(ctx context.Context, request SysConfigParam) ([]string, error)
	GetTxGasLimit(ctx context.Context) (uint64, error)
	GetBlockGasLimit(ctx context.Context) (uint64, error)
	GetGasContractName(ctx context.Context) (string, error)
	GetIsProduceEmptyBlock(ctx context.Context) (uint32, error)
	GetCheckContractDeployPermission(ctx context.Context) (uint32, error)
	GetAllowAnyAccountDeployContract(ctx context.Context) (uint32, error)
	GetIsApproveDeployedContract(ctx context.Context) (uint32, error)
	GetIsTxUseGas(ctx context.Context) (uint32, error)
	GetVRFParams(ctx context.Context) (string, error)
}

package client

import (
	"context"
	"encoding/json"

	"github.com/PlatONE_Network/PlatONE-SDK-Go/packet"

	"github.com/PlatONE_Network/PlatONE-SDK-Go/precompiled/syscontracts"

	"github.com/PlatONE_Network/PlatONE-SDK-Go/common"
	common_sdk "github.com/PlatONE_Network/PlatONE-SDK-Go/common"
	"github.com/PlatONE_Network/PlatONE-SDK-Go/platone/rpc"
)

type IClient interface {
	GetRpcClient() *rpc.Client
	RPCSend(ctx context.Context, method string, args ...interface{}) (json.RawMessage, error)
	ClientSend(action string, params []interface{}) (interface{}, error)
}

type IContract interface {
	Deploy(ctx context.Context, txparam common.TxParams, codepath string, consParams []string) (interface{}, error)
	ListContractMethods() (packet.ContractAbi, error)
	Execute(ctx context.Context, txparam common.TxParams, funcName string, funcParams []string, address string) (interface{}, error)
	IsFuncNameInContract(funcName string) (bool, error)
	GetReceipt(txhash string) (*packet.Receipt, error)
}

type IAccount interface {
	UserAdd(ctx context.Context, txparam common_sdk.TxParams, name, phone, email, organization string) (string, error)
	UserUpdate(ctx context.Context, txparam common_sdk.TxParams, phone, email, organization string) (string, error)
	QueryUser(ctx context.Context, txparam common_sdk.TxParams, user string) (string, error)
	Lock(ctx context.Context) (bool, error)
	UnLock(ctx context.Context) (bool, error)
}

type ICns interface {
	CnsExecute(ctx context.Context, txparam common.TxParams, funcName string, funcParams []string, cns string) (interface{}, error)
	CnsRegister(ctx context.Context, txparam common_sdk.TxParams, version, address string) (string, error)
	CnsResolve(ctx context.Context, txparam common_sdk.TxParams, version string) (string, error)
	CnsRedirect(ctx context.Context, txparam common_sdk.TxParams, version string) (string, error)
	CnsQueryAll(ctx context.Context, txparam common_sdk.TxParams) (string, error)
	CnsQueryByName(ctx context.Context, txparam common_sdk.TxParams) (string, error)
	CnsQueryByAddress(ctx context.Context, txparam common_sdk.TxParams, address string) (string, error)
	CnsQueryByAccount(ctx context.Context, txparam common_sdk.TxParams, account string) (string, error)
	CnsStateByAddress(ctx context.Context, txparam common_sdk.TxParams, address string) (int32, error)
	CnsState(ctx context.Context, txparam common_sdk.TxParams) (int32, error)
}

type IFireWall interface {
	FwStatus(ctx context.Context, txparam common_sdk.TxParams) (string, error)
	FwStart(ctx context.Context, txparam common_sdk.TxParams) (string, error)
	FwClose(ctx context.Context, txparam common_sdk.TxParams) (string, error)
	FwExport(ctx context.Context, txparam common_sdk.TxParams, filePath string) (bool, error)
	FwImport(ctx context.Context, txparam common_sdk.TxParams, filePath string) (string, error)
	FwNew(ctx context.Context, txparam common_sdk.TxParams, action, targetAddr, api string) (string, error)
	FwDelete(ctx context.Context, txparam common_sdk.TxParams, action, targetAddr, api string) (string, error)
	FwReset(ctx context.Context, txparam common_sdk.TxParams, action, targetAddr, api string) (string, error)
	FwClear(ctx context.Context, txparam common_sdk.TxParams, action string) (string, error)
}

type INode interface {
	NodeAdd(ctx context.Context, txparam common_sdk.TxParams, requestNodeInfo syscontracts.NodeInfo) (string, error)
	NodeDelete(ctx context.Context, txparam common_sdk.TxParams) (string, error)
	NodeUpdate(ctx context.Context, txparam common_sdk.TxParams, request syscontracts.NodeUpdateInfo) (string, error)
	NodeQuery(ctx context.Context, txparam common_sdk.TxParams, request *syscontracts.NodeQueryInfo) (string, error)
	NodeStat(ctx context.Context, txparam common_sdk.TxParams, request *syscontracts.NodeStatInfo) (int32, error)
}

type IRole interface {
	SetSuperAdmin(ctx context.Context, txparam common_sdk.TxParams) (string, error)
	TransferSuperAdmin(ctx context.Context, txparam common_sdk.TxParams, address string) (string, error)
	AddChainAdmin(ctx context.Context, txparam common_sdk.TxParams, address string) (string, error)
	DelChainAdmin(ctx context.Context, txparam common_sdk.TxParams, address string) (string, error)
	AddGroupAdmin(ctx context.Context, txparam common_sdk.TxParams, address string) (string, error)
	DelGroupAdmin(ctx context.Context, txparam common_sdk.TxParams, address string) (string, error)
	AddNodeAdmin(ctx context.Context, txparam common_sdk.TxParams, address string) (string, error)
	DelNodeAdmin(ctx context.Context, txparam common_sdk.TxParams, address string) (string, error)
	AddContractAdmin(ctx context.Context, txparam common_sdk.TxParams, address string) (string, error)
	DelContractAdmin(ctx context.Context, txparam common_sdk.TxParams, address string) (string, error)
	AddContractDeployer(ctx context.Context, txparam common_sdk.TxParams, address string) (string, error)
	DelContractDeployer(ctx context.Context, txparam common_sdk.TxParams, address string) (string, error)
	GetAddrListOfRole(ctx context.Context, txparam common_sdk.TxParams, role string) (string, error)
	GetRoles(ctx context.Context, txparam common_sdk.TxParams, address string) (string, error)
	HasRole(ctx context.Context, txparam common_sdk.TxParams, address, role string) (int32, error)
}

type ISysconfig interface {
	SetSysConfig(ctx context.Context, txparam common_sdk.TxParams, request SysConfigParam) ([]string, error)
	GetTxGasLimit(ctx context.Context, txparam common_sdk.TxParams) (uint64, error)
	GetBlockGasLimit(ctx context.Context, txparam common_sdk.TxParams) (uint64, error)
	GetGasContractName(ctx context.Context, txparam common_sdk.TxParams) (string, error)
	GetIsProduceEmptyBlock(ctx context.Context, txparam common_sdk.TxParams) (uint32, error)
	GetCheckContractDeployPermission(ctx context.Context, txparam common_sdk.TxParams) (uint32, error)
	GetAllowAnyAccountDeployContract(ctx context.Context, txparam common_sdk.TxParams) (uint32, error)
	GetIsApproveDeployedContract(ctx context.Context, txparam common_sdk.TxParams) (uint32, error)
	GetIsTxUseGas(ctx context.Context, txparam common_sdk.TxParams) (uint32, error)
	GetVRFParams(ctx context.Context, txparam common_sdk.TxParams) (string, error)
}

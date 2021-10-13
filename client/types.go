package client

import (
	"github.com/PlatONE_Network/PlatONE-SDK-Go/common"
	"github.com/PlatONE_Network/PlatONE-SDK-Go/packet"
)

type IContract interface {
	Deploy(txparam common.TxParams, consParams []string) ([]interface{}, error)
	ListContractMethods() (packet.ContractAbi, error)
	Execute(txparam common.TxParams, funcName string, funcParams []string) ([]interface{}, error)
	IsFuncNameInContract(funcName string) (bool, error)
	GetReceipt(txhash string) (*packet.Receipt, error)
	contractCallWrap(txparam common.TxParams, funcParams []string, funcName, contract string) ([]interface{}, error)
}

type IAccount interface {
}

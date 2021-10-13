package client

import common_sdk "github.com/PlatONE_Network/PlatONE-SDK-Go/common"

type SysConfigClient struct {
	ContractClient
}

func (sysConfigClient SysConfigClient) SetSysConfig(txparam common_sdk.TxParams) (string, error) {

}

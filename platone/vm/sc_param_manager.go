package vm

import "errors"

var (
	GasContractNameKey                 string = "GasContractName"
	IsProduceEmptyBlockKey             string = "IsProduceEmptyBlock"
	TxGasLimitKey                      string = "TxGasLimit"
	BlockGasLimitKey                   string = "BlockGasLimit"
	IsCheckContractDeployPermissionKey string = "IsCheckContractDeployPermission"
	IsApproveDeployedContractKey       string = "IsApproveDeployedContract"
	IsTxUseGasKey                      string = "IsTxUseGas"
	VrfParamsKey                       string = "VRFParams"
	IsBlockUseTrieHashKey              string = "IsBlockUseTrieHash"
)

var (
	errDataTypeInvalid = errors.New("the data type invalid")
	errUnsupported     = errors.New("the operation is unsupported")
)

const (
	TxGasLimitMinValue        uint64 = 12771596 * 100 // 12771596 大致相当于 0.012772s
	TxGasLimitMaxValue        uint64 = 2e9            // 相当于 2s
	txGasLimitDefaultValue    uint64 = 1.5e9          // 相当于 1.5s
	BlockGasLimitMinValue     uint64 = 12771596 * 100 // 12771596 大致相当于 0.012772s
	BlockGasLimitMaxValue     uint64 = 2e10           // 相当于 20s
	blockGasLimitDefaultValue uint64 = 1e10           // 相当于 10s
	failFlag                         = -1
	sucFlag                          = 0
)

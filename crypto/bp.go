package crypto

import (
	"math/big"

	"git-c.i.wxblockchain.com/vena/src/client-sdk-go/venachain/crypto"
	"git-c.i.wxblockchain.com/vena/src/client-sdk-go/venachain/crypto/bp"
)

// 输入两个值和范围，生成proof
func GetProof(value []*big.Int, scope string) (string, error) {
	range_hash := crypto.Keccak256([]byte(scope))
	param := bp.GenerateAggBpStatement_range(2, 16, range_hash)
	proof, err := bp.AggBpProve_s(param, value)
	if err != nil {
		return "", err
	}
	return proof, nil
}

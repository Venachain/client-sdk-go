package crypto

import (
	"math/big"

	"github.com/PlatONE_Network/PlatONE-SDK-Go/platone/crypto"
	"github.com/PlatONE_Network/PlatONE-SDK-Go/platone/crypto/bp"
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

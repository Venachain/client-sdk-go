package crypto

import (
	"fmt"
	"math/big"
	"testing"
)

func TestGetProof(t *testing.T) {
	v := make([]*big.Int, 2)
	v[0] = big.NewInt(3)
	v[1] = big.NewInt(66)
	scope := "test"
	proof, err := GetProof(v, scope)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(proof)
}

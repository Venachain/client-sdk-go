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
func BenchmarkGetProofParallel(b *testing.B) {
	v := make([]*big.Int, 2)
	v[0] = big.NewInt(3)
	v[1] = big.NewInt(66)
	scope := "1.000000000000000000000000000000_100.000000000000000000000000000000"
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := GetProof(v, scope)
			if err != nil {
				fmt.Println(err)
			}
		}
	})
}

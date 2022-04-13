package common

import (
	"encoding/binary"
	"math"
	"math/big"

	math2 "git-c.i.wxblockchain.com/vena/src/client-sdk-go/venachain/common/math"
)

func CallResAsUint128(bts []byte) *big.Int {
	if len(bts) < 32 {
		return &big.Int{}
	}
	return Byte128ToBig(bts[len(bts)-16:], false)
}

func CallResAsUint64(bts []byte) uint64 {
	if len(bts) < 32 {
		return 0
	}

	var n uint64 = 0
	for _, b := range bts[:32] {
		n = n*256 + uint64(b)
	}
	return n
}

func CallResAsUint32(bts []byte) uint32 {
	if len(bts) < 32 {
		return 0
	}

	var n uint32 = 0
	for _, b := range bts[:32] {
		n = n*256 + uint32(b)
	}
	return n
}

func CallResAsFloat128(bts []byte) *big.Float {
	if len(bts) < 32 {
		return &big.Float{}
	}

	bytes := bts[len(bts)-16:]
	low := binary.BigEndian.Uint64(bytes[8:])
	high := binary.BigEndian.Uint64(bytes[:8])

	F, _ := math2.NewFromBits(high, low).Big()

	return F
}

func CallResAsFloat64(bts []byte) float64 {
	if len(bts) < 32 {
		return 0
	}
	bits := binary.BigEndian.Uint64(bts[len(bts)-8:])
	return math.Float64frombits(bits)
}

func CallResAsFloat32(bts []byte) float32 {
	if len(bts) < 32 {
		return 0
	}
	bits := binary.BigEndian.Uint32(bts[len(bts)-4:])
	return math.Float32frombits(bits)
}

func CallResAsInt128(bts []byte) *big.Int {

	if len(bts) < 32 {
		return new(big.Int).SetInt64(0)
	}
	return Byte128ToBig(bts[len(bts)-16:], true)
}

func CallResAsInt64(bts []byte) int64 {
	if len(bts) < 32 {
		return 0
	}

	var n int64 = 0
	for _, b := range bts[:32] {
		n = n*256 + int64(b)
	}
	return n
}

func CallResAsInt32(bts []byte) int32 {
	if len(bts) < 32 {
		return 0
	}

	var n int32 = 0
	for _, b := range bts[:32] {
		n = n*256 + int32(b)
	}
	return n
}

func CallResAsBool(bts []byte) bool {
	if len(bts) < 32 {
		return false
	}

	if bts[31] == 1 {
		return true
	} else {
		return false
	}
}

func CallResAsString(bts []byte) string {
	if len(bts) < 64 {
		return ""
	}

	slen := int(bts[61])*256*256 + int(bts[62])*256 + int(bts[63])
	if slen > len(bts)-64 {
		return ""
	}
	return string(bts[64 : 64+slen])
}

func RevertBytes(bts []byte) {
	for i, j := 0, len(bts)-1; i < j; {
		temp := bts[i]
		bts[i] = bts[j]
		bts[j] = temp
		i++
		j--
	}
}

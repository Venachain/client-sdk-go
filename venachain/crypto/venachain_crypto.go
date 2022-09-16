package crypto

import (
	"crypto/elliptic"
	"hash"
	"math/big"

	"github.com/Venachain/client-sdk-go/venachain/common"
)

// Hasher 通用哈希接口
type Hasher interface {
	// Hash256
	// 返回 256 位哈希码的哈希字节数组
	Hash256(data ...[]byte) []byte
	// Hash256Hash
	// 返回 256 位哈希码的哈希对象
	Hash256Hash(data ...[]byte) common.Hash

	// ... 其他安全长度的哈希方法
}

// Key 通用密钥接口
type Key interface {
	// Bytes 转为 []byte
	Bytes() ([]byte, error)
	// Equals 比较两个 key 是否相同
	Equals(key Key) bool
	// SaveToFile 保存私钥到给定文件
	SaveToFile(file string) error
	// LoadByFile 从给定文件中加载私钥
	LoadByFile(file string) error
	// ENRKey 获取密钥曲线名称
	ENRKey() string
}

// PrivKey 通用私钥接口
type PrivKey interface {
	Key

	// Sign 私钥签名
	// data 待签名的数据
	Sign(data []byte) ([]byte, error)

	// GetPubKey 获取该私钥对应的公钥
	GetPubKey() PubKey

	// FromBytes 从私钥字节数据中还原出私钥对象
	// d 私钥字节数据
	FromBytes(d []byte) error

	// FromHex 从私钥哈希还原出私钥对象
	// hexkey 私钥哈希值
	FromHex(hexkey string) error

	// Zero 将当前 PrivKey 置为零值
	Zero()

	// GetD 返回私钥数据
	GetD() *big.Int
}

// PubKey 通用公钥接口
type PubKey interface {
	Key

	// Verify 公钥验证签名
	// hash 原数据的哈希值
	// sig 用私钥对原始数据签名后的签名数据
	Verify(hash []byte, sig []byte) (bool, error)

	// GetAddress 获取公钥对应的地址
	GetAddress() common.Address

	// FromBytes 从哈希中还原出公钥对象
	// d 公钥字节数据
	FromBytes(d []byte) error

	// GetCurve 获取曲线数据
	GetCurve() elliptic.Curve

	// GetX 获取 x
	GetX() *big.Int

	// GetY 获取 y
	GetY() *big.Int

	// IsOnCurve 判断点（x, y）是否在椭圆曲线上
	// x, y 点坐标
	IsOnCurve(x, y *big.Int) bool
}

// Signature 通用签名接口
type Signature interface {
	// GetSig 获取签名数据
	GetSig() []byte

	// GetPubKey 获取公钥数据
	GetPubKey() []byte

	// RecoverPubkeyBytes 签名恢复为[]byte公钥
	// hash 签名数据的哈希值
	RecoverPubkeyBytes(hash []byte) ([]byte, error)

	// RecoverPubKey 签名恢复为PubKey公钥
	// hash 签名数据的哈希值
	RecoverPubKey(hash []byte) (PubKey, error)
}

// HashState trie/hasher.go keccakState
type HashState interface {
	hash.Hash
	// Read 从hash state中读取数据，作用与Sum相同（Sum的返回结果需要取后半部分）
	// 不同的是 Read 改变hash state，Sum将改变的结果拼接在hash buf后并返回
	Read([]byte) (int, error)
}

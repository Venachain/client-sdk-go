package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"io"
	"io/ioutil"
	"math/big"
	"os"

	"git-c.i.wxblockchain.com/vena/src/client-sdk-go/venachain/common"
	"git-c.i.wxblockchain.com/vena/src/client-sdk-go/venachain/crypto/secp256k1"
	"git-c.i.wxblockchain.com/vena/src/client-sdk-go/venachain/rlp"
)

// SHA3Hasher SHA3哈希实现
type SHA3Hasher struct {
}

// Hash256 返回 256 位哈希码的哈希字节数组
func (sha3 *SHA3Hasher) Hash256(data ...[]byte) []byte {
	return Keccak256(data...)
}

// Hash256Hash 返回 256 位哈希码的哈希对象
func (sha3 *SHA3Hasher) Hash256Hash(data ...[]byte) common.Hash {
	return Keccak256Hash(data...)
}

// ECDSAPrivKey ECDSA私钥实现
type ECDSAPrivKey struct {
	*ecdsa.PrivateKey
}

func NewECDSAPrivKey(pub *ECDSAPubKey, D *big.Int) *ECDSAPrivKey {
	return &ECDSAPrivKey{
		&ecdsa.PrivateKey{
			PublicKey: *pub.pub,
			D:         D,
		},
	}
}

// Bytes 转为 []byte
func (privKey *ECDSAPrivKey) Bytes() ([]byte, error) {
	return FromECDSA(privKey.PrivateKey), nil
}

// Equals 比较两个 key 是否相同
func (privKey *ECDSAPrivKey) Equals(key Key) bool {
	x, ok := key.(*ECDSAPrivKey)
	if !ok {
		return false
	}
	return privKey.PrivateKey.Equal(x.PrivateKey)
}

// Sign 私钥签名
// data 待签名的数据
func (privKey *ECDSAPrivKey) Sign(data []byte) ([]byte, error) {
	return Sign(data, privKey.PrivateKey)
}

// GetPubKey 获取该私钥对应的公钥
func (privKey *ECDSAPrivKey) GetPubKey() PubKey {
	return &ECDSAPubKey{&privKey.PublicKey}
}

// FromBytes 从私钥字节数据中还原出私钥对象
// d 私钥字节数据
func (privKey *ECDSAPrivKey) FromBytes(d []byte) error {
	privateKey, err := toECDSA(d, true)
	if err != nil {
		return err
	}
	privKey.PrivateKey = privateKey
	return nil
}

// FromHex 从私钥哈希还原出私钥对象
// hexkey 私钥哈希值
func (privKey *ECDSAPrivKey) FromHex(hexkey string) error {
	privateKey, err := HexToECDSA(hexkey)
	if err != nil {
		return err
	}
	privKey.PrivateKey = privateKey
	return nil
}

// LoadByFile 从给定文件中加载私钥
func (privKey *ECDSAPrivKey) LoadByFile(file string) error {
	privateKey, err := LoadECDSA(file)
	if err != nil {
		return err
	}
	privKey.PrivateKey = privateKey
	return nil
}

// SaveToFile 保存私钥到给定文件
func (privKey *ECDSAPrivKey) SaveToFile(file string) error {
	err := SaveECDSA(file, privKey.PrivateKey)
	if err != nil {
		return err
	}
	return nil
}

// Zero 将当前 PrivKey 置为零值
func (privKey *ECDSAPrivKey) Zero() {
	b := privKey.D.Bits()
	for i := range b {
		b[i] = 0
	}
}

// GetD 返回私钥数据
func (privKey *ECDSAPrivKey) GetD() *big.Int {
	return privKey.D
}

func (privKey *ECDSAPrivKey) ENRKey() string {
	return "secp256k1"
}

func (privKey *ECDSAPrivKey) EncodeRLP(w io.Writer) error {
	bytes, err := privKey.Bytes()
	if err != nil {
		return err
	}
	return rlp.Encode(w, bytes)
}

func (privKey *ECDSAPrivKey) DecodeRLP(s *rlp.Stream) error {
	bytes, err := s.Bytes()
	if err != nil {
		return err
	}
	return privKey.FromBytes(bytes)
}

// ECDSAPubKey ECDSA公钥实现
type ECDSAPubKey struct {
	pub *ecdsa.PublicKey
}

func NewECDSAPubKey(X *big.Int, Y *big.Int) *ECDSAPubKey {
	return &ECDSAPubKey{
		pub: &ecdsa.PublicKey{
			Curve: secp256k1.S256(),
			X:     X,
			Y:     Y,
		},
	}
}

// LoadByFile 从给定文件中加载私钥
func (pubKey *ECDSAPubKey) LoadByFile(file string) error {
	var buf []byte
	fd, err := os.Open(file)
	if err != nil {
		return err
	}
	defer fd.Close()
	if buf, err = io.ReadAll(fd); err != nil {
		return err
	}
	key, err := hex.DecodeString(string(buf))
	if err != nil {
		return err
	}
	return pubKey.FromBytes(key)
}

// SaveToFile 保存私钥到给定文件
func (pubKey *ECDSAPubKey) SaveToFile(file string) error {
	bytes, err := pubKey.Bytes()
	if err != nil {
		return err
	}
	key := hex.EncodeToString(bytes)
	return ioutil.WriteFile(file, []byte(key), 0600)
}

// FromBytes 从哈希中还原出公钥对象
// d 公钥字节数据
func (pubKey *ECDSAPubKey) FromBytes(pub []byte) error {
	publicKey, err := UnmarshalPubkey(pub)
	if err != nil {
		return err
	}
	pubKey.pub = publicKey
	return nil
}

// Bytes 转为 []byte
func (pubKey *ECDSAPubKey) Bytes() ([]byte, error) {
	return FromECDSAPub(pubKey.pub), nil
}

// Equals 比较两个 key 是否相同
func (pubKey *ECDSAPubKey) Equals(key Key) bool {
	x, ok := key.(*ECDSAPubKey)
	if !ok {
		return false
	}
	return pubKey.pub.Equal(x.pub)
}

// Verify 公钥验证签名
// hash 原数据的哈希值
// sig 用私钥对原始数据签名后的签名数据
func (pubKey *ECDSAPubKey) Verify(hash []byte, sig []byte) (bool, error) {
	pubkey, err := pubKey.Bytes()
	if err != nil {
		return false, err
	}
	// TODO：ECDSA 的签名数据是 65 字节，最后一位被塞进了一个 recid，不知道是什么
	// TODO：验签的时候需要去掉最后一位才能验签成功
	return VerifySignature(pubkey, hash, sig[:len(sig)-1]), nil
}

// GetAddress 获取公钥对应的地址
func (pubKey *ECDSAPubKey) GetAddress() common.Address {
	pubBytes := FromECDSAPub(pubKey.pub)
	return common.BytesToAddress(Keccak256(pubBytes[1:])[12:])
}

// GetCurve 获取曲线数据
func (pubKey *ECDSAPubKey) GetCurve() elliptic.Curve {
	return pubKey.pub.Curve
}

// GetX 获取 x
func (pubKey *ECDSAPubKey) GetX() *big.Int {
	return pubKey.pub.X
}

// GetY 获取 y
func (pubKey *ECDSAPubKey) GetY() *big.Int {
	return pubKey.pub.Y
}

func (pubKey *ECDSAPubKey) IsOnCurve(x, y *big.Int) bool {
	return pubKey.pub.Curve.IsOnCurve(x, y)
}

func (pubKey *ECDSAPubKey) ENRKey() string {
	return "secp256k1"
}

func (pubKey *ECDSAPubKey) EncodeRLP(w io.Writer) error {
	bytes, err := pubKey.Bytes()
	if err != nil {
		return err
	}
	return rlp.Encode(w, bytes)
}

func (pubKey *ECDSAPubKey) DecodeRLP(s *rlp.Stream) error {
	bytes, err := s.Bytes()
	if err != nil {
		return err
	}
	return pubKey.FromBytes(bytes)
}

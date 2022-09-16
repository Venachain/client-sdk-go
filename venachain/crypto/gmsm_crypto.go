package crypto

import (
	"crypto/elliptic"
	"crypto/subtle"
	"encoding/asn1"
	"encoding/hex"
	"errors"
	"hash"
	"io"
	"io/ioutil"
	"math/big"
	"os"

	"git-c.i.wxblockchain.com/vena/src/client-sdk-go/venachain/common"
	"git-c.i.wxblockchain.com/vena/src/client-sdk-go/venachain/common/math"
	"git-c.i.wxblockchain.com/vena/src/client-sdk-go/venachain/rlp"
	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/sm3"
)

const (
	SM3HashAlgo = 66
)

var (
	sm2P256N, _   = new(big.Int).SetString("FFFFFFFEFFFFFFFFFFFFFFFFFFFFFFFF7203DF6B21C6052B53BBF40939D54123", 16)
	sm2P256NhalfN = new(big.Int).Div(sm2P256N, big.NewInt(2))
)

// SM2Sig
// 持有 r、s 值的 SM2 签名结构
type SM2Sig struct {
	R, S *big.Int
}

// GMSMHasher
// GMSM 对 Hasher 接口的实现，采用的是 SM3 哈希算法
type GMSMHasher struct {
}

// Hash256 返回 SM3 哈希算法生成的 256 位哈希码的哈希字节数组
func (gmsmHasher *GMSMHasher) Hash256(data ...[]byte) []byte {
	d := sm3.New()
	for _, b := range data {
		d.Write(b)
	}
	return d.Sum(nil)
}

// Hash256Hash 返回 SM3 哈希算法生成的 256 位哈希码的哈希对象
func (gmsmHasher *GMSMHasher) Hash256Hash(data ...[]byte) (h common.Hash) {
	d := sm3.New()
	for _, b := range data {
		d.Write(b)
	}
	d.Sum(h[:0])
	return h
}

// GMSMSignature GMSM 对 Signature 接口的实现
type GMSMSignature struct {
	pubKey []byte
	sig    []byte
}

// GetSig
// 获取签名数据
func (signature *GMSMSignature) GetSig() []byte {
	return signature.sig
}

// GetPubKey
// 获取 Signature 的公钥数据
func (signature *GMSMSignature) GetPubKey() []byte {
	return signature.pubKey
}

// RecoverPubkeyBytes 返回字节数组形式的公钥，参数 hash 只是为了适配 RecoverPubKey()，
// 在国密实现方式中无需传入任何参数，传 nil 即可
// hash 原数据的哈希值
func (signature *GMSMSignature) RecoverPubkeyBytes(hash []byte) ([]byte, error) {
	return signature.pubKey, nil
}

// RecoverPubKey 返回公钥对象，参数 hash 只是为了适配 RecoverPubKey()，
// 在国密实现方式中无需传入任何参数，传 nil 即可
// hash 原数据的哈希值
func (signature *GMSMSignature) RecoverPubKey(hash []byte) (PubKey, error) {
	pubKey, err := BytesToPubKey(signature.pubKey)
	if err != nil {
		return nil, err
	}
	return pubKey, nil
}

// GMSMPrivKey SM2私钥实现
type GMSMPrivKey struct {
	*sm2.PrivateKey
}

func NewGMSMPrivKey(pub *GMSMPubKey, D *big.Int) *GMSMPrivKey {
	return &GMSMPrivKey{
		&sm2.PrivateKey{
			PublicKey: *pub.pub,
			D:         D,
		},
	}
}

// Bytes 返回私钥对应的 x509 字节数组
func (privKey *GMSMPrivKey) Bytes() ([]byte, error) {
	return math.PaddedBigBytes(privKey.D, privKey.Params().BitSize/8), nil
}

// Equals 判断两个私钥是否相等
func (privKey *GMSMPrivKey) Equals(key Key) bool {
	a, err := privKey.Bytes()
	if err != nil {
		return false
	}
	b, err := key.Bytes()
	if err != nil {
		return false
	}
	return subtle.ConstantTimeCompare(a, b) == 1
}

// Sign 用私钥对输入数据进行签名，返回签名数据
// data 待签名的数据
func (privKey *GMSMPrivKey) Sign(data []byte) ([]byte, error) {
	hash := sm3.Sm3Sum(data)
	r, s, err := sm2.Sign(privKey.PrivateKey, hash[:])
	if err != nil {
		return nil, err
	}
	return asn1.Marshal(SM2Sig{
		R: r,
		S: s,
	})
}

// GetPubKey 获取该私钥对应的公钥
func (privKey *GMSMPrivKey) GetPubKey() PubKey {
	return &GMSMPubKey{&privKey.PrivateKey.PublicKey}
}

// FromBytes 从 x509 字节数据中还原出私钥
// d 私钥字节数据
func (privKey *GMSMPrivKey) FromBytes(d []byte) error {
	priv, err := toGMPrivKey(d)
	if err != nil {
		return err
	}
	privKey.PrivateKey = priv
	return nil
}

// FromHex 从 16 进制地址中还原出私钥
// hexkey 私钥哈希值
func (privKey *GMSMPrivKey) FromHex(hexkey string) error {
	b, err := hex.DecodeString(hexkey)
	if err != nil {
		return errors.New("invalid hex string")
	}
	privateKey, err := toGMPrivKey(b)
	if err != nil {
		return err
	}
	privKey.PrivateKey = privateKey
	return nil
}

// LoadByFile 从给定的文件中加载出私钥
func (privKey *GMSMPrivKey) LoadByFile(file string) error {
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
	//if err != nil {
	//	return err
	//}
	return privKey.FromBytes(key)
}

// SaveToFile 保存私钥到给定文件
func (privKey *GMSMPrivKey) SaveToFile(file string) error {
	b, err := privKey.Bytes()
	if err != nil {
		return err
	}
	k := hex.EncodeToString(b)
	return ioutil.WriteFile(file, []byte(k), 0600)
}

// Zero 将当前 PrivKey 置为零值
func (privKey *GMSMPrivKey) Zero() {
	b := privKey.D.Bits()
	for i := range b {
		b[i] = 0
	}
}

// GetD 返回私钥数据
func (privKey *GMSMPrivKey) GetD() *big.Int {
	return privKey.D
}

func (privKey *GMSMPrivKey) ENRKey() string {
	return "sm2p256"
}
func (privKey *GMSMPrivKey) EncodeRLP(w io.Writer) error {
	bytes, err := privKey.Bytes()
	if err != nil {
		return err
	}
	return rlp.Encode(w, bytes)
}
func (privKey *GMSMPrivKey) DecodeRLP(s *rlp.Stream) error {
	bytes, err := s.Bytes()
	if err != nil {
		return err
	}
	return privKey.FromBytes(bytes)
}

// GMSMPubKey SM2公钥实现
type GMSMPubKey struct {
	pub *sm2.PublicKey
}

func NewGMSMPubKey(X *big.Int, Y *big.Int) *GMSMPubKey {
	return &GMSMPubKey{
		pub: &sm2.PublicKey{
			Curve: sm2.P256Sm2(),
			X:     X,
			Y:     Y,
		},
	}
}

// Bytes 将公钥转为字节数组
func (pubKey *GMSMPubKey) Bytes() ([]byte, error) {
	if pubKey == nil || pubKey.pub == nil || pubKey.pub.X == nil || pubKey.pub.Y == nil {
		return nil, nil
	}
	return elliptic.Marshal(P256Sm2(), pubKey.GetX(), pubKey.GetY()), nil
}

// Equals 判断两个公钥是否相等
func (pubKey *GMSMPubKey) Equals(key Key) bool {
	a, err := pubKey.Bytes()
	if err != nil {
		return false
	}
	b, err := key.Bytes()
	if err != nil {
		return false
	}
	return subtle.ConstantTimeCompare(a, b) == 1
}

// Verify 用公钥对给定数据和签名进行验签
// hash 原数据的哈希值
// sig 用私钥对原始数据签名后的签名数据
func (pubKey *GMSMPubKey) Verify(hash []byte, sig []byte) (bool, error) {
	s := new(SM2Sig)
	if _, err := asn1.Unmarshal(sig, s); err != nil {
		return false, err
	}
	if s == nil {
		return false, errors.New("sig is nil")
	}
	h := sm3.Sm3Sum(hash)
	return sm2.Verify(pubKey.pub, h[:], s.R, s.S), nil
}

// GetAddress 通过公钥中获取账户地址
func (pubKey *GMSMPubKey) GetAddress() common.Address {
	pubBytes, _ := pubKey.Bytes()
	return common.BytesToAddress(DefaultHasher.Hash256(pubBytes[1:])[12:])
}

// LoadByFile 从给定文件中加载公钥
func (pubKey *GMSMPubKey) LoadByFile(file string) error {
	var buf []byte
	fd, err := os.Open(file)
	if err != nil {
		return err
	}
	defer fd.Close()
	if buf, err = io.ReadAll(fd); err != nil {
		return err
	}

	key, _ := hex.DecodeString(string(buf))
	if err != nil {
		return err
	}
	return pubKey.FromBytes(key)
}

// SaveToFile 把公钥保存到给定文件中
func (pubKey *GMSMPubKey) SaveToFile(file string) error {
	b, err := pubKey.Bytes()
	if err != nil {
		return err
	}
	k := hex.EncodeToString(b)
	return ioutil.WriteFile(file, []byte(k), 0600)
}

// FromBytes 从给定的字节数据中还原出公钥
// d 公钥字节数据
func (pubKey *GMSMPubKey) FromBytes(d []byte) error {
	x, y := elliptic.Unmarshal(P256Sm2(), d)
	if x == nil {
		return errInvalidPubkeyGM
	}
	pubKey.pub = &sm2.PublicKey{Curve: P256Sm2(), X: x, Y: y}
	return nil
}

// GetCurve 获取曲线数据
func (pubKey *GMSMPubKey) GetCurve() elliptic.Curve {
	return pubKey.pub.Curve
}

// GetX 获取 x
func (pubKey *GMSMPubKey) GetX() *big.Int {
	return pubKey.pub.X
}

// GetY 获取 y
func (pubKey *GMSMPubKey) GetY() *big.Int {
	return pubKey.pub.Y
}

func (pubKey *GMSMPubKey) IsOnCurve(x, y *big.Int) bool {
	return pubKey.pub.IsOnCurve(x, y)
}

func (pubKey *GMSMPubKey) ENRKey() string {
	return "sm2p256"
}

func (pubKey *GMSMPubKey) EncodeRLP(w io.Writer) error {
	bytes, err := pubKey.Bytes()
	if err != nil {
		return err
	}
	return rlp.Encode(w, bytes)
}

func (pubKey *GMSMPubKey) DecodeRLP(s *rlp.Stream) error {
	bytes, err := s.Bytes()
	if err != nil {
		return err
	}
	return pubKey.FromBytes(bytes)
}

// ValidateSignatureValuesByGMSM SM2验证签名算法的曲线值
func ValidateSignatureValuesByGMSM(v byte, r, s *big.Int, homestead bool) bool {
	if r.Cmp(common.Big1) < 0 || s.Cmp(common.Big1) < 0 {
		return false
	}

	// 国密中不需要对 V 值进行处理
	// Frontier: allow s to be in full N range
	return r.Cmp(sm2P256N) < 0 && s.Cmp(sm2P256N) < 0
}

type SM3State struct {
	hash.Hash
}

func (hs *SM3State) Read(out []byte) (int, error) {
	return len(out), nil
}

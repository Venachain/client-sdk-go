package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"errors"
	"fmt"
	"hash"
	"math/big"
	"sync"

	"github.com/Venachain/client-sdk-go/venachain/common"
	"github.com/Venachain/client-sdk-go/venachain/crypto/secp256k1"
	"github.com/Venachain/client-sdk-go/venachain/crypto/sha3"
	"github.com/Venachain/client-sdk-go/venachain/rlp"
	"github.com/sirupsen/logrus"
	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/sm3"
)

const (
	GMPKLength = 65
)

var (
	GM    = "GM"
	ECDSA = "ECDSA"
)

var (
	once sync.Once
	// Encryption 当前链使用的密码学算法实现
	Encryption string
	// DefaultHasher 默认的哈希器实现
	DefaultHasher Hasher

	NoCryptoConfigErr = errors.New("unknown crypto type for server")
)

func init() {
	Encryption = ECDSA // defaults to ECDSA
	common.SetEncryption(Encryption)
	initHasher()
}

func IsGM() bool {
	return Encryption == GM
}

// 初始化 DefaultHasher
func initHasher() {
	switch Encryption {
	case GM:
		DefaultHasher = &GMSMHasher{}
	case ECDSA:
		DefaultHasher = &SHA3Hasher{}
	default:
		panic("unknown crypto type for server")
	}
}

// SetEncryption 指定链的密码学算法实现
func SetEncryption(encryption string) {
	// 只能在链启动的时候修改一次
	once.Do(func() {
		if encryption != ECDSA && encryption != GM {
			panic(fmt.Sprintf("unknown encryption type for server: %v", encryption))
		}
		Encryption = encryption
		common.SetEncryption(Encryption)
		initHasher() // re-init hash method since encryption could change
	})
}

// BytesToPubKey 从字节数组中恢复出公钥
// d 公钥字节数据
func BytesToPubKey(d []byte) (PubKey, error) {
	var pubkey PubKey
	switch Encryption {
	case GM:
		pubkey = &GMSMPubKey{}
	case ECDSA:
		pubkey = &ECDSAPubKey{}
	default:
		return nil, NoCryptoConfigErr
	}
	err := pubkey.FromBytes(d)
	if err != nil {
		return nil, err
	}
	return pubkey, nil
}

// BytesToPrivKey 从字节数组中恢复出私钥
// d 私钥字节数据
func BytesToPrivKey(d []byte) (PrivKey, error) {
	var privKey PrivKey
	switch Encryption {
	case GM:
		privKey = &GMSMPrivKey{}
	case ECDSA:
		privKey = &ECDSAPrivKey{}
	default:
		return nil, NoCryptoConfigErr
	}
	err := privKey.FromBytes(d)
	if err != nil {
		return nil, err
	}
	return privKey, nil
}

// BytesToPrivKeyUnsafe 从字节数组中恢复出私钥，不做严格检查（不安全）
// d 私钥字节数据
func BytesToPrivKeyUnsafe(d []byte) PrivKey {
	var privKey PrivKey
	switch Encryption {
	case GM:
		key, err := BytesToPrivKey(d)
		if err != nil {
			logrus.Warn("err")
		}
		privKey = key
	case ECDSA:
		key := ToECDSAUnsafe(d)
		privKey = &ECDSAPrivKey{key}
	default:
		return nil
	}
	return privKey
}

// HexToPrivKey 从私钥哈希还原出私钥对象
// hexkey 私钥哈希
func HexToPrivKey(hexkey string) (PrivKey, error) {
	var privKey PrivKey
	switch Encryption {
	case GM:
		privKey = &GMSMPrivKey{}
	case ECDSA:
		privKey = &ECDSAPrivKey{}
	default:
		return nil, NoCryptoConfigErr
	}
	err := privKey.FromHex(hexkey)
	if err != nil {
		return nil, err
	}
	return privKey, nil
}

// LoadPrivKeyFromFile 从给定文件中加载私钥
// file 目标文件路径
func LoadPrivKeyFromFile(file string) (PrivKey, error) {
	var key PrivKey
	switch Encryption {
	case GM:
		key = &GMSMPrivKey{}
	case ECDSA:
		key = &ECDSAPrivKey{}
	default:
		return nil, NoCryptoConfigErr
	}
	err := key.LoadByFile(file)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// LoadPubKeyFromFile 从给定文件中加载公钥
// file 目标文件路径
func LoadPubKeyFromFile(file string) (PubKey, error) {
	var key PubKey
	switch Encryption {
	case GM:
		key = &GMSMPubKey{}
	case ECDSA:
		key = &ECDSAPubKey{}
	default:
		return nil, NoCryptoConfigErr
	}
	err := key.LoadByFile(file)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// SavePrivKeyToFile 保存私钥到给定文件中
// file 目标文件路径
// privKey 私钥
func SavePrivKeyToFile(file string, privKey PrivKey) error {
	return saveKeyToFile(file, privKey)
}

// SavePubKeyToFile 保存公钥到给定文件中
// file 目标文件路径
// pubKey 公钥
func SavePubKeyToFile(file string, pubKey PubKey) error {
	return saveKeyToFile(file, pubKey)
}

// saveKeyToFile 保存密钥到指定文件中
// file 目标文件路径
// key 密钥
func saveKeyToFile(file string, key Key) error {
	if key == nil {
		return errors.New("key must not be nil")
	}
	return key.SaveToFile(file)
}

// PubKeyToAddress 获取公钥对应的地址
// pubKey 公钥
func PubKeyToAddress(pubKey PubKey) common.Address {
	return pubKey.GetAddress()
}

// GenerateKeypair 生成密钥对
func GenerateKeypair() (PrivKey, error) {
	switch Encryption {
	case GM:
		privateKey, err := sm2.GenerateKey()
		if err != nil {
			return nil, err
		}
		return &GMSMPrivKey{privateKey}, nil
	case ECDSA:
		privateKey, err := ecdsa.GenerateKey(S256(), rand.Reader)
		if err != nil {
			return nil, err
		}
		return &ECDSAPrivKey{privateKey}, nil
	default:
		return nil, NoCryptoConfigErr
	}
}

// Sign2 私钥签名
// privKey 私钥
// data 待签名的数据
func Sign2(privKey PrivKey, data []byte) ([]byte, error) {
	if privKey == nil {
		return nil, errors.New("privKey must not be nil")
	}
	return privKey.Sign(data)
}

// Verify 验证签名
// pubkey 公钥
// hash 原数据的哈希值
// sig 签名数据
func Verify(pubkey, hash, sig []byte) (bool, error) {
	pubKey, err := BytesToPubKey(pubkey)
	if err != nil {
		return false, err
	}
	return pubKey.Verify(hash, sig)
}

// ValidateSignatureValues2 验证签名算法椭圆曲线的值是否合法
func ValidateSignatureValues2(v byte, r, s *big.Int, homestead bool) bool {
	switch Encryption {
	case GM:
		return ValidateSignatureValuesByGMSM(v, r, s, homestead)
	case ECDSA:
		return ValidateSignatureValues(v, r, s, homestead)
	default:
		return false
	}
}

// NewPubKey 创建公钥
func NewPubKey(X *big.Int, Y *big.Int) PubKey {
	switch Encryption {
	case GM:
		return NewGMSMPubKey(X, Y)
	case ECDSA:
		return NewECDSAPubKey(X, Y)
	default:
		return nil
	}
}

func RecoverPubKey(payload, hash, sig []byte) ([]byte, error) {
	switch Encryption {
	case GM:
		if len(payload) < GMPKLength {
			return nil, errors.New("invalid sm2 public key")
		}
		pubkey := make([]byte, GMPKLength)
		copy(pubkey, payload[:GMPKLength])
		return pubkey, nil
	case ECDSA:
		return secp256k1.RecoverPubkey(hash, sig)
	default:
		return nil, NoCryptoConfigErr
	}
}

func NewHash() hash.Hash {
	switch Encryption {
	case GM:
		return sm3.New()
	case ECDSA:
		return sha3.NewKeccak256()
	default:
		return nil
	}
}

func NewHashState() HashState {
	switch Encryption {
	case GM:
		return &SM3State{sm3.New()}
	default:
		return sha3.NewKeccak256().(HashState)
	}
}

// NewAddress creates an PlatONE address given the bytes and the nonce
func NewAddress(b common.Address, nonce uint64) common.Address {
	data, _ := rlp.EncodeToBytes([]interface{}{b, nonce})
	bytes := DefaultHasher.Hash256(data)[12:]
	return common.BytesToAddress(bytes)
}

// NewAddress2 creates an PlatONE address given the address bytes, initial
// contract code and a salt.
func NewAddress2(b common.Address, salt [32]byte, code []byte) common.Address {
	bytes := DefaultHasher.Hash256([]byte{0xff}, b.Bytes(), salt[:], DefaultHasher.Hash256(code))[12:]
	return common.BytesToAddress(bytes)
}
func CompressPubkeyNew(pubkey PubKey) []byte {
	switch Encryption {
	case GM:
		buf := []byte{}
		yp := pubkey.GetY().Bit(0)
		buf = append(buf, pubkey.GetX().Bytes()...)
		if n := len(pubkey.GetX().Bytes()); n < 32 {
			buf = append(zeroByteSlice()[:(32-n)], buf...)
		}
		buf = append([]byte{byte(yp + 2)}, buf...)

		return buf
	case ECDSA:
		return secp256k1.CompressPubkey(pubkey.GetX(), pubkey.GetY())
	default:
		return nil
	}
}

func DecompressPubkeyNew(pubkey []byte) (PubKey, error) {
	switch Encryption {
	case GM:
		return &GMSMPubKey{pub: sm2.Decompress(pubkey)}, nil
	case ECDSA:
		x, y := secp256k1.DecompressPubkey(pubkey)
		if x == nil {
			return nil, fmt.Errorf("invalid public key")
		}
		rawKey := ecdsa.PublicKey{X: x, Y: y, Curve: S256()}

		return &ECDSAPubKey{pub: &rawKey}, nil
	default:
		return nil, errInvalidPubkeyECDSA
	}
}

// PreHandleTXPayload 对交易数据进行预处理，看需要把公钥放入 payload 的前面
func PreHandleTXPayload(pubkey, payload []byte) []byte {
	switch Encryption {
	case GM:
		return append(pubkey, payload...)
	}
	return payload
}

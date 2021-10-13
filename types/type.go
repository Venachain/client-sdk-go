package types

import (
	"github.com/PlatONE_Network/PlatONE-SDK-Go/platone/common"
	"github.com/PlatONE_Network/PlatONE-SDK-Go/platone/rpc"
	"math/big"
	"sync/atomic"
)

// Client 链 RPC 连接客户端
type Client struct {
	//ethClient   *ethclient.Client
	RpcClient   *rpc.Client
	Passphrase  string
	KeyfilePath string
}

type Response struct {
	Jsonrpc string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
	Id      int         `json:"id"`
}

type IContract interface {
	Deploy() bool
}

type Transaction struct {
	data txdata
	// caches
	hash   atomic.Value
	size   atomic.Value
	from   atomic.Value
	router int32
}

type txdata struct {
	AccountNonce uint64          `json:"nonce"    gencodec:"required"`
	Price        *big.Int        `json:"gasPrice" gencodec:"required"`
	GasLimit     uint64          `json:"gas"      gencodec:"required"`
	Recipient    *common.Address `json:"to"       rlp:"nil"` // nil means contract creation
	Amount       *big.Int        `json:"value"    gencodec:"required"`
	Payload      []byte          `json:"input"    gencodec:"required"`
	//CnsData      []byte          `json:"cnsData"`

	// Signature values
	V *big.Int `json:"v" gencodec:"required"`
	R *big.Int `json:"r" gencodec:"required"`
	S *big.Int `json:"s" gencodec:"required"`

	// This is only used when marshaling to JSON.
	Hash *common.Hash `json:"hash" rlp:"-"`
}

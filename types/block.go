package types

import (
	"math/big"
	"sync/atomic"
	"time"

	"git-c.i.wxblockchain.com/vena/src/client-sdk-go/venachain/common"
	"git-c.i.wxblockchain.com/vena/src/client-sdk-go/venachain/rlp"
)

// 记录有交易依赖关系的交易依赖情况
// 数组中的index用于标识交易在body中的位置
type DAG []Dependency

// 该数组中的内容为直接依赖的交易在body中的位置（不考虑间接依赖关系）
type Dependency []uint

type StorageSize float64
type writeCounter common.StorageSize

type GetBlockResponse struct {
	ParentHash       common.Hash    `json:"parentHash"       gencodec:"required"`
	Coinbase         common.Address `json:"miner"            gencodec:"required"`
	Root             common.Hash    `json:"stateRoot"        gencodec:"required"`
	TransactionsRoot common.Hash    `json:"transactionsRoot" gencodec:"required"`
	ReceiptHash      common.Hash    `json:"receiptsRoot"     gencodec:"required"`
	Bloom            string         `json:"logsBloom"        gencodec:"required"`
	Number           string         `json:"number"           gencodec:"required"`
	GasLimit         string         `json:"gasLimit"         gencodec:"required"`
	GasUsed          string         `json:"gasUsed"          gencodec:"required"`
	Time             string         `json:"timestamp"        gencodec:"required"`
	Extra            string         `json:"extraData"        gencodec:"required"`
	MixDigest        string         `json:"mixHash"          gencodec:"required"`
	Nonce            string         `json:"nonce"            gencodec:"required"`
	Hash             string         `json:"hash"             gencodec:"hash"`
	Transactions     []string       `json:"transactions"     gencodec:"transactions"`
}

// Block represents an entire block in the Ethereum blockchain.
type Block struct {
	Header       *Header
	Transactions Transactions
	//dag          DAG
	// caches
	hash atomic.Value
	size atomic.Value

	// These fields are used by package eth to track
	// inter-peer block relay.
	ReceivedAt   time.Time
	ReceivedFrom interface{}
	ConfirmSigns []*common.BlockConfirmSign
}

// Body is a simple (mutable, non-safe) data container for storing and moving
// a block's data contents (transactions) together.
type Body struct {
	Transactions []*Transaction
}

// NewBlockWithHeader creates a block with the given header data. The
// header data is copied, changes to header and to the field values
// will not affect the block.
func NewBlockWithHeader(header *HeaderVena) *BlockVena {
	return &BlockVena{header: CopyHeader(header)}
}

// CopyHeader creates a deep copy of a block header to prevent side effects from
// modifying a header variable.
func CopyHeader(h *HeaderVena) *HeaderVena {
	cpy := *h
	if cpy.Time = new(big.Int); h.Time != nil {
		cpy.Time.Set(h.Time)
	}
	if cpy.Number = new(big.Int); h.Number != nil {
		cpy.Number.Set(h.Number)
	}
	if len(h.Extra) > 0 {
		cpy.Extra = make([]byte, len(h.Extra))
		copy(cpy.Extra, h.Extra)
	}
	return &cpy
}

// Size returns the true RLP encoded storage size of the block, either by encoding
// and returning it, or returning a previsouly cached value.
func (b *Block) Size() common.StorageSize {
	if size := b.size.Load(); size != nil {
		return size.(common.StorageSize)
	}
	c := writeCounter(0)
	rlp.Encode(&c, b)
	b.size.Store(common.StorageSize(c))
	return common.StorageSize(c)
}

func (c *writeCounter) Write(b []byte) (int, error) {
	*c += writeCounter(len(b))
	return len(b), nil
}

// Block represents an entire block in the Ethereum blockchain.
type BlockVena struct {
	header       *HeaderVena
	transactions Transactions
	//dag          DAG
	// caches
	hash atomic.Value
	size atomic.Value

	// These fields are used by package eth to track
	// inter-peer block relay.
	ReceivedAt   time.Time
	ReceivedFrom interface{}
	ConfirmSigns []*common.BlockConfirmSign
	dag          DAG
}

// Header represents a block header in the Ethereum blockchain.
type Header struct {
	ParentHash  common.Hash    `json:"parentHash"       gencodec:"required"`
	Coinbase    common.Address `json:"miner"            gencodec:"required"`
	Root        common.Hash    `json:"stateRoot"        gencodec:"required"`
	TxHash      common.Hash    `json:"transactionsRoot" gencodec:"required"`
	ReceiptHash common.Hash    `json:"receiptsRoot"     gencodec:"required"`
	Bloom       string         `json:"logsBloom"        gencodec:"required"`
	Number      string         `json:"number"           gencodec:"required"`
	GasLimit    string         `json:"gasLimit"         gencodec:"required"`
	GasUsed     string         `json:"gasUsed"          gencodec:"required"`
	Time        string         `json:"timestamp"        gencodec:"required"`
	Extra       string         `json:"extraData"        gencodec:"required"`
	MixDigest   string         `json:"mixHash"          gencodec:"required"`
	Nonce       string         `json:"nonce"            gencodec:"required"`
	Hash        string         `json:"hash"            gencodec:"hash"`
	// caches
	sealHash atomic.Value `json:"-" rlp:"-"`
}

// Header represents a block header in the Ethereum blockchain.
type HeaderVena struct {
	ParentHash  common.Hash    `json:"parentHash"       gencodec:"required"`
	Coinbase    common.Address `json:"miner"            gencodec:"required"`
	Root        common.Hash    `json:"stateRoot"        gencodec:"required"`
	TxHash      common.Hash    `json:"transactionsRoot" gencodec:"required"`
	ReceiptHash common.Hash    `json:"receiptsRoot"     gencodec:"required"`
	Bloom       Bloom          `json:"logsBloom"        gencodec:"required"`
	Number      *big.Int       `json:"number"           gencodec:"required"`
	GasLimit    uint64         `json:"gasLimit"         gencodec:"required"`
	GasUsed     uint64         `json:"gasUsed"          gencodec:"required"`
	Time        *big.Int       `json:"timestamp"        gencodec:"required"`
	Extra       []byte         `json:"extraData"        gencodec:"required"`
	MixDigest   common.Hash    `json:"mixHash"          gencodec:"required"`
	Nonce       BlockNonce     `json:"nonce"            gencodec:"required"`

	// caches
	sealHash atomic.Value `json:"-" rlp:"-"`
}

// WithBody returns a new block with the given transaction.
func (b *BlockVena) WithBody(transactions []*Transaction, dag DAG) *BlockVena {
	block := &BlockVena{
		header:       CopyHeader(b.header),
		transactions: make([]*Transaction, len(transactions)),
		dag:          dag,
	}
	copy(block.transactions, transactions)
	return block
}

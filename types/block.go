package types

import "git-c.i.wxblockchain.com/vena/src/client-sdk-go/venachain/common"

type GetBlockResponse struct {
	ParentHash   common.Hash    `json:"parentHash"       gencodec:"required"`
	Coinbase     common.Address `json:"miner"            gencodec:"required"`
	Root         common.Hash    `json:"stateRoot"        gencodec:"required"`
	TxHash       common.Hash    `json:"transactionsRoot" gencodec:"required"`
	ReceiptHash  common.Hash    `json:"receiptsRoot"     gencodec:"required"`
	Bloom        string         `json:"logsBloom"        gencodec:"required"`
	Number       string         `json:"number"           gencodec:"required"`
	GasLimit     string         `json:"gasLimit"         gencodec:"required"`
	GasUsed      string         `json:"gasUsed"          gencodec:"required"`
	Time         string         `json:"timestamp"        gencodec:"required"`
	Extra        string         `json:"extraData"        gencodec:"required"`
	MixDigest    string         `json:"mixHash"          gencodec:"required"`
	Nonce        string         `json:"nonce"            gencodec:"required"`
	Hash         string         `json:"hash"             gencodec:"hash"`
	Transactions []string       `json:"transactions"     gencodec:"transactions"`
}

// Body is a simple (mutable, non-safe) data container for storing and moving
// a block's data contents (transactions) together.
type Body struct {
	Transactions []*Transaction
}

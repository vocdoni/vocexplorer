package dbtypes

import (
	"time"
)

// Struct types intended to mirror the database's proto types. This is done so that protobuf is not required by the frontend, drastically reducing the binary size.

// BlockchainInfo mirrors the BlockchainInfo proto type
type BlockchainInfo struct {
	Network           string    `json:"Network,omitempty"`
	Version           string    `json:"Version,omitempty"`
	LatestBlockHeight int64     `json:"LatestBlockHeight,omitempty"`
	GenesisTimeStamp  time.Time `json:"GenesisTimeStamp,omitempty"`
	ChainID           string    `json:"ChainID,omitempty"`
	BlockTime         []int32   `json:"BlockTime,omitempty"`
	BlockTimeStamp    int32     `json:"BlockTimeStamp,omitempty"`
	Height            int64     `json:"Height,omitempty"`
	MaxBytes          int64     `json:"MaxBytes,omitempty"`
	Syncing           bool      `json:"syncing,omitempty"`
}

// Height mirrors the Height proto type
type Height struct {
	Height int64 `json:"Height,omitempty"`
}

type HeightMap struct {
	Heights map[string]int64 `json:"heights,omitempty"`
}

// Envelope mirrors the Envelope proto type
type Envelope struct {
	EncryptionKeyIndexes []int32 `json:"EncryptionKeyIndexes,omitempty"`
	Nullifier            string  `json:"Nullifier,omitempty"`
	ProcessID            string  `json:"ProcessID,omitempty"`
	Package              string  `json:"Package,omitempty"`
	ProcessHeight        int64   `json:"ProcessHeight,omitempty"`
	GlobalHeight         int64   `json:"GlobalHeight,omitempty"`
	TxHeight             int64   `json:"TxHeight,omitempty"`
}

// StoreBlock mirrors the StoreBlock proto type
type StoreBlock struct {
	Hash     []byte    `json:"Hash,omitempty"`
	Height   int64     `json:"Height,omitempty"`
	NumTxs   int64     `json:"NumTxs,omitempty"`
	Time     time.Time `json:"Time,omitempty"`
	Proposer []byte    `json:"Proposer,omitempty"`
}

// Transaction mirrors the Transaction proto type
type Transaction struct {
	Height    int64  `json:"Height,omitempty"`
	Index     uint32 `json:"Index,omitempty"`
	Tx        []byte `json:"Tx,omitempty"`
	TxHeight  int64  `json:"TxHeight,omitempty"`
	Nullifier string `json:"Nullifier,omitempty"`
	Hash      []byte `json:"Hash,omitempty"`
}

// ItemList mirrors the ItemList proto type
type ItemList struct {
	Items [][]byte `json:"items,omitempty"`
}

// Validator mirrors the Validator proto type
type Validator struct {
	Address          []byte  `json:"Address,omitempty"`
	PubKey           []byte  `json:"PubKey,omitempty"`
	VotingPower      int64   `json:"VotingPower,omitempty"`
	ProposerPriority int64   `json:"ProposerPriority,omitempty"`
	Height           *Height `json:"height,omitempty"`
}

// Process mirrors the Process proto type
type Process struct {
	ID          string  `json:"ID,omitempty"`
	EntityID    string  `json:"EntityID,omitempty"`
	LocalHeight *Height `json:"LocalHeight,omitempty"`
}

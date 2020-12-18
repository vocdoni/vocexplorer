package api

import (
	"time"

	"go.vocdoni.io/dvote/types"
)

// Pkeys is the set of cryptographic keys for a process
type Pkeys struct {
	Pub  []types.Key
	Priv []types.Key
	Comm []types.Key
	Rev  []types.Key
}

// VochainStats is the type used by the public stats api
type VochainStats struct {
	BlockHeight       int64     `json:"block_height"`
	EntityCount       int64     `json:"entity_count"`
	EnvelopeCount     int64     `json:"envelope_count"`
	ProcessCount      int64     `json:"process_count"`
	TransactionHeight int64     `json:"transaction_height"`
	ValidatorCount    int64     `json:"validator_count"`
	BlockTime         *[5]int32 `json:"block_time"`
	BlockTimeStamp    int32     `json:"block_time_stamp"`
	ChainID           string    `json:"chain_id"`
	GenesisTimeStamp  time.Time `json:"genesis_time_stamp"`
	Height            int64     `json:"height"`
	Network           string    `json:"network"`
	Version           string    `json:"version"`
	LatestBlockHeight int64     `json:"latest_block_height"`
	AvgTxsPerBlock    float64   `json:"avg_txs_per_block"`
	AvgTxsPerMinute   float64   `json:"avg_txs_per_minute"`
	// The hash of the block with the most txs
	MaxBytes          int64  `json:"max_bytes"`
	MaxTxsBlockHash   string `json:"max_txs_block_hash"`
	MaxTxsBlockHeight int64  `json:"max_txs_block_height"`
	// The start of the minute with the most txs
	MaxTxsMinute    time.Time `json:"max_txs_minute"`
	MaxTxsPerBlock  int64     `json:"max_txs_per_block"`
	MaxTxsPerMinute int64     `json:"max_txs_per_minute"`
	Syncing         bool      `json:"syncing"`
}

// ProcessResults holds the results of a process
type ProcessResults struct {
	State   string
	Type    string
	Results [][]uint32
}

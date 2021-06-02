package storeutil

import (
	"time"

	indexertypes "go.vocdoni.io/dvote/vochain/scrutinizer/indexertypes"
	"go.vocdoni.io/proto/build/go/models"
)

// Transactions stores all data about blockchain transactions
type Transactions struct {
	CurrentTransaction        *indexertypes.TxPackage
	Count                     int
	CurrentDecodedTransaction *DecodedTransaction
	CurrentBlock              *indexertypes.BlockMetadata
	Pagination                PageStore
	Transactions              []*FullTransaction
}

// DecodedTransaction stores human-readable decoded transaction data
type DecodedTransaction struct {
	RawTxContents string
	RawTx         *models.Tx
	Time          time.Time
	ProcessID     string
	EntityID      string
	Nullifier     string
}

// FullTransaction stores a TxPackage and DecodedTransaction
type FullTransaction struct {
	Decoded *DecodedTransaction
	Package *indexertypes.TxPackage
}

package storeutil

import (
	"time"

	indexertypes "go.vocdoni.io/dvote/vochain/scrutinizer/indexertypes"
	"go.vocdoni.io/proto/build/go/models"
)

// Transactions stores all data about blockchain transactions
type Transactions struct {
	CurrentTransaction        *indexertypes.TxPackage
	CurrentDecodedTransaction *DecodedTransaction
	CurrentBlock              *indexertypes.BlockMetadata
	Pagination                PageStore
	Transactions              []*indexertypes.TxMetadata
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

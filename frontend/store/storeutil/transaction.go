package storeutil

import (
	"time"

	"go.vocdoni.io/dvote/types"
	"go.vocdoni.io/proto/build/go/models"
)

// Transactions stores all data about blockchain transactions
type Transactions struct {
	CurrentTransaction        *types.TxPackage
	CurrentDecodedTransaction *DecodedTransaction
	CurrentBlock              *types.BlockMetadata
	Pagination                PageStore
	Transactions              []*types.TxPackage
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

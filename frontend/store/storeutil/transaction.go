package storeutil

import (
	"time"

	"go.vocdoni.io/proto/build/go/models"
)

// Transactions stores all data about blockchain transactions
type Transactions struct {
	CurrentTransaction        *models.TxPackage
	CurrentDecodedTransaction *DecodedTransaction
	CurrentBlock              *models.BlockHeader
	Pagination                PageStore
	Transactions              []*models.TxPackage
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

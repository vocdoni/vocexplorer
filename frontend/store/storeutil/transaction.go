package storeutil

import (
	"time"

	"go.vocdoni.io/proto/build/go/models"
)

// Transactions stores all data about blockchain transactions
type Transactions struct {
	// Count                     int
	CurrentTransactionRef     TransactionReference
	CurrentTransaction        *models.SignedTx
	CurrentDecodedTransaction *DecodedTransaction
	CurrentBlock              *models.BlockHeader
	Pagination                PageStore
	Transactions              []*models.SignedTx
}

// DecodedTransaction stores human-readable decoded transaction data
type DecodedTransaction struct {
	Hash          []byte
	RawTxContents string
	RawTx         models.Tx
	Time          time.Time
	ProcessID     string
	EntityID      string
	Nullifier     string
}

type TransactionReference struct {
	BlockHeight uint32
	TxIndex     int32
}

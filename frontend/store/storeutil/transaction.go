package storeutil

import (
	"time"

	"go.vocdoni.io/proto/build/go/models"
)

// Transactions stores all data about blockchain transactions
type Transactions struct {
	Count                     int
	CurrentTransactionHeight  int64
	CurrentTransaction        *models.SignedTx
	CurrentDecodedTransaction *DecodedTransaction
	CurrentBlock              []byte
	Pagination                PageStore
	Transactions              *models.SignedTxList
}

// DecodedTransaction stores human-readable decoded transaction data
type DecodedTransaction struct {
	EnvelopeHeight int64
	RawTxContents  string
	RawTx          models.Tx
	Time           time.Time
	ProcessID      string
	EntityID       string
	Nullifier      string
}

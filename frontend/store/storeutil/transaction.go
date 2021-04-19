package storeutil

import (
	"time"

	"gitlab.com/vocdoni/vocexplorer/api/dbtypes"
	"gitlab.com/vocdoni/vocexplorer/config"
	"go.vocdoni.io/proto/build/go/models"
)

// Transactions stores all data about blockchain transactions
type Transactions struct {
	Count                     int
	CurrentTransactionHeight  int64
	CurrentTransaction        *dbtypes.Transaction
	CurrentDecodedTransaction *DecodedTransaction
	CurrentBlock              *dbtypes.StoreBlock
	Pagination                PageStore
	Transactions              [config.ListSize]*dbtypes.Transaction
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

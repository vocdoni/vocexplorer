package storeutil

import (
	"time"

	sctypes "go.vocdoni.io/dvote/vochain/scrutinizer/types"
	"go.vocdoni.io/proto/build/go/models"
)

// Transactions stores all data about blockchain transactions
type Transactions struct {
	CurrentTransaction        *sctypes.TxPackage
	CurrentDecodedTransaction *DecodedTransaction
	CurrentBlock              *sctypes.BlockMetadata
	Pagination                PageStore
	Transactions              []*sctypes.TxMetadata
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

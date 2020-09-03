package storeutil

import (
	"time"

	dvotetypes "gitlab.com/vocdoni/go-dvote/types"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/types"
)

// Transactions stores all data about blockchain transactions
type Transactions struct {
	Count                     int
	CurrentTransactionHeight  int64
	CurrentTransaction        *types.SendTx
	CurrentDecodedTransaction *DecodedTransaction
	CurrentBlock              *types.StoreBlock
	Pagination                PageStore
	Transactions              [config.ListSize]*types.SendTx
}

// DecodedTransaction stores human-readable decoded transaction data
type DecodedTransaction struct {
	EnvelopeHeight int64
	Metadata       []byte
	RawTxContents  []byte
	RawTx          dvotetypes.Tx
	Time           time.Time
	ProcessID      string
	EntityID       string
	Nullifier      string
}

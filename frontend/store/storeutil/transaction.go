package storeutil

import (
	"time"

	"gitlab.com/vocdoni/vocexplorer/api/dvotetypes"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/proto"
)

// Transactions stores all data about blockchain transactions
type Transactions struct {
	Count                     int
	CurrentTransactionHeight  int64
	CurrentTransaction        *proto.SendTx
	CurrentDecodedTransaction *DecodedTransaction
	CurrentBlock              *proto.StoreBlock
	Pagination                PageStore
	Transactions              [config.ListSize]*proto.SendTx
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

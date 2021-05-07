package storeutil

import indexertypes "go.vocdoni.io/dvote/vochain/scrutinizer/indexertypes"

// Blocks stores all data about blockchain blocks
type Blocks struct {
	Blocks                []*indexertypes.BlockMetadata
	Count                 int
	CurrentBlock          *indexertypes.BlockMetadata
	CurrentBlockHeight    uint32
	CurrentTxs            []*indexertypes.TxMetadata
	TransactionPagination PageStore
	Pagination            PageStore
}

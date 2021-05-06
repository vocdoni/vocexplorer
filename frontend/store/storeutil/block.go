package storeutil

import sctypes "go.vocdoni.io/dvote/vochain/scrutinizer/types"

// Blocks stores all data about blockchain blocks
type Blocks struct {
	Blocks                []*sctypes.BlockMetadata
	Count                 int
	CurrentBlock          *sctypes.BlockMetadata
	CurrentBlockHeight    uint32
	CurrentTxs            []*sctypes.TxMetadata
	TransactionPagination PageStore
	Pagination            PageStore
}

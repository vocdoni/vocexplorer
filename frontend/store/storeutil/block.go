package storeutil

import (
	"go.vocdoni.io/dvote/types"
)

// Blocks stores all data abotu blockchain blocks
type Blocks struct {
	Blocks                []*types.BlockMetadata
	Count                 int
	CurrentBlock          *types.BlockMetadata
	CurrentBlockHeight    uint32
	CurrentTxs            []*types.TxPackage
	TransactionPagination PageStore
	Pagination            PageStore
}

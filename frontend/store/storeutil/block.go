package storeutil

import (
	"go.vocdoni.io/proto/build/go/models"
)

// Blocks stores all data abotu blockchain blocks
type Blocks struct {
	Blocks                []*models.BlockHeader
	Count                 int
	CurrentBlock          *models.BlockHeader
	CurrentBlockHeight    uint32
	CurrentTxs            []*models.TxPackage
	TransactionPagination PageStore
	Pagination            PageStore
}

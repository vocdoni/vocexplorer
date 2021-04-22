package storeutil

import (
	"go.vocdoni.io/proto/build/go/models"
)

// Blocks stores all data abotu blockchain blocks
type Blocks struct {
	Blocks                []*models.BlockHeader
	Count                 int
	CurrentBlock          *models.BlockHeader
	CurrentTxs            *models.SignedTxList
	TransactionPagination PageStore
	Pagination            PageStore
}

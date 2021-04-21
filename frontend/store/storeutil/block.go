package storeutil

import (
	"go.vocdoni.io/proto/build/go/models"
)

// Blocks stores all data abotu blockchain blocks
type Blocks struct {
	Blocks                *models.TendermintHeaderList
	Count                 int
	CurrentBlock          *models.TendermintHeader
	CurrentTxs            *models.SignedTxList
	TransactionPagination PageStore
	Pagination            PageStore
}

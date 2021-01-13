package storeutil

import (
	"github.com/vocdoni/vocexplorer/api"
	"github.com/vocdoni/vocexplorer/api/dbtypes"
	"github.com/vocdoni/vocexplorer/config"
)

// Blocks stores all data abotu blockchain blocks
type Blocks struct {
	Blocks                [config.ListSize]*dbtypes.StoreBlock
	Count                 int
	CurrentBlock          *api.Block
	CurrentBlockHeight    int64
	CurrentTxs            [config.ListSize]*dbtypes.Transaction
	TransactionPagination PageStore
	Pagination            PageStore
}

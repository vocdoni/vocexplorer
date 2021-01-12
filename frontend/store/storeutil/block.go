package storeutil

import (
	"gitlab.com/vocdoni/vocexplorer/api"
	"gitlab.com/vocdoni/vocexplorer/api/dbtypes"
	"gitlab.com/vocdoni/vocexplorer/config"
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

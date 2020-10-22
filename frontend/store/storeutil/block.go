package storeutil

import (
	"gitlab.com/vocdoni/vocexplorer/api/dbtypes"
	"gitlab.com/vocdoni/vocexplorer/api/tmtypes"
	"gitlab.com/vocdoni/vocexplorer/config"
)

// Blocks stores all data abotu blockchain blocks
type Blocks struct {
	Blocks                [config.ListSize]*dbtypes.StoreBlock
	Count                 int
	CurrentBlock          *tmtypes.ResultBlock
	CurrentBlockHeight    int64
	CurrentTxs            [config.ListSize]*dbtypes.Transaction
	TransactionPagination PageStore
	Pagination            PageStore
}

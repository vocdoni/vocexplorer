package storeutil

import (
	tmtypes "github.com/tendermint/tendermint/types"
	"gitlab.com/vocdoni/vocexplorer/api/dbtypes"
	"gitlab.com/vocdoni/vocexplorer/config"
)

// Blocks stores all data abotu blockchain blocks
type Blocks struct {
	Blocks                [config.ListSize]*dbtypes.StoreBlock
	Count                 int
	CurrentBlock          *tmtypes.Block
	CurrentBlockHeight    int64
	CurrentTxs            [config.ListSize]*dbtypes.Transaction
	TransactionPagination PageStore
	Pagination            PageStore
}

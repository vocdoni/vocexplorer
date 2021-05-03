package storeutil

import (
	tmtypes "github.com/tendermint/tendermint/types"
	"go.vocdoni.io/dvote/types"
)

// Blocks stores all data abotu blockchain blocks
type Blocks struct {
	Blocks                []*tmtypes.Block
	Count                 int
	CurrentBlock          *tmtypes.Block
	CurrentBlockHeight    uint32
	CurrentTxs            []*types.TxPackage
	TransactionPagination PageStore
	Pagination            PageStore
}

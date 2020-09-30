package storeutil

import (
	"gitlab.com/vocdoni/vocexplorer/api/tmtypes"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/proto"
)

// Blocks stores all data abotu blockchain blocks
type Blocks struct {
	Blocks                [config.ListSize]*proto.StoreBlock
	Count                 int
	CurrentBlock          *tmtypes.ResultBlock
	CurrentBlockHeight    int64
	CurrentBlockTxHeights []int64
	Pagination            PageStore
}

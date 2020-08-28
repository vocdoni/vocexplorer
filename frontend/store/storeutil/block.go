package storeutil

import (
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/types"
)

// Blocks stores all data abotu blockchain blocks
type Blocks struct {
	BlockList   [config.ListSize]*types.StoreBlock
	Pagination  PageStore
	TotalBlocks int
}

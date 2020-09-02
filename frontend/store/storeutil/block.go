package storeutil

import (
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/types"
)

// Blocks stores all data abotu blockchain blocks
type Blocks struct {
	Blocks             [config.ListSize]*types.StoreBlock
	Count              int
	CurrentBlock       *coretypes.ResultBlock
	CurrentBlockHeight int64
	Pagination         PageStore
}

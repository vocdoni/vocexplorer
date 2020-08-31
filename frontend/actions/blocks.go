package actions

import (
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/types"
)

// BlocksTabChange is the action to change the current blocks tab
type BlocksTabChange struct {
	Tab string
}

// SetBlockList is the action to set the list of current blocks
type SetBlockList struct {
	BlockList [config.ListSize]*types.StoreBlock
}

// BlocksHeightUpdate is the action to change the current block height
type BlocksHeightUpdate struct {
	Height int
}

// SetCurrentBlock is the action to set the current block
type SetCurrentBlock struct {
	Block *coretypes.ResultBlock
}

// DisableBlockUpdate is the action to set the disable update status for blocks
type DisableBlockUpdate struct {
	Disabled bool
}

// On initialization, register actions
func init() {
	dispatcher.Register(blockActions)
}

// blockActions is the handler for all block-related store actions
func blockActions(action interface{}) {
	switch a := action.(type) {
	case *BlocksTabChange:
		store.Blocks.Pagination.Tab = a.Tab

	case *BlocksHeightUpdate:
		store.Blocks.Count = a.Height

	case *SetBlockList:
		store.Blocks.Blocks = a.BlockList

	case *SetCurrentBlock:
		store.Blocks.CurrentBlock = a.Block

	case *DisableBlockUpdate:
		store.Blocks.Pagination.DisableUpdate = a.Disabled
	default:
		return // don't fire listeners
	}

	store.Listeners.Fire()
}

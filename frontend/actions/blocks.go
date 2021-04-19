package actions

import (
	"gitlab.com/vocdoni/vocexplorer/api"
	"gitlab.com/vocdoni/vocexplorer/api/dbtypes"
	"gitlab.com/vocdoni/vocexplorer/config"
)

// BlocksIndexChange is the action to set the pagination index
type BlocksIndexChange struct {
	Index int
}

// BlockTransactionsIndexChange is the action to set the pagination index
type BlockTransactionsIndexChange struct {
	Index int
}

// BlocksTabChange is the action to change the current blocks tab
type BlocksTabChange struct {
	Tab string
}

// SetBlockList is the action to set the list of current blocks
type SetBlockList struct {
	BlockList [config.ListSize]*dbtypes.StoreBlock
}

// BlocksHeightUpdate is the action to change the current block height
type BlocksHeightUpdate struct {
	Height int
}

// SetCurrentBlock is the action to set the current block
type SetCurrentBlock struct {
	Block *api.Block
}

// SetCurrentBlockHeight is the action to set the current block height
type SetCurrentBlockHeight struct {
	Height int64
}

// SetCurrentBlockTransactionList is the action to set the current list of transactions
type SetCurrentBlockTransactionList struct {
	TransactionList [config.ListSize]*dbtypes.Transaction
}

package actions

import (
	"go.vocdoni.io/dvote/types"
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
	BlockList []*types.BlockMetadata
}

// BlocksHeightUpdate is the action to change the current block height
type BlocksHeightUpdate struct {
	Height int
}

// SetCurrentBlock is the action to set the current block
type SetCurrentBlock struct {
	Block *types.BlockMetadata
}

// SetCurrentBlockHeight is the action to set the current block height
type SetCurrentBlockHeight struct {
	Height uint32
}

// SetCurrentBlockTransactionList is the action to set the current list of transactions
type SetCurrentBlockTransactionList struct {
	TransactionList []*types.TxPackage
}

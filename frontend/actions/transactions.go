package actions

import (
	"gitlab.com/vocdoni/vocexplorer/frontend/store/storeutil"
	"go.vocdoni.io/dvote/types"
)

// TransactionTabChange is the action to change between tabs in transaction view details
type TransactionTabChange struct {
	Tab string
}

// TransactionsIndexChange is the action to set the pagination index
type TransactionsIndexChange struct {
	Index int
}

// SetTransactionList is the action to set the transaction list
type SetTransactionList struct {
	TransactionList []*types.TxPackage
}

// SetCurrentTransaction is the action to set the current transaction
type SetCurrentTransaction struct {
	Transaction *types.TxPackage
}

// SetTransactionBlock is the action to set the block associated with the current transaction
type SetTransactionBlock struct {
	Block *types.BlockMetadata
}

// SetCurrentDecodedTransaction is the action to set the decoded contents associated with the current transaction
type SetCurrentDecodedTransaction struct {
	Transaction *storeutil.DecodedTransaction
}

package actions

import (
	"github.com/vocdoni/vocexplorer/api/dbtypes"
	"github.com/vocdoni/vocexplorer/config"
	"github.com/vocdoni/vocexplorer/frontend/store/storeutil"
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
	TransactionList [config.ListSize]*dbtypes.Transaction
}

// SetTransactionCount is the action to set the Transaction count
type SetTransactionCount struct {
	Count int
}

// SetCurrentTransactionHeight is the action to set the height of the current transaction
type SetCurrentTransactionHeight struct {
	Height int64
}

// SetCurrentTransaction is the action to set the current transaction
type SetCurrentTransaction struct {
	Transaction *dbtypes.Transaction
}

// SetTransactionBlock is the action to set the block associated with the current transaction
type SetTransactionBlock struct {
	Block *dbtypes.StoreBlock
}

// SetCurrentDecodedTransaction is the action to set the decoded contents associated with the current transaction
type SetCurrentDecodedTransaction struct {
	Transaction *storeutil.DecodedTransaction
}

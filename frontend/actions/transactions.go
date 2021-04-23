package actions

import (
	"gitlab.com/vocdoni/vocexplorer/frontend/store/storeutil"
	"go.vocdoni.io/proto/build/go/models"
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
	TransactionList []*models.SignedTx
}

// // SetTransactionCount is the action to set the Transaction count
// type SetTransactionCount struct {
// 	Count int
// }

// SetCurrentTransactionRef is the action to set the reference of the current transaction
type SetCurrentTransactionRef struct {
	Height uint32
	Index  int32
}

// SetCurrentTransaction is the action to set the current transaction
type SetCurrentTransaction struct {
	Transaction *models.SignedTx
}

// SetTransactionBlock is the action to set the block associated with the current transaction
type SetTransactionBlock struct {
	Block *models.BlockHeader
}

// SetCurrentDecodedTransaction is the action to set the decoded contents associated with the current transaction
type SetCurrentDecodedTransaction struct {
	Transaction *storeutil.DecodedTransaction
}

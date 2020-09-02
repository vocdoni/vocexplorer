package actions

import (
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/frontend/store/storeutil"
	"gitlab.com/vocdoni/vocexplorer/types"
)

// SetTransactionList is the action to set the transaction list
type SetTransactionList struct {
	TransactionList [config.ListSize]*types.SendTx
}

// SetTransactionCount is the action to set the Transaction count
type SetTransactionCount struct {
	Count int
}

// DisableTransactionsUpdate is the action to set the disable update status for transactions
type DisableTransactionsUpdate struct {
	Disabled bool
}

// SetCurrentTransactionHeight is the action to set the height of the current transaction
type SetCurrentTransactionHeight struct {
	Height int64
}

// SetCurrentTransaction is the action to set the current transaction
type SetCurrentTransaction struct {
	Transaction *types.SendTx
}

// SetTransactionBlock is the action to set the block associated with the current transaction
type SetTransactionBlock struct {
	Block *types.StoreBlock
}

// SetCurrentDecodedTransaction is the action to set the decoded contents associated with the current transaction
type SetCurrentDecodedTransaction struct {
	Transaction *storeutil.DecodedTransaction
}

// On initialization, register actions
func init() {
	dispatcher.Register(transactionActions)
}

// transactionActions is the handler for all transaction-related store actions
func transactionActions(action interface{}) {
	switch a := action.(type) {
	case *SetTransactionList:
		store.Transactions.Transactions = a.TransactionList

	case *SetTransactionCount:
		store.Transactions.Count = a.Count

	case *DisableTransactionsUpdate:
		store.Transactions.Pagination.DisableUpdate = a.Disabled

	case *SetCurrentTransactionHeight:
		store.Transactions.CurrentTransactionHeight = a.Height

	case *SetCurrentTransaction:
		store.Transactions.CurrentTransaction = a.Transaction

	case *SetTransactionBlock:
		store.Transactions.CurrentBlock = a.Block

	case *SetCurrentDecodedTransaction:
		store.Transactions.CurrentDecodedTransaction = a.Transaction

	default:
		return // don't fire listeners
	}

	store.Listeners.Fire()
}

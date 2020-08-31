package actions

import (
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
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

	default:
		return // don't fire listeners
	}

	store.Listeners.Fire()
}

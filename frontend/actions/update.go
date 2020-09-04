package actions

import (
	"gitlab.com/vocdoni/vocexplorer/api"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
)

// UpdateCounts updates the values of all item counts (eg. validator count)
func UpdateCounts() {
	newVal, ok := api.GetBlockHeight()
	if ok {
		dispatcher.Dispatch(&BlocksHeightUpdate{Height: int(newVal)})
	}
	newVal, ok = api.GetTxHeight()
	if ok {
		dispatcher.Dispatch(&SetTransactionCount{Count: int(newVal)})
	}
	newVal, ok = api.GetEntityHeight()
	if ok {
		dispatcher.Dispatch(&SetEntityCount{Count: int(newVal)})
	}
	newVal, ok = api.GetProcessHeight()
	if ok {
		dispatcher.Dispatch(&SetProcessCount{Count: int(newVal)})
	}
	newVal, ok = api.GetEnvelopeHeight()
	if ok {
		dispatcher.Dispatch(&SetEnvelopeCount{Count: int(newVal)})
	}
	newVal, ok = api.GetValidatorCount()
	if ok {
		dispatcher.Dispatch(&SetValidatorCount{Count: int(newVal)})
	}
}

// EnableUpdates resets all components' 'disable update' flags
func EnableUpdates() {
	dispatcher.Dispatch(&DisableBlockUpdate{Disabled: false})
	dispatcher.Dispatch(&DisableEntityUpdate{Disabled: false})
	dispatcher.Dispatch(&DisableEnvelopeUpdate{Disabled: false})
	dispatcher.Dispatch(&DisableTransactionsUpdate{Disabled: false})
	dispatcher.Dispatch(&DisableValidatorUpdate{Disabled: false})
}

package actions

import (
	"gitlab.com/vocdoni/vocexplorer/api"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
)

// DisableUpdate is the action to set the disable update status for given disableupdate boolean
type DisableUpdate struct {
	Updater  *bool
	Disabled bool
}

// EnableAllUpdates is the action to set all disable updates bools to false
type EnableAllUpdates struct {
}

// UpdateCounts updates the values of all item counts (eg. validator count)
func UpdateCounts() {
	newVal, ok := api.GetBlockHeight()
	if ok {
		// Blocks indexed by height, so convert to count
		dispatcher.Dispatch(&BlocksHeightUpdate{Height: int(newVal) - 1})
	}
	newVal, ok = api.GetTxHeight()
	if ok {
		// Transactions indexed by height, so convert to count
		dispatcher.Dispatch(&SetTransactionCount{Count: int(newVal) - 1})
	}
	newVal, ok = api.GetEntityCount()
	if ok {
		dispatcher.Dispatch(&SetEntityCount{Count: int(newVal)})
	}
	newVal, ok = api.GetProcessCount()
	if ok {
		dispatcher.Dispatch(&SetProcessCount{Count: int(newVal)})
	}
	newVal, ok = api.GetEnvelopeCount()
	if ok {
		dispatcher.Dispatch(&SetEnvelopeCount{Count: int(newVal)})
	}
	newVal, ok = api.GetValidatorCount()
	if ok {
		dispatcher.Dispatch(&SetValidatorCount{Count: int(newVal)})
	}
}

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
	newVal, ok := api.GetStats()
	if ok {
		dispatcher.Dispatch(&BlocksHeightUpdate{Height: int(newVal.BlockHeight) - 1})
		dispatcher.Dispatch(&SetTransactionCount{Count: int(newVal.TransactionHeight) - 1})
		dispatcher.Dispatch(&SetEntityCount{Count: int(newVal.EntityCount)})
		dispatcher.Dispatch(&SetEnvelopeCount{Count: int(newVal.EnvelopeCount)})
		dispatcher.Dispatch(&SetProcessCount{Count: int(newVal.ProcessCount)})
		dispatcher.Dispatch(&SetValidatorCount{Count: int(newVal.ValidatorCount)})
		dispatcher.Dispatch(&SetTransactionStats{
			AvgTxsPerBlock:    newVal.AvgTxsPerBlock,
			AvgTxsPerMinute:   newVal.AvgTxsPerMinute,
			MaxTxsBlockHash:   newVal.MaxTxsBlockHash,
			MaxTxsBlockHeight: newVal.MaxTxsBlockHeight,
			MaxTxsMinute:      newVal.MaxTxsMinute,
			MaxTxsPerBlock:    newVal.MaxTxsPerBlock,
			MaxTxsPerMinute:   newVal.MaxTxsPerMinute,
		})

	}
	// newVal, ok := api.GetBlockHeight()
	// if ok {
	// 	// Blocks indexed by height, so convert to count
	// 	dispatcher.Dispatch(&BlocksHeightUpdate{Height: int(newVal) - 1})
	// }
	// newVal, ok = api.GetTxHeight()
	// if ok {
	// 	// Transactions indexed by height, so convert to count
	// 	dispatcher.Dispatch(&SetTransactionCount{Count: int(newVal) - 1})
	// }
	// newVal, ok = api.GetEntityCount()
	// if ok {
	// 	dispatcher.Dispatch(&SetEntityCount{Count: int(newVal)})
	// }
	// newVal, ok = api.GetProcessCount()
	// if ok {
	// 	dispatcher.Dispatch(&SetProcessCount{Count: int(newVal)})
	// }
	// newVal, ok = api.GetEnvelopeCount()
	// if ok {
	// 	dispatcher.Dispatch(&SetEnvelopeCount{Count: int(newVal)})
	// }
	// newVal, ok = api.GetValidatorCount()
	// if ok {
	// 	dispatcher.Dispatch(&SetValidatorCount{Count: int(newVal)})
	// }
}

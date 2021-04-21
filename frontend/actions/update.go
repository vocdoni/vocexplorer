package actions

// DisableUpdate is the action to set the disable update status for given disableupdate boolean
type DisableUpdate struct {
	Updater  *bool
	Disabled bool
}

// EnableAllUpdates is the action to set all disable updates bools to false
type EnableAllUpdates struct {
}

// // UpdateCounts updates the values of all item counts (eg. validator count)
// func UpdateCounts() {
// 	stats, ok := api.GetStats()
// 	if ok {
// 		dispatcher.Dispatch(&BlocksHeightUpdate{Height: int(stats.BlockHeight) - 1})
// 		dispatcher.Dispatch(&SetTransactionCount{Count: int(stats.TransactionHeight) - 1})
// 		dispatcher.Dispatch(&SetEntityCount{Count: int(stats.EntityCount)})
// 		dispatcher.Dispatch(&SetEnvelopeCount{Count: int(stats.EnvelopeCount)})
// 		dispatcher.Dispatch(&SetProcessCount{Count: int(stats.ProcessCount)})
// 		dispatcher.Dispatch(&SetValidatorCount{Count: int(stats.ValidatorCount)})
// 		dispatcher.Dispatch(&SetStats{Stats: stats})
// 	}
// }

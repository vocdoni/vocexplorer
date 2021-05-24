package update

import (
	"strings"
	"time"

	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/frontend/store/storeutil"
	"gitlab.com/vocdoni/vocexplorer/logger"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// EnvelopeProcessResults updates auxilary info for all process id's belonging to currently displayed envelopes
func EnvelopeProcessResults() {
	for _, envelope := range store.Envelopes.Envelopes {
		if envelope != nil {
			ID := strings.ToLower(util.TrimHex(util.HexToString(envelope.ProcessId)))
			if ID != "" {
				if _, ok := store.Processes.ProcessResults[ID]; !ok {
					results, state, tp, final, err := store.Client.GetResults(envelope.ProcessId)
					if err != nil {
						logger.Error(err)
					}
					if err == nil {
						dispatcher.Dispatch(&actions.SetProcessResults{
							PID: ID,
							Results: storeutil.ProcessResults{
								Results: results,
								State:   state,
								Type:    tp,
								Final:   final,
							},
						})
					}
				}
			}
		}
	}
}

// CurrentProcessResults updates current process information
func CurrentProcessResults() {
	results, state, tp, final, err := store.Client.GetResults(store.Processes.CurrentProcess.Process.ID)
	if err != nil {
		logger.Error(err)
		return
	}
	dispatcher.Dispatch(&actions.SetProcessResults{
		PID: util.HexToString(store.Processes.CurrentProcess.Process.ID),
		Results: storeutil.ProcessResults{
			Results: results,
			State:   state,
			Type:    tp,
			Final:   final,
		},
	})

}

// CheckCurrentPage returns true and stops ticker if the current page is title
func CheckCurrentPage(title string, ticker *time.Ticker) bool {
	if store.CurrentPage != title {
		logger.Info("redirecting")
		ticker.Stop()
		return false
	}
	return true
}

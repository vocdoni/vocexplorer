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

// ProcessResults updates auxilary info for all currently displayed process id's
func ProcessResults() {
	for _, process := range store.Processes.Processes {
		if process != nil {
			results, state, tp, final, err := store.Client.GetResults(process.Process.ID)
			if err != nil {
				logger.Error(err)
			}
			if results != nil && err == nil {
				dispatcher.Dispatch(&actions.SetProcessResults{
					PID: util.HexToString(process.Process.ID),
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
					if results != nil && err == nil {
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

// EntityProcessResults ensures the given entity's processes' results are all stored
func EntityProcessResults() {
	for _, pid := range store.Entities.CurrentEntity.ProcessIds {
		if _, ok := store.Processes.ProcessResults[pid]; !ok {
			results, state, tp, final, err := store.Client.GetResults(util.StringToHex(pid))
			if err != nil {
				logger.Error(err)
				return
			}
			if results != nil {
				dispatcher.Dispatch(&actions.SetProcessResults{
					PID: pid,
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

// CheckCurrentPage returns true and stops ticker if the current page is title
func CheckCurrentPage(title string, ticker *time.Ticker) bool {
	if store.CurrentPage != title {
		logger.Info("redirecting")
		ticker.Stop()
		return false
	}
	return true
}

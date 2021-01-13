package update

import (
	"strings"
	"time"

	"github.com/vocdoni/vocexplorer/api"
	"github.com/vocdoni/vocexplorer/frontend/actions"
	"github.com/vocdoni/vocexplorer/frontend/dispatcher"
	"github.com/vocdoni/vocexplorer/frontend/store"
	"github.com/vocdoni/vocexplorer/frontend/store/storeutil"
	"github.com/vocdoni/vocexplorer/logger"
	"github.com/vocdoni/vocexplorer/util"
)

// ProcessResults updates auxilary info for all currently displayed process id's
func ProcessResults() {
	for _, process := range store.Processes.Processes {
		if process != nil {
			ID := process.ID
			if ID != "" {
				if _, ok := store.Processes.ProcessResults[ID]; !ok {
					results, ok := api.GetProcessResults(strings.ToLower(ID))
					if ok && results != nil {
						dispatcher.Dispatch(&actions.SetProcessContents{
							ID: ID,
							Process: storeutil.Process{
								ProcessType: results.Type,
								State:       results.State,
								Results:     results.Results},
						})
					}
				}
			}
		}
	}
}

// EnvelopeProcessResults updates auxilary info for all process id's belonging to currently displayed envelopes
func EnvelopeProcessResults() {
	for _, envelope := range store.Envelopes.Envelopes {
		if envelope != nil {
			ID := strings.ToLower(util.TrimHex(envelope.ProcessID))
			if ID != "" {
				if _, ok := store.Processes.ProcessResults[ID]; !ok {
					results, ok := api.GetProcessResults(strings.ToLower(ID))
					if ok && results != nil {
						dispatcher.Dispatch(&actions.SetProcessContents{
							ID: ID,
							Process: storeutil.Process{
								ProcessType: results.Type,
								State:       results.State,
								Results:     results.Results},
						})
					}
				}
			}
		}
	}
}

// CurrentProcessResults updates current process information
func CurrentProcessResults() {
	results, ok := api.GetProcessResults(strings.ToLower(store.Processes.CurrentProcess.ID))
	if !ok || results == nil {
		return
	}
	newVal, ok := api.GetProcessEnvelopeCount(store.Processes.CurrentProcess.ID)
	if !ok {
		return
	}
	if ok && results != nil {
		dispatcher.Dispatch(&actions.SetCurrentProcessResults{
			Process: storeutil.Process{
				EnvelopeCount: int(newVal),
				ProcessType:   results.Type,
				State:         results.State,
				Results:       results.Results},
		})
	}
}

// EntityProcessResults ensures the given entity's processes' results are all stored
func EntityProcessResults() {
	for _, process := range store.Entities.CurrentEntity.Processes {
		if process != nil {
			ID := process.ID
			if ID != "" {
				if _, ok := store.Processes.ProcessResults[ID]; !ok {
					results, ok := api.GetProcessResults(strings.ToLower(ID))
					if ok && results != nil {
						dispatcher.Dispatch(&actions.SetProcessContents{
							ID: ID,
							Process: storeutil.Process{
								ProcessType: results.Type,
								State:       results.State,
								Results:     results.Results},
						})
					}
				}
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

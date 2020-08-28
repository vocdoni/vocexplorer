package actions

import (
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/frontend/store/storeutil"
)

// ProcessesTabChange is the action to change processes tabs
type ProcessesTabChange struct {
	Tab string
}

// SetProcessIDs is the action to set the process list
type SetProcessIDs struct {
	ProcessIDs [config.ListSize]string
}

// SetProcessCount is the action to set the process count
type SetProcessCount struct {
	ProcessCount int
}

// SetEnvelopeHeights is the action to set the envelope heights map
type SetEnvelopeHeights struct {
	EnvelopeHeights map[string]int64
}

// SetProcessContents is the action to set a single process contents
type SetProcessContents struct {
	Process storeutil.Process
	ID      string
}

// On initialization, register actions
func init() {
	dispatcher.Register(envelopeActions)
}

// processActions is the handler for all process-related store actions
func processActions(action interface{}) {
	switch a := action.(type) {
	case *SetProcessIDs:
		store.Processes.ProcessIDs = a.ProcessIDs

	case *ProcessesTabChange:
		store.Processes.Pagination.Tab = a.Tab

	case *SetProcessCount:
		store.Processes.ProcessCount = a.ProcessCount

	case *SetEnvelopeHeights:
		store.Processes.EnvelopeHeights = a.EnvelopeHeights

	case *SetProcessContents:
		store.Processes.ProcessResults[a.ID] = a.Process

	default:
		return // don't fire listeners
	}

	store.Listeners.Fire()
}

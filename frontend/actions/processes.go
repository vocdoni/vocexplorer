package actions

import (
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/api"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/frontend/store/storeutil"
	"gitlab.com/vocdoni/vocexplorer/types"
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
	Count int
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

// SetProcessKeys is the action to set the keys for a single process
type SetProcessKeys struct {
	Keys *api.Pkeys
	ID   string
}

// DisableProcessUpdate is the action to set the disable update status for processes
type DisableProcessUpdate struct {
	Disabled bool
}

// SetProcessState is the action to set the current process state
type SetProcessState struct {
	State string
}

// SetProcessType is the action to set the current process type
type SetProcessType struct {
	Type string
}

// SetCurrentProcessEnvelopeHeight is the action to set the current process' envelope height
type SetCurrentProcessEnvelopeHeight struct {
	Height int
}

// SetCurrentProcess is the action to set the current process
type SetCurrentProcess struct {
	Process storeutil.Process
}

// SetCurrentProcessID is the action to set the current process ID
type SetCurrentProcessID struct {
	ID string
}

// SetCurrentProcessEnvelopes is the action to set the envelope list for the current process
type SetCurrentProcessEnvelopes struct {
	EnvelopeList [config.ListSize]*types.Envelope
}

// On initialization, register actions
func init() {
	dispatcher.Register(processActions)
}

// processActions is the handler for all process-related store actions
func processActions(action interface{}) {
	switch a := action.(type) {
	case *SetProcessIDs:
		store.Processes.ProcessIDs = a.ProcessIDs

	case *ProcessesTabChange:
		store.Processes.Pagination.Tab = a.Tab

	case *SetProcessCount:
		store.Processes.Count = a.Count

	case *SetEnvelopeHeights:
		store.Processes.EnvelopeHeights = a.EnvelopeHeights

	case *SetProcessContents:
		store.Processes.ProcessResults[a.ID] = a.Process

	case *SetProcessKeys:
		store.Processes.ProcessKeys[a.ID] = a.Keys

	case *DisableProcessUpdate:
		store.Processes.Pagination.DisableUpdate = a.Disabled

	case *SetProcessState:
		store.Processes.CurrentProcess.State = a.State

	case *SetProcessType:
		store.Processes.CurrentProcess.ProcessType = a.Type

	case *SetCurrentProcessEnvelopeHeight:
		store.Processes.CurrentProcess.EnvelopeCount = a.Height

	case *SetCurrentProcess:
		store.Processes.CurrentProcess = a.Process

	case *SetCurrentProcessID:
		store.Processes.CurrentProcessID = a.ID

	case *SetCurrentProcessEnvelopes:
		store.Processes.CurrentProcess.Envelopes = a.EnvelopeList

	default:
		return // don't fire listeners
	}

	store.Listeners.Fire()
}

package actions

import (
	"gitlab.com/vocdoni/vocexplorer/frontend/store/storeutil"
	indexertypes "go.vocdoni.io/dvote/vochain/scrutinizer/indexertypes"
)

// ProcessesIndexChange is the action to set the pagination index
type ProcessesIndexChange struct {
	Index int
}

// ProcessEnvelopesIndexChange is the action to set the pagination index for the current process' envelope list
type ProcessEnvelopesIndexChange struct {
	Index int
}

// ProcessTabChange is the action to change processes tabs
type ProcessTabChange struct {
	Tab string
}

// SetProcess is the action to set a single process
type SetProcess struct {
	PID     string
	Process *storeutil.Process
}

// SetProcessIds is the action to set the process ids
type SetProcessIds struct {
	Processes []string
}

// SetProcessCount is the action to set the process count
type SetProcessCount struct {
	Count int
}

// SetProcessResults is the action to set a single process results
type SetProcessResults struct {
	Results storeutil.ProcessResults
	PID     string
}

// SetProcessState is the action to set the current process state
type SetProcessState struct {
	State string
}

// SetProcessType is the action to set the current process type
type SetProcessType struct {
	Type string
}

// SetCurrentProcessStruct is the action to set the current process
type SetCurrentProcessStruct struct {
	Process *storeutil.Process
}

// SetCurrentProcessEnvelopeCount is the action to set the current process envelope cou t
type SetCurrentProcessEnvelopeCount struct {
	Count int
}

// SetCurrentProcessEnvelopes is the action to set the envelope list for the current process
type SetCurrentProcessEnvelopes struct {
	EnvelopeList []*indexertypes.EnvelopeMetadata
}

// SetProcessStatusFilter sets the status filter for processes
type SetProcessStatusFilter struct {
	StatusFilter string
}

// SetProcessSrcNetworkIDFilter sets the source network id filter for processes
type SetProcessSrcNetworkIDFilter struct {
	SrcNetworkIDFilter string
}

// SetProcessResultsFilter sets the results filter for processes
type SetProcessResultsFilter struct {
	ResultsFilter bool
}

// SetProcessNamespaceFilter sets the status filter for namespaces
type SetProcessNamespaceFilter struct {
	NamespaceFilter int
}

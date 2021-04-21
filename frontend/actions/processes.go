package actions

import (
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/store/storeutil"
	"go.vocdoni.io/proto/build/go/models"
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

// SetProcessList is the action to set the process list
type SetProcessList struct {
	Processes [config.ListSize]*storeutil.Process
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

// SetCurrentProcessInfo is the action to set the current process info
type SetCurrentProcessInfo struct {
	Process storeutil.Process
}

// SetCurrentProcessStruct is the action to set the current process
type SetCurrentProcessStruct struct {
	Process *storeutil.Process
}

// SetCurrentProcessEnvelopes is the action to set the envelope list for the current process
type SetCurrentProcessEnvelopes struct {
	EnvelopeList *models.EnvelopePackageList
}

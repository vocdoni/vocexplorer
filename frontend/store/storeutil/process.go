package storeutil

import (
	"gitlab.com/vocdoni/vocexplorer/api"
	"gitlab.com/vocdoni/vocexplorer/api/dbtypes"
	"gitlab.com/vocdoni/vocexplorer/config"
)

// Processes stores the current processes information
type Processes struct {
	Count                   int
	ProcessResults          map[string]Process
	ProcessKeys             map[string]*api.Pkeys
	Processes               [config.ListSize]*dbtypes.Process
	Pagination              PageStore
	EnvelopePagination      PageStore
	EnvelopeHeights         map[string]int64
	CurrentProcessResults   Process
	CurrentProcess          *dbtypes.Process
	CurrentProcessEnvelopes [config.ListSize]*dbtypes.Envelope
}

// Process holds info about one vochain process, including votes and results
type Process struct {
	Envelopes     [config.ListSize]*dbtypes.Envelope
	EnvelopeCount int
	ProcessInfo   api.ProcessResults
}

package storeutil

import (
	"go.vocdoni.io/dvote/api"
	indexertypes "go.vocdoni.io/dvote/vochain/scrutinizer/indexertypes"
)

// Processes stores the current processes information
type Processes struct {
	Count              int
	ProcessResults     map[string]ProcessResults
	Processes          map[string]*Process
	ProcessIds         []string
	Pagination         PageStore
	EnvelopePagination PageStore
	CurrentProcess     *Process
}

// Process holds info about one vochain process, including the process and envelope info
type Process struct {
	Envelopes      []*indexertypes.EnvelopeMetadata
	EnvelopeCount  int
	Process        *indexertypes.Process
	ProcessSummary *api.ProcessSummary
	ProcessID      string
}

type ProcessResults struct {
	Results [][]string
	State   string
	Type    string
	Final   bool
}

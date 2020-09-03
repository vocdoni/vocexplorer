package storeutil

import (
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/api"
	"gitlab.com/vocdoni/vocexplorer/types"
)

// Processes stores the current processes information
type Processes struct {
	CurrentProcessID string
	CurrentProcess   Process
	Count            int
	EnvelopesIndex   int
	EnvelopesPage    int
	ProcessIDs       [config.ListSize]string
	ProcessResults   map[string]Process
	ProcessKeys      map[string]*api.Pkeys
	Pagination       PageStore
	EnvelopeHeights  map[string]int64
}

// Process holds info about one vochain process, including votes and results
type Process struct {
	Envelopes     [config.ListSize]*types.Envelope
	EnvelopeCount int
	ProcessType   string
	Results       [][]uint32
	State         string
}

package storeutil

import (
	"gitlab.com/vocdoni/vocexplorer/api"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/proto"
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
	Envelopes     [config.ListSize]*proto.Envelope
	EnvelopeCount int
	ProcessType   string
	Results       [][]uint32
	State         string
}

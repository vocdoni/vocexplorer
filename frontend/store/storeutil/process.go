package storeutil

import (
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/types"
)

// Processes stores the current processes information
type Processes struct {
	CurrentProcessID string
	ProcessCount     int
	ProcessIDs       [config.ListSize]string
	ProcessResults   map[string]Process
	Pagination       PageStore
	EnvelopeHeights  map[string]int64
}

// Process holds info about one vochain process, including votes and results
type Process struct {
	EnvelopeList   [config.ListSize]*types.Envelope
	EnvelopeHeight int
	ProcessType    string
	Results        [][]uint32
	State          string
}

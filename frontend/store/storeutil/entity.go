package storeutil

import (
	"gitlab.com/vocdoni/vocexplorer/config"
)

// Entities stores the current entities information
type Entities struct {
	EntityCount    int
	ProcessHeights map[string]int64
	EntityIDs      [config.ListSize]string
	Entities       [config.ListSize]Entity
	Pagination     PageStore
}

// Entity holds info about one vochain entity
type Entity struct {
	EnvelopeHeights map[string]int64
	ProcessIDs      [config.ListSize]string
	Processes       map[string]Process
	ProcessCount    int
}

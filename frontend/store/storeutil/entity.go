package storeutil

import (
	"gitlab.com/vocdoni/vocexplorer/config"
)

// Entities stores the current entities information
type Entities struct {
	Count           int
	CurrentEntity   Entity
	CurrentEntityID string
	EntityIDs       [config.ListSize]string
	Pagination      PageStore
	ProcessHeights  map[string]int64
}

// Entity holds info about one vochain entity
type Entity struct {
	ProcessCount int
	ProcessIDs   [config.ListSize]string
}

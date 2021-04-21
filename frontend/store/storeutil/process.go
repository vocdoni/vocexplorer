package storeutil

import (
	"gitlab.com/vocdoni/vocexplorer/config"
	"go.vocdoni.io/proto/build/go/models"
)

// Processes stores the current processes information
type Processes struct {
	Count              int
	Processes          [config.ListSize]*Process
	Pagination         PageStore
	EnvelopePagination PageStore
	EnvelopeHeights    map[string]int64
	CurrentProcess     *Process
}

// Process holds info about one vochain process, including votes and results
type Process struct {
	Envelopes     *models.EnvelopePackageList
	EnvelopeCount int
	Process       *models.Process
}

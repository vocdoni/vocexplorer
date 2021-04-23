package storeutil

import (
	"gitlab.com/vocdoni/vocexplorer/config"
	"go.vocdoni.io/proto/build/go/models"
)

// Processes stores the current processes information
type Processes struct {
	Count              int
	ProcessResults     map[string]ProcessResults
	Processes          map[string]*Process
	ProcessIds         [config.ListSize]string
	Pagination         PageStore
	EnvelopePagination PageStore
	EnvelopeHeights    map[string]int64
	CurrentProcess     *Process
}

// Process holds info about one vochain process, including the process and envelope info
type Process struct {
	Envelopes     []*models.EnvelopePackage
	EnvelopeCount int
	Process       *models.Process
}

type ProcessResults struct {
	Results [][]string
	State   string
	Type    string
	Final   bool
}

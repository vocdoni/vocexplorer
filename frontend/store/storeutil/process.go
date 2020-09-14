package storeutil

import (
	"gitlab.com/vocdoni/vocexplorer/api"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/proto"
)

// Processes stores the current processes information
type Processes struct {
	Count                 int
	ProcessResults        map[string]Process
	ProcessKeys           map[string]*api.Pkeys
	Processes             [config.ListSize]*proto.Process
	Pagination            PageStore
	EnvelopePagination    PageStore
	EnvelopeHeights       map[string]int64
	CurrentProcessResults Process
	CurrentProcess        *proto.Process
}

// Process holds info about one vochain process, including votes and results
type Process struct {
	Envelopes     [config.ListSize]*proto.Envelope
	EnvelopeCount int
	ProcessType   string
	Results       [][]uint32
	State         string
}

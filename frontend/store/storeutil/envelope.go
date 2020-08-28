package storeutil

import (
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/types"
)

// Envelopes stores the current envelopes information
type Envelopes struct {
	EnvelopeList  [config.ListSize]*types.Envelope
	EnvelopeCount int
	Pagination    PageStore
}

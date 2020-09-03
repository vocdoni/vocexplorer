package storeutil

import (
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/types"
)

// Envelopes stores the current envelopes information
type Envelopes struct {
	Count                 int
	CurrentEnvelope       *types.Envelope
	CurrentEnvelopeHeight int64
	Envelopes             [config.ListSize]*types.Envelope
	Pagination            PageStore
}

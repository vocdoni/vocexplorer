package storeutil

import (
	"github.com/vocdoni/vocexplorer/api/dbtypes"
	"github.com/vocdoni/vocexplorer/config"
)

// Envelopes stores the current envelopes information
type Envelopes struct {
	Count                 int
	CurrentEnvelope       *dbtypes.Envelope
	CurrentEnvelopeHeight int64
	Envelopes             [config.ListSize]*dbtypes.Envelope
	Pagination            PageStore
}

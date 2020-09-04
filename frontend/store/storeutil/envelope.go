package storeutil

import (
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/proto"
)

// Envelopes stores the current envelopes information
type Envelopes struct {
	Count                 int
	CurrentEnvelope       *proto.Envelope
	CurrentEnvelopeHeight int64
	Envelopes             [config.ListSize]*proto.Envelope
	Pagination            PageStore
}

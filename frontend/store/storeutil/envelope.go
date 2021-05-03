package storeutil

import "go.vocdoni.io/dvote/types"

// Envelopes stores the current envelopes information
type Envelopes struct {
	Count                    int
	CurrentEnvelope          *types.EnvelopePackage
	CurrentEnvelopeNullifier []byte
	Envelopes                []*types.EnvelopePackage
	Pagination               PageStore
}

package storeutil

import sctypes "go.vocdoni.io/dvote/vochain/scrutinizer/types"

// Envelopes stores the current envelopes information
type Envelopes struct {
	Count                    int
	CurrentEnvelope          *sctypes.EnvelopePackage
	CurrentEnvelopeNullifier []byte
	Envelopes                []*sctypes.EnvelopeMetadata
	Pagination               PageStore
}

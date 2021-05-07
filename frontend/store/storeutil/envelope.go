package storeutil

import indexertypes "go.vocdoni.io/dvote/vochain/scrutinizer/indexertypes"

// Envelopes stores the current envelopes information
type Envelopes struct {
	Count                    int
	CurrentEnvelope          *indexertypes.EnvelopePackage
	CurrentEnvelopeNullifier []byte
	Envelopes                []*indexertypes.EnvelopeMetadata
	Pagination               PageStore
}

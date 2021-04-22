package storeutil

import (
	"go.vocdoni.io/proto/build/go/models"
)

// Envelopes stores the current envelopes information
type Envelopes struct {
	Count                    int
	CurrentEnvelope          *models.EnvelopePackage
	CurrentEnvelopeNullifier []byte
	Envelopes                *models.EnvelopePackageList
	Pagination               PageStore
}

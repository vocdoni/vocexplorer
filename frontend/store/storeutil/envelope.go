package storeutil

import (
	"go.vocdoni.io/proto/build/go/models"
)

// Envelopes stores the current envelopes information
type Envelopes struct {
	Count                 int
	CurrentEnvelope       *models.EnvelopePackage
	CurrentEnvelopeHeight int64
	Envelopes             *models.EnvelopePackageList
	Pagination            PageStore
}

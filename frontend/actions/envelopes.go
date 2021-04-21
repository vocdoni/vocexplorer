package actions

import (
	"go.vocdoni.io/proto/build/go/models"
)

// TransactionTabChange is the action to change between tabs in transaction view details
type EnvelopesTabChange struct {
	Tab string
}

// EnvelopesIndexChange is the action to set the pagination index
type EnvelopesIndexChange struct {
	Index int
}

// SetEnvelopeList is the action to set the envelope list
type SetEnvelopeList struct {
	EnvelopeList *models.VoteEnvelopeList
}

// SetEnvelopeCount is the action to set the Envelope count
type SetEnvelopeCount struct {
	Count int
}

// SetCurrentEnvelope is the action to set the current envelope
type SetCurrentEnvelope struct {
	Envelope *models.VoteEnvelope
}

// SetCurrentEnvelopeHeight is the action to set the current envelope height
type SetCurrentEnvelopeHeight struct {
	Height int64
}

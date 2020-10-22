package actions

import (
	"gitlab.com/vocdoni/vocexplorer/api/dbtypes"
	"gitlab.com/vocdoni/vocexplorer/config"
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
	EnvelopeList [config.ListSize]*dbtypes.Envelope
}

// SetEnvelopeCount is the action to set the Envelope count
type SetEnvelopeCount struct {
	Count int
}

// SetCurrentEnvelope is the action to set the current envelope
type SetCurrentEnvelope struct {
	Envelope *dbtypes.Envelope
}

// SetCurrentEnvelopeHeight is the action to set the current envelope height
type SetCurrentEnvelopeHeight struct {
	Height int64
}

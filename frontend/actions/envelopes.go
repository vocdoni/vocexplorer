package actions

import "go.vocdoni.io/dvote/types"

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
	EnvelopeList []*types.EnvelopePackage
}

// SetEnvelopeCount is the action to set the Envelope count
type SetEnvelopeCount struct {
	Count int
}

// SetCurrentEnvelope is the action to set the current envelope
type SetCurrentEnvelope struct {
	Envelope *types.EnvelopePackage
}

// SetCurrentEnvelopeNullifier is the action to set the current envelope nullifier
type SetCurrentEnvelopeNullifier struct {
	Nullifier []byte
}

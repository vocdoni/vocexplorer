package actions

import (
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/proto"
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
	EnvelopeList [config.ListSize]*proto.Envelope
}

// SetEnvelopeCount is the action to set the Envelope count
type SetEnvelopeCount struct {
	Count int
}

// SetCurrentEnvelope is the action to set the current envelope
type SetCurrentEnvelope struct {
	Envelope *proto.Envelope
}

// SetCurrentEnvelopeHeight is the action to set the current envelope height
type SetCurrentEnvelopeHeight struct {
	Height int64
}

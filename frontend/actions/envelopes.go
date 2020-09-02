package actions

import (
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/types"
)

// EnvelopesIndexChange is the action to set the pagination index
type EnvelopesIndexChange struct {
	Index int
}

// SetEnvelopeList is the action to set the envelope list
type SetEnvelopeList struct {
	EnvelopeList [config.ListSize]*types.Envelope
}

// SetEnvelopeCount is the action to set the Envelope count
type SetEnvelopeCount struct {
	Count int
}

// SetCurrentEnvelope is the action to set the current envelope
type SetCurrentEnvelope struct {
	Envelope *types.Envelope
}

// SetCurrentEnvelopeHeight is the action to set the current envelope height
type SetCurrentEnvelopeHeight struct {
	Height int64
}

// DisableEnvelopeUpdate is the action to set the disable update status for envelopes
type DisableEnvelopeUpdate struct {
	Disabled bool
}

// On initialization, register actions
func init() {
	dispatcher.Register(envelopeActions)
}

// envelopeActions is the handler for all envelope-related store actions
func envelopeActions(action interface{}) {
	switch a := action.(type) {
	case *EnvelopesIndexChange:
		store.Envelopes.Pagination.Index = a.Index

	case *SetEnvelopeList:
		store.Envelopes.Envelopes = a.EnvelopeList

	case *SetEnvelopeCount:
		store.Envelopes.Count = a.Count

	case *SetCurrentEnvelope:
		store.Envelopes.CurrentEnvelope = a.Envelope

	case *SetCurrentEnvelopeHeight:
		store.Envelopes.CurrentEnvelopeHeight = a.Height

	case *DisableEnvelopeUpdate:
		store.Envelopes.Pagination.DisableUpdate = a.Disabled

	default:
		return // don't fire listeners
	}

	store.Listeners.Fire()
}

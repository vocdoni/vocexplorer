package actions

import (
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/types"
)

// SetEnvelopeList is the action to set the envelope list
type SetEnvelopeList struct {
	EnvelopeList [config.ListSize]*types.Envelope
}

// SetEnvelopeCount is the action to set the Envelope count
type SetEnvelopeCount struct {
	EnvelopeCount int
}

// On initialization, register actions
func init() {
	dispatcher.Register(envelopeActions)
}

// envelopeActions is the handler for all envelope-related store actions
func envelopeActions(action interface{}) {
	switch a := action.(type) {
	case *SetEnvelopeList:
		store.Envelopes.EnvelopeList = a.EnvelopeList

	case *SetEnvelopeCount:
		store.Envelopes.EnvelopeCount = a.EnvelopeCount

	default:
		return // don't fire listeners
	}

	store.Listeners.Fire()
}

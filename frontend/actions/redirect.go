package actions

import (
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
)

// SignalRedirect is the action to signal a page redirect
type SignalRedirect struct {
}

// On initialization, register actions
func init() {
	dispatcher.Register(redirectActions)
}

// redirectActions is the handler for all redirect-related actions
func redirectActions(action interface{}) {
	switch action.(type) {
	case *SignalRedirect:
		store.RedirectChan <- struct{}{}

	default:
		return // don't fire listeners
	}

	store.Listeners.Fire()
}

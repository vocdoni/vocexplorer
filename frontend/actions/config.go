package actions

import (
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
)

// StoreConfig is the action to store a config object
type StoreConfig struct {
	Config config.Cfg
}

// On initialization, register actions
func init() {
	dispatcher.Register(configActions)
}

// configActions is the handler for all config-related actions
func configActions(action interface{}) {
	switch a := action.(type) {
	case *StoreConfig:
		store.Config = a.Config

	default:
		return // don't fire listeners
	}

	store.Listeners.Fire()
}

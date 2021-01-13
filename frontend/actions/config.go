package actions

import (
	"github.com/vocdoni/vocexplorer/config"
)

// StoreConfig is the action to store a config object
type StoreConfig struct {
	Config config.Cfg
}

package actions

import (
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
)

// EntitiesTabChange is the action to change entities tabs
type EntitiesTabChange struct {
	Tab string
}

// SetEntityIDs is the action to set the entity list
type SetEntityIDs struct {
	EntityIDs [config.ListSize]string
}

// SetEntityCount is the action to set the entity count
type SetEntityCount struct {
	EntityCount int
}

// SetProcessHeights is the action to set the entity count
type SetProcessHeights struct {
	ProcessHeights map[string]int64
}

// On initialization, register actions
func init() {
	dispatcher.Register(envelopeActions)
}

// entityActions is the handler for all entity-related store actions
func entityActions(action interface{}) {
	switch a := action.(type) {
	case *SetEntityIDs:
		store.Entities.EntityIDs = a.EntityIDs

	case *EntitiesTabChange:
		store.Entities.Pagination.Tab = a.Tab

	case *SetEntityCount:
		store.Entities.EntityCount = a.EntityCount

	case *SetProcessHeights:
		store.Entities.ProcessHeights = a.ProcessHeights

	default:
		return // don't fire listeners
	}

	store.Listeners.Fire()
}

package actions

import (
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
)

// EntitiesIndexChange is the action to set the pagination index
type EntitiesIndexChange struct {
	Index int
}

// EntityProcessesIndexChange is the action to set the pagination index for the current entity's process list
type EntityProcessesIndexChange struct {
	Index int
}

// EntityProcessesPageChange is the action to set the pagination page for the current entity's process list
type EntityProcessesPageChange struct {
	Index int
}

// EntitiesTabChange is the action to change entities tabs
type EntitiesTabChange struct {
	Tab string
}

// SetEntityIDs is the action to set the entity list
type SetEntityIDs struct {
	EntityIDs [config.ListSize]string
}

// SetCurrentEntityID is the action to set the current active entity ID
type SetCurrentEntityID struct {
	EntityID string
}

// SetEntityCount is the action to set the entity count
type SetEntityCount struct {
	Count int
}

// SetEntityProcessCount is the action to set the current entity's process count
type SetEntityProcessCount struct {
	Count int
}

// SetProcessHeights is the action to set the entity count
type SetProcessHeights struct {
	ProcessHeights map[string]int64
}

// SetEntityProcessList is the action to set the current entity's process list
type SetEntityProcessList struct {
	ProcessList [config.ListSize]string
}

// DisableEntityUpdate is the action to set the disable update status for entities
type DisableEntityUpdate struct {
	Disabled bool
}

// On initialization, register actions
func init() {
	dispatcher.Register(entityActions)
}

// entityActions is the handler for all entity-related store actions
func entityActions(action interface{}) {
	switch a := action.(type) {
	case *EntitiesIndexChange:
		store.Entities.Pagination.Index = a.Index

	case *EntityProcessesIndexChange:
		store.Entities.ProcessesIndex = a.Index

	case *EntityProcessesPageChange:
		store.Entities.ProcessesPage = a.Index

	case *SetEntityIDs:
		store.Entities.EntityIDs = a.EntityIDs

	case *SetCurrentEntityID:
		store.Entities.CurrentEntityID = a.EntityID

	case *EntitiesTabChange:
		store.Entities.Pagination.Tab = a.Tab

	case *SetEntityCount:
		store.Entities.Count = a.Count

	case *SetProcessHeights:
		store.Entities.ProcessHeights = a.ProcessHeights

	case *SetEntityProcessList:
		store.Entities.CurrentEntity.ProcessIDs = a.ProcessList

	case *DisableEntityUpdate:
		store.Entities.Pagination.DisableUpdate = a.Disabled

	case *SetEntityProcessCount:
		store.Entities.CurrentEntity.ProcessCount = a.Count

	default:
		return // don't fire listeners
	}

	store.Listeners.Fire()
}

package actions

import (
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/store/storeutil"
)

// EntitiesIndexChange is the action to set the pagination index
type EntitiesIndexChange struct {
	Index int
}

// EntityProcessesIndexChange is the action to set the pagination index
type EntityProcessesIndexChange struct {
	Index int
}

// EntityTabChange is the action to change entities tabs
type EntityTabChange struct {
	Tab string
}

// SetEntityIDs is the action to set the entity list
type SetEntityIDs struct {
	EntityIDs []string
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

// SetEntityProcessIds is the action to set the current entity's process ids
type SetEntityProcessIds struct {
	ProcessList []string
}

// SetEntityProcessIds is the action to set the current entity's process ids
type SetEntityProcessList struct {
	ProcessList [config.ListSize]*storeutil.Process
}

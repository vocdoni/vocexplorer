package actions

import (
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/types"
)

// SetValidatorList is the action to set the validator list
type SetValidatorList struct {
	List [config.ListSize]*types.Validator
}

// SetValidatorCount is the action to set the Validator count
type SetValidatorCount struct {
	Count int
}

// DisableValidatorUpdate is the action to set the disable update status for validators
type DisableValidatorUpdate struct {
	Disabled bool
}

// On initialization, register actions
func init() {
	dispatcher.Register(validatorActions)
}

// validatorActions is the handler for all validator-related store actions
func validatorActions(action interface{}) {
	switch a := action.(type) {
	case *SetValidatorList:
		store.Validators.Validators = a.List

	case *SetValidatorCount:
		store.Validators.Count = a.Count

	case *DisableValidatorUpdate:
		store.Validators.Pagination.DisableUpdate = a.Disabled

	default:
		return // don't fire listeners
	}

	store.Listeners.Fire()
}

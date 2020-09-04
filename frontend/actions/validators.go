package actions

import (
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/proto"
)

// ValidatorsIndexChange is the action to set the pagination index
type ValidatorsIndexChange struct {
	Index int
}

// SetValidatorList is the action to set the validator list
type SetValidatorList struct {
	List [config.ListSize]*proto.Validator
}

// SetValidatorCount is the action to set the Validator count
type SetValidatorCount struct {
	Count int
}

// DisableValidatorUpdate is the action to set the disable update status for validators
type DisableValidatorUpdate struct {
	Disabled bool
}

// SetCurrentValidator is the action to set the currently displayed validator
type SetCurrentValidator struct {
	Validator *proto.Validator
}

// SetCurrentValidatorID is the action to set the currently displayed validator ID
type SetCurrentValidatorID struct {
	ID string
}

// SetCurrentValidatorBlockCount is the action to set the currently displayed validator's block count
type SetCurrentValidatorBlockCount struct {
	Count int
}

// SetCurrentValidatorBlockList is the action to set the list of blocks belonging to the current validator
type SetCurrentValidatorBlockList struct {
	BlockList [config.ListSize]*proto.StoreBlock
}

// On initialization, register actions
func init() {
	dispatcher.Register(validatorActions)
}

// validatorActions is the handler for all validator-related store actions
func validatorActions(action interface{}) {
	switch a := action.(type) {
	case *ValidatorsIndexChange:
		store.Validators.Pagination.Index = a.Index

	case *SetValidatorList:
		store.Validators.Validators = a.List

	case *SetValidatorCount:
		store.Validators.Count = a.Count

	case *DisableValidatorUpdate:
		store.Validators.Pagination.DisableUpdate = a.Disabled

	case *SetCurrentValidator:
		store.Validators.CurrentValidator = a.Validator

	case *SetCurrentValidatorID:
		store.Validators.CurrentValidatorID = a.ID

	case *SetCurrentValidatorBlockCount:
		store.Validators.CurrentBlockCount = a.Count

	case *SetCurrentValidatorBlockList:
		store.Validators.CurrentBlockList = a.BlockList

	default:
		return // don't fire listeners
	}

	store.Listeners.Fire()
}

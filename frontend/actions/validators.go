package actions

import (
	"gitlab.com/vocdoni/vocexplorer/config"
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

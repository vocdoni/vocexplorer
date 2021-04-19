package actions

import (
	"gitlab.com/vocdoni/vocexplorer/api/dbtypes"
	"gitlab.com/vocdoni/vocexplorer/config"
)

// ValidatorsIndexChange is the action to set the pagination index
type ValidatorsIndexChange struct {
	Index int
}

// ValidatorBlocksIndexChange is the action to set the pagination index
type ValidatorBlocksIndexChange struct {
	Index int
}

// SetValidatorList is the action to set the validator list
type SetValidatorList struct {
	List [config.ListSize]*dbtypes.Validator
}

// SetValidatorCount is the action to set the Validator count
type SetValidatorCount struct {
	Count int
}

// SetCurrentValidator is the action to set the currently displayed validator
type SetCurrentValidator struct {
	Validator *dbtypes.Validator
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
	BlockList [config.ListSize]*dbtypes.StoreBlock
}

// SetValidatorBlockHeightMap is the action to set the map of block heights associated with each validator
type SetValidatorBlockHeightMap struct {
	HeightMap map[string]int64
}

package storeutil

import (
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/types"
)

// Validators stores all data about blockchain validators
type Validators struct {
	Count              int
	CurrentValidatorID string
	CurrentValidator   *types.Validator
	CurrentBlockCount  int
	CurrentBlockList   [config.ListSize]*types.StoreBlock
	Pagination         PageStore
	Validators         [config.ListSize]*types.Validator
}
package storeutil

import (
	"github.com/vocdoni/vocexplorer/api/dbtypes"
	"github.com/vocdoni/vocexplorer/config"
)

// Validators stores all data about blockchain validators
type Validators struct {
	BlockHeights       map[string]int64
	Count              int
	CurrentBlockCount  int
	CurrentBlockList   [config.ListSize]*dbtypes.StoreBlock
	CurrentValidator   *dbtypes.Validator
	CurrentValidatorID string
	Pagination         PageStore
	BlockPagination    PageStore
	Validators         [config.ListSize]*dbtypes.Validator
}

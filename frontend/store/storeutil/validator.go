package storeutil

import (
	"go.vocdoni.io/proto/build/go/models"
)

// Validators stores all data about blockchain validators
type Validators struct {
	BlockHeights       map[string]int64
	Count              int
	CurrentBlockCount  int
	CurrentBlockList   *models.TendermintHeaderList
	CurrentValidator   *models.Validator
	CurrentValidatorID string
	Pagination         PageStore
	BlockPagination    PageStore
	Validators         *models.ValidatorList
}

package storeutil

import (
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/proto"
)

// Validators stores all data about blockchain validators
type Validators struct {
	BlockHeights       map[string]int64
	Count              int
	CurrentBlockCount  int
	CurrentBlockList   [config.ListSize]*proto.StoreBlock
	CurrentValidator   *proto.Validator
	CurrentValidatorID string
	Pagination         PageStore
	Validators         [config.ListSize]*proto.Validator
}

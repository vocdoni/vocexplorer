package storeutil

import (
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/proto"
)

// Validators stores all data about blockchain validators
type Validators struct {
	Count              int
	CurrentValidatorID string
	CurrentValidator   *proto.Validator
	CurrentBlockCount  int
	CurrentBlockList   [config.ListSize]*proto.StoreBlock
	Pagination         PageStore
	Validators         [config.ListSize]*proto.Validator
}

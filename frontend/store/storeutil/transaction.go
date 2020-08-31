package storeutil

import (
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/types"
)

// Transactions stores all data about blockchain transactions
type Transactions struct {
	Count        int
	Pagination   PageStore
	Transactions [config.ListSize]*types.SendTx
}

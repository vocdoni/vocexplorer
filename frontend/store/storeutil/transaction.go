package storeutil

import (
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/types"
)

// Transactions stores all data about blockchain transactions
type Transactions struct {
	Pagination PageStore
	TxCount    int
	TxList     [config.ListSize]*types.SendTx
}

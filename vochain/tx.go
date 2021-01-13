package vochain

import (
	"fmt"

	"github.com/vocdoni/vocexplorer/proto"
)

// GetTransactions retrieves all transactions from a single block
func (vs *VochainService) GetTransactions(blockHeight int64) ([]*proto.Transaction, error) {
	block := vs.app.Node.BlockStore().LoadBlock(blockHeight)
	if block == nil {
		return nil, fmt.Errorf("block %d does not exist", blockHeight)
	}
	var txList []*proto.Transaction
	for i, tmTx := range block.Txs {
		txList = append(txList, &proto.Transaction{
			Hash:   tmTx.Hash(),
			Height: blockHeight,
			Index:  uint32(i),
			Tx:     tmTx,
		})
	}
	return txList, nil
}

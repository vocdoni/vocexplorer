package vochain

import (
	"fmt"
	"reflect"

	"gitlab.com/vocdoni/vocexplorer/api/tmtypes"
	"gitlab.com/vocdoni/vocexplorer/proto"
)

// GetTransactions retrieves all transactions from a single block
func (vs *VochainService) GetTransactions(blockHeight int64) ([]*proto.Transaction, error) {
	block := vs.app.Node.BlockStore().LoadBlock(blockHeight)
	if block == nil {
		return nil, fmt.Errorf("block %d does not exist", blockHeight)
	}
	var txList []*proto.Transaction
	for i, tmTx := range block.Txs {
		var tx tmtypes.Tx = reflect.ValueOf(tmTx).Bytes()
		txList = append(txList, &proto.Transaction{
			Hash:   tx.Hash(),
			Height: blockHeight,
			Index:  uint32(i),
			Tx:     tx,
		})
	}
	return txList, nil
}

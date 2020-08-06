package types

import (
	"time"

	abci "github.com/tendermint/tendermint/abci/types"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	tmtypes "github.com/tendermint/tendermint/types"
)

//StoreBlock stores the parts of a block relevant to our database
type StoreBlock struct {
	NumTxs int              `json:",omitempty"`
	Hash   tmbytes.HexBytes `json:",omitempty"`
	Height int64            `json:",omitempty"`
	Time   time.Time        `json:",omitempty"`
}

//StoreTx stores the parts of a tx relevant to our database
type StoreTx struct {
	Height   int64                  `json:",omitempty"`
	TxHeight int64                  `json:",omitempty"`
	Tx       tmtypes.Tx             `json:",omitempty"`
	TxResult abci.ResponseDeliverTx `json:",omitempty"`
	Index    uint32                 `json:",omitempty"`
}

//SendTx stores the parts of a tx relevant to our database
type SendTx struct {
	Hash  tmbytes.HexBytes `json:",omitempty"`
	Store StoreTx          `json:",omitempty"`
}

//IsEmpty returns true if block is empty
func (s StoreBlock) IsEmpty() bool {
	if len(s.Hash) == 0 && s.Height == 0 && s.NumTxs == 0 {
		return true
	}
	return false
}

//IsEmpty returns true if tx is empty
func (s SendTx) IsEmpty() bool {
	if len(s.Hash) == 0 && s.Store.TxHeight == 0 && s.Store.Height == 0 && s.Store.Index == 0 {
		return true
	}
	return false
}

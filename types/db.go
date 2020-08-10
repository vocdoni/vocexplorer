package types

import (
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
)

//SendTx stores the parts of a tx relevant to our database
type SendTx struct {
	Hash  tmbytes.HexBytes `json:",omitempty"`
	Store *StoreTx         `json:",omitempty"`
}

//BlockIsEmpty returns true if block is empty
func BlockIsEmpty(s *StoreBlock) bool {
	if len(s.GetHash()) == 0 && s.GetHeight() == 0 && s.GetNumTxs() == 0 {
		return true
	}
	return false
}

//IsEmpty returns true if tx is empty
func (s SendTx) IsEmpty() bool {
	if len(s.Hash) == 0 && s.Store.GetTxHeight() == 0 && s.Store.GetHeight() == 0 && s.Store.GetIndex() == 0 {
		return true
	}
	return false
}

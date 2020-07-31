package types

import (
	"time"

	tmbytes "github.com/tendermint/tendermint/libs/bytes"
)

//StoreBlock stores the parts of a block relevant to our database
type StoreBlock struct {
	NumTxs int
	Hash   tmbytes.HexBytes
	Height int64
	Time   time.Time
}

//IsEmpty returns true if block is empty
func (s StoreBlock) IsEmpty() bool {
	if len(s.Hash) == 0 && s.Height == 0 && s.NumTxs == 0 {
		return true
	}
	return false
}

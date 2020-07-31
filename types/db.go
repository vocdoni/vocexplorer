package types

import (
	"time"

	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	ttypes "github.com/tendermint/tendermint/types"
)

//StoreBlock stores the parts of a block relevant to our database
type StoreBlock struct {
	Data   ttypes.Data
	Hash   tmbytes.HexBytes
	Height int64
	Time   time.Time
}

//IsEmpty returns true if block is empty
func (s StoreBlock) IsEmpty() bool {
	if len(s.Hash) == 0 {
		return true
	}
	return false
}

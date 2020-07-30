package types

import (
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	ttypes "github.com/tendermint/tendermint/types"
)

//StoreBlock stores the parts of a block relevant to our database
type StoreBlock struct {
	Hash   tmbytes.HexBytes
	Header ttypes.Header
}

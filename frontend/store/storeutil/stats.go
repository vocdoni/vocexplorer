package storeutil

import (
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

// Stats stores the blockchain statistics
type Stats struct {
	APIList        []string
	AppVersion     int
	BlockTime      *[5]int32
	BlockTimeStamp int32
	Health         int32
	Height         int64
	MaxBlockSize   int
	Genesis        *tmtypes.GenesisDoc
	ResultStatus   *coretypes.ResultStatus
}

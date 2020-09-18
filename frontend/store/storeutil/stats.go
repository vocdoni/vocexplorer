package storeutil

import (
	"time"

	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

// Stats stores the blockchain statistics
type Stats struct {
	APIList           []string
	AppVersion        int
	AvgTxsPerBlock    float64
	AvgTxsPerMinute   float64
	BlockTime         *[5]int32
	BlockTimeStamp    int32
	Genesis           *tmtypes.GenesisDoc
	Health            int32
	Height            int64
	MaxBlockSize      int
	MaxTxsBlockHash   string
	MaxTxsBlockHeight int64
	MaxTxsMinute      time.Time
	MaxTxsPerBlock    int64
	MaxTxsPerMinute   int64
	ResultStatus      *coretypes.ResultStatus
}

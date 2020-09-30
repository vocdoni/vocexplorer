package actions

import (
	"time"

	"gitlab.com/vocdoni/vocexplorer/api/tmtypes"
)

//SetResultStatus is the action to set the blockchain stats result status
type SetResultStatus struct {
	Status *tmtypes.ResultStatus
}

//SetGenesis is the action to set the blockchain genesis block
type SetGenesis struct {
	Genesis *tmtypes.GenesisDoc
}

//SetGatewayInfo is the action to set the gateway statistic info
type SetGatewayInfo struct {
	APIList []string
	Health  int32
}

//SetBlockStatus is the action to set the latest block status
type SetBlockStatus struct {
	BlockTime      *[5]int32
	BlockTimeStamp int32
	Height         int64
}

//SetTransactionStats sets the transaction stats (avg, max, etc)
type SetTransactionStats struct {
	AvgTxsPerBlock    float64
	AvgTxsPerMinute   float64
	MaxTxsBlockHash   string
	MaxTxsBlockHeight int64
	MaxTxsMinute      time.Time
	MaxTxsPerBlock    int64
	MaxTxsPerMinute   int64
}

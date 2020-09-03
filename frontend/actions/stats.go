package actions

import (
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
)

//SetResultStatus is the action to set the blockchain stats result status
type SetResultStatus struct {
	Status *coretypes.ResultStatus
}

//SetGenesis is the action to set the blockchain genesis block
type SetGenesis struct {
	Genesis *tmtypes.GenesisDoc
}

//SetGatewayInfo is the action to set the gateway statistic info
type SetGatewayInfo struct {
	APIList []string
	Ok      bool
	Health  int32
}

//SetBlockStatus is the action to set the latest block status
type SetBlockStatus struct {
	BlockTime      *[5]int32
	BlockTimeStamp int32
	Height         int64
}

// On initialization, register actions
func init() {
	dispatcher.Register(statsActions)
}

// statsActions is the handler for all stats-related store actions
func statsActions(action interface{}) {
	switch a := action.(type) {
	case *SetResultStatus:
		store.Stats.ResultStatus = a.Status

	case *SetGenesis:
		store.Stats.Genesis = a.Genesis

	case *SetGatewayInfo:
		store.Stats.APIList = a.APIList
		store.Stats.Ok = a.Ok
		store.Stats.Health = a.Health

	case *SetBlockStatus:
		store.Stats.BlockTime = a.BlockTime
		store.Stats.BlockTimeStamp = a.BlockTimeStamp
		store.Stats.Height = a.Height

	default:
		return // don't fire listeners
	}

	store.Listeners.Fire()
}

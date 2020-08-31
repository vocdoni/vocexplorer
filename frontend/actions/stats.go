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

	default:
		return // don't fire listeners
	}

	store.Listeners.Fire()
}

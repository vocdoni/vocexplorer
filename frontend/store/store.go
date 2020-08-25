package store

import (
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store/storeutil"
)

var (
	BlockTabActive = "transactions"
	// CurrentBlockHeight stores the latest known block height
	CurrentBlockHeight int64

	// Listeners is the listeners that will be invoked when the store changes.
	Listeners = storeutil.NewListenerRegistry()

	Processes struct {
		CurrentPage   int
		PagChannel    chan int
		DisableUpdate bool
	}
)

func init() {
	dispatcher.Register(onAction)
}

func onAction(action interface{}) {
	switch a := action.(type) {
	case *actions.BlocksTabChange:
		BlockTabActive = a.Tab

	case *actions.BlocksHeightUpdate:
		CurrentBlockHeight = a.Height

	default:
		return // don't fire listeners
	}

	Listeners.Fire()
}

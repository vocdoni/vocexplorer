package store

import (
	"fmt"

	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store/storeutil"
)

// type Store struct {
// 	BlockTabActive string
// 	Listeners      *storeutil.ListenerRegistry
// }

// var Instance = Store{
// 	Listeners: storeutil.NewListenerRegistry(),
// }

var (
	BlockTabActive = "transactions"
	// CurrentBlockHeight stores the latest known block height
	CurrentBlockHeight int64

	// Listeners is the listeners that will be invoked when the store changes.
	Listeners = storeutil.NewListenerRegistry()
)

func init() {
	dispatcher.Register(onAction)
}

func onAction(action interface{}) {
	switch a := action.(type) {
	case *actions.BlocksTabChange:
		BlockTabActive = a.Tab

	case *actions.BlocksHeightUpdate:
		fmt.Println(a.Height)
		CurrentBlockHeight = a.Height

	default:
		return // don't fire listeners
	}

	Listeners.Fire()
}

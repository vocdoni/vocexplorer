package store

import (
	"github.com/tendermint/tendermint/rpc/client/http"
	"gitlab.com/vocdoni/vocexplorer/client"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store/storeutil"
)

var (
	// BlockTabActive stores the current active block tab
	BlockTabActive string
	// CurrentBlockHeight stores the latest known block height
	CurrentBlockHeight int64

	// Listeners is the listeners that will be invoked when the store changes.
	Listeners = storeutil.NewListenerRegistry()

	// RedirectChan is the channel which signals a page redirect
	RedirectChan chan struct{}

	// GatewayClient is the global gateway client
	GatewayClient *client.Client
	// TendermintClient is the global tendermint client
	TendermintClient *http.HTTP

	// Processes stores the current processes information
	Processes struct {
		Tab           string
		PagChannel    chan int
		CurrentPage   int
		DisableUpdate bool
	}

	// Entities stores the current entities information
	Entities struct {
		Tab           string
		CurrentPage   int
		PagChannel    chan int
		DisableUpdate bool
	}
)

func init() {
	BlockTabActive = "transactions"
	Processes.Tab = "results"
	Entities.Tab = "processes"
	RedirectChan = make(chan struct{}, 100)

	dispatcher.Register(onAction)
}

func onAction(action interface{}) {
	switch a := action.(type) {
	case *actions.BlocksTabChange:
		BlockTabActive = a.Tab

	case *actions.BlocksHeightUpdate:
		CurrentBlockHeight = a.Height

	case *actions.ProcessesTabChange:
		Processes.Tab = a.Tab

	case *actions.SignalRedirect:
		RedirectChan <- struct{}{}

	default:
		return // don't fire listeners
	}

	Listeners.Fire()
}

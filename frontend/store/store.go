package store

import (
	"github.com/tendermint/tendermint/rpc/client/http"
	"gitlab.com/vocdoni/vocexplorer/client"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store/storeutil"
	"gitlab.com/vocdoni/vocexplorer/rpc/rpcinit"
)

var (
	// Config stores the application configuration
	Config config.Cfg
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

	// Entities holds all entity information
	Entities storeutil.Entities
	// Processes holds all entity information
	Processes storeutil.Processes
	// Envelopes holds all entity information
	Envelopes storeutil.Envelopes
	// Stats holds all blockchain stats
	Stats storeutil.Stats
	// Blocks holds all blockchain Blocks
	Blocks storeutil.Blocks
	// Transactions holds all blockchain transactions
	Transactions storeutil.Transactions
	// Validators holds all blockchain Validators
	Validators storeutil.Validators
)

func init() {
	BlockTabActive = "transactions"
	Processes.Pagination.Tab = "results"
	Entities.Pagination.Tab = "processes"
	RedirectChan = make(chan struct{}, 100)
	Entities.Pagination.PagChannel = make(chan int, 50)
	Processes.Pagination.PagChannel = make(chan int, 50)
	Envelopes.Pagination.PagChannel = make(chan int, 50)
	Blocks.Pagination.PagChannel = make(chan int, 50)
	Transactions.Pagination.PagChannel = make(chan int, 50)
	Validators.Pagination.PagChannel = make(chan int, 50)

	dispatcher.Register(onAction)
}

func onAction(action interface{}) {
	switch a := action.(type) {
	case *actions.TendermintClientInit:
		TendermintClient = rpcinit.StartClient(Config.TendermintHost)

	case *actions.GatewayClientInit:
		GatewayClient, _ = client.InitGateway(Config.GatewayHost)

	case *actions.BlocksTabChange:
		BlockTabActive = a.Tab

	case *actions.BlocksHeightUpdate:
		CurrentBlockHeight = a.Height

	case *actions.ProcessesTabChange:
		Processes.Pagination.Tab = a.Tab

	case *actions.SignalRedirect:
		RedirectChan <- struct{}{}

	case *actions.StoreConfig:
		Config = a.Config

	default:
		// return // don't fire listeners
	}

	Listeners.Fire()
}

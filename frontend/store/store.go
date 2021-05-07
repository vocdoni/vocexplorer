package store

import (
	"gitlab.com/vocdoni/vocexplorer/client"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/store/storeutil"
	"go.vocdoni.io/dvote/api"
)

var (
	// Config stores the application configuration
	Config config.Cfg

	Client *client.Client

	// Listeners is the listeners that will be invoked when the store changes.
	Listeners = storeutil.NewListenerRegistry()

	// RedirectChan is the channel which signals a page redirect
	RedirectChan chan struct{}
	// CurrentPage holds the current page title, used in signaling updater goroutines to exit
	CurrentPage string

	// ServerConnected is true if the webserver is connected
	ServerConnected bool

	// Entities holds all entity information
	Entities storeutil.Entities
	// Processes holds all entity information
	Processes storeutil.Processes
	// Envelopes holds all entity information
	Envelopes storeutil.Envelopes
	// Stats holds all blockchain stats
	Stats *api.VochainStats
	// Blocks holds all blockchain Blocks
	Blocks storeutil.Blocks
	// Transactions holds all blockchain transactions
	Transactions storeutil.Transactions
	// Validators holds all blockchain Validators
	Validators storeutil.Validators
	// ProcessDomain is the link for process profiles
	ProcessDomain string
	// EntityDomain is the link for entity profiles
	EntityDomain string
)

func init() {
	Blocks.Pagination.Tab = "transactions"
	Processes.Pagination.Tab = "results"
	Entities.Pagination.Tab = "processes"
	Transactions.Pagination.Tab = "contents"
	Envelopes.Pagination.Tab = "contents"

	RedirectChan = make(chan struct{}, 50)
	Entities.Pagination.PagChannel = make(chan int, 50)
	Entities.ProcessPagination.PagChannel = make(chan int, 50)
	Processes.Pagination.PagChannel = make(chan int, 50)
	Processes.EnvelopePagination.PagChannel = make(chan int, 50)
	Envelopes.Pagination.PagChannel = make(chan int, 50)
	Blocks.Pagination.PagChannel = make(chan int, 50)
	Blocks.TransactionPagination.PagChannel = make(chan int, 50)
	Transactions.Pagination.PagChannel = make(chan int, 50)
	Validators.Pagination.PagChannel = make(chan int, 50)

	Entities.Pagination.SearchChannel = make(chan string, 50)
	Entities.ProcessPagination.SearchChannel = make(chan string, 50)
	Processes.Pagination.SearchChannel = make(chan string, 50)
	Processes.EnvelopePagination.SearchChannel = make(chan string, 50)
	Envelopes.Pagination.SearchChannel = make(chan string, 50)
	Blocks.Pagination.SearchChannel = make(chan string, 50)
	Blocks.TransactionPagination.SearchChannel = make(chan string, 50)
	Transactions.Pagination.SearchChannel = make(chan string, 50)
	Validators.Pagination.SearchChannel = make(chan string, 50)

	Entities.Pagination.Search = false
	Entities.ProcessPagination.Search = false
	Processes.Pagination.Search = false
	Processes.EnvelopePagination.Search = false
	Envelopes.Pagination.Search = false
	Blocks.Pagination.Search = false
	Blocks.TransactionPagination.Search = false
	Transactions.Pagination.Search = false
	Validators.Pagination.Search = false

	Processes.ProcessResults = make(map[string]storeutil.ProcessResults)
	Processes.Processes = make(map[string]*storeutil.Process)
	Entities.ProcessHeights = make(map[string]int64)

	ServerConnected = true

	ProcessDomain = "https://app.vocdoni.net/processes/#/0x"
	EntityDomain = "https://vocdoni.link/entities/0x"

	Reduce()
}

package router

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/db"
)

// RegisterRoutes takes a mux and registers all the routes callbacks within this package
func RegisterRoutes(m *mux.Router, cfg *config.Cfg, d *db.ExplorerDB) {

	// Page Routes
	m.HandleFunc("/", indexHandler)
	m.HandleFunc("/processes", indexHandler)
	m.HandleFunc("/process/{id}", indexHandler)
	m.HandleFunc("/entities", indexHandler)
	m.HandleFunc("/entity/{id}", indexHandler)
	m.HandleFunc("/envelopes", indexHandler)
	m.HandleFunc("/envelope/{id}", indexHandler)
	m.HandleFunc("/blocks", indexHandler)
	m.HandleFunc("/block/{id}", indexHandler)
	m.HandleFunc("/transactions", indexHandler)
	m.HandleFunc("/transaction/{id}", indexHandler)
	m.HandleFunc("/stats", indexHandler)
	m.HandleFunc("/validators", indexHandler)
	m.HandleFunc("/validator/{id}", indexHandler)

	// API Routes
	m.HandleFunc("/ping", PingHandler())
	m.HandleFunc("/config", configHandler(cfg))

	// Blocks
	m.HandleFunc("/api/block/", GetBlockHandler(d))
	m.HandleFunc("/api/blockheader/", GetBlockHeaderHandler(d))
	m.HandleFunc("/api/listblocks/", ListBlocksHandler(d))
	m.HandleFunc("/api/listblocksvalidator/", ListBlocksByValidatorHandler(d))
	m.HandleFunc("/api/numblocksvalidator/", NumBlocksByValidatorHandler(d))
	m.HandleFunc("/api/blocksearch/", SearchBlocksHandler(d))
	m.HandleFunc("/api/validatorblocksearch/", SearchBlocksByValidatorHandler(d))

	// Transactions
	m.HandleFunc("/api/txbyheight/", GetTxByHeightHandler(d))
	m.HandleFunc("/api/txbyhash/", GetTxByHashHandler(d))
	m.HandleFunc("/api/listtxs/", ListTxsHandler(d))
	m.HandleFunc("/api/txhash/", TxHeightFromHashHandler(d))
	m.HandleFunc("/api/transactionsearch/", SearchTransactionsHandler(d))

	// Processes
	m.HandleFunc("/api/process/", GetProcessHandler(d))
	m.HandleFunc("/api/processresults/", GetProcessResultsHandler(d))
	m.HandleFunc("/api/processkeys/", GetProcessKeysHandler(d))
	m.HandleFunc("/api/listprocesses/", ListProcessesHandler(d))
	m.HandleFunc("/api/listprocessesbyentity/", ListProcessesByEntityHandler(d))
	m.HandleFunc("/api/procenvheight/", EnvelopeHeightByProcessHandler(d))
	m.HandleFunc("/api/processsearch/", SearchProcessesHandler(d))

	// Envelopes
	m.HandleFunc("/api/envelope/", GetEnvelopeHandler(d))
	m.HandleFunc("/api/listenvelopes/", ListEnvelopesHandler(d))
	m.HandleFunc("/api/listenvelopesprocess/", ListEnvelopesByProcessHandler(d))
	m.HandleFunc("/api/envelopenullifier/", EnvelopeHeightFromNullifierHandler(d))
	m.HandleFunc("/api/envelopesearch/", SearchEnvelopesHandler(d))

	// Entities
	m.HandleFunc("/api/entityprocheight/", ProcessHeightByEntityHandler(d))
	m.HandleFunc("/api/entitysearch/", SearchEntitiesHandler(d))
	m.HandleFunc("/api/listentities/", ListEntitiesHandler(d))

	// Validators
	m.HandleFunc("/api/validator/", GetValidatorHandler(d))
	m.HandleFunc("/api/listvalidators/", ListValidatorsHandler(d))
	m.HandleFunc("/api/validatorsearch/", SearchValidatorsHandler(d))

	// Other
	m.HandleFunc("/api/height/", HeightHandler(d))
	m.HandleFunc("/api/heightmap/", HeightMapHandler(d))
	m.HandleFunc("/api/stats", StatsHandler(d, cfg))

	m.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	m.NotFoundHandler = http.Handler(http.NotFoundHandler())
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/index.html")
}

func configHandler(cfg *config.Cfg) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(w).Encode(cfg); err != nil {
			panic(err)
		}
	}
}

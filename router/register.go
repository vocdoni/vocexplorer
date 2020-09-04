package router

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	dvotedb "gitlab.com/vocdoni/go-dvote/db"
	"gitlab.com/vocdoni/vocexplorer/config"
)

// RegisterRoutes takes a mux and registers all the routes callbacks within this package
func RegisterRoutes(m *mux.Router, cfg *config.Cfg, d *dvotedb.BadgerDB) {

	// Page Routes
	m.HandleFunc("/", indexHandler)
	m.HandleFunc("/vocdash", indexHandler)
	m.HandleFunc("/process/{id}", indexHandler)
	m.HandleFunc("/entity/{id}", indexHandler)
	m.HandleFunc("/envelope/{id}", indexHandler)
	m.HandleFunc("/blocktxs", indexHandler)
	m.HandleFunc("/blocks", indexHandler)
	m.HandleFunc("/block/{id}", indexHandler)
	m.HandleFunc("/tx/{id}", indexHandler)
	m.HandleFunc("/stats", indexHandler)
	m.HandleFunc("/validators", indexHandler)
	m.HandleFunc("/validator/{id}", indexHandler)

	// API Routes
	m.HandleFunc("/ping", PingHandler())
	m.HandleFunc("/config", configHandler(cfg))
	m.HandleFunc("/db/listblocks/", ListBlocksHandler(d))
	m.HandleFunc("/db/listblocksvalidator/", ListBlocksByValidatorHandler(d))
	m.HandleFunc("/db/numblocksvalidator/", NumBlocksByValidatorHandler(d))
	m.HandleFunc("/db/envelope/", GetEnvelopeHandler(d))
	m.HandleFunc("/db/listenvelopes/", ListEnvelopesHandler(d))
	m.HandleFunc("/db/listenvelopesprocess/", ListEnvelopesByProcessHandler(d))
	m.HandleFunc("/db/envprocheight/", EnvelopeHeightByProcessHandler(d))
	m.HandleFunc("/db/entityprocheight/", ProcessHeightByEntityHandler(d))
	m.HandleFunc("/db/block/", GetBlockHandler(d))
	m.HandleFunc("/db/height/", HeightHandler(d))
	m.HandleFunc("/db/heightmap/", HeightMapHandler(d))
	m.HandleFunc("/db/listtxs/", ListTxsHandler(d))
	m.HandleFunc("/db/tx/", GetTxHandler(d))
	m.HandleFunc("/db/txhash/", TxHeightFromHashHandler(d))
	m.HandleFunc("/db/envelopenullifier/", EnvelopeHeightFromNullifierHandler(d))
	m.HandleFunc("/db/validator/", GetValidatorHandler(d))
	m.HandleFunc("/db/listvalidators/", ListValidatorsHandler(d))
	m.HandleFunc("/db/entity/", GetEntityHandler(d))
	m.HandleFunc("/db/process/", GetProcessHandler(d))
	m.HandleFunc("/db/listentities/", ListEntitiesHandler(d))
	m.HandleFunc("/db/listprocesses/", ListProcessesHandler(d))
	m.HandleFunc("/db/listprocessesbyentity/", ListProcessesByEntityHandler(d))
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

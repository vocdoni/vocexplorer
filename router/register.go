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
	m.HandleFunc("/api/listblocks/", ListBlocksHandler(d))
	m.HandleFunc("/api/listblocksvalidator/", ListBlocksByValidatorHandler(d))
	m.HandleFunc("/api/numblocksvalidator/", NumBlocksByValidatorHandler(d))
	m.HandleFunc("/api/envelope/", GetEnvelopeHandler(d))
	m.HandleFunc("/api/listenvelopes/", ListEnvelopesHandler(d))
	m.HandleFunc("/api/listenvelopesprocess/", ListEnvelopesByProcessHandler(d))
	m.HandleFunc("/api/envprocheight/", EnvelopeHeightByProcessHandler(d))
	m.HandleFunc("/api/entityprocheight/", ProcessHeightByEntityHandler(d))
	m.HandleFunc("/api/block/", GetBlockHandler(d))
	m.HandleFunc("/api/height/", HeightHandler(d))
	m.HandleFunc("/api/heightmap/", HeightMapHandler(d))
	m.HandleFunc("/api/listtxs/", ListTxsHandler(d))
	m.HandleFunc("/api/tx/", GetTxHandler(d))
	m.HandleFunc("/api/txhash/", TxHeightFromHashHandler(d))
	m.HandleFunc("/api/envelopenullifier/", EnvelopeHeightFromNullifierHandler(d))
	m.HandleFunc("/api/validator/", GetValidatorHandler(d))
	m.HandleFunc("/api/listvalidators/", ListValidatorsHandler(d))
	m.HandleFunc("/api/entity/", GetEntityHandler(d))
	m.HandleFunc("/api/process/", GetProcessHandler(d))
	m.HandleFunc("/api/listentities/", ListEntitiesHandler(d))
	m.HandleFunc("/api/listprocesses/", ListProcessesHandler(d))
	m.HandleFunc("/api/listprocessesbyentity/", ListProcessesByEntityHandler(d))
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

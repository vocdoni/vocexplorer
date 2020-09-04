package router

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	dvotedb "gitlab.com/vocdoni/go-dvote/db"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/db"
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
	m.HandleFunc("/ping", db.PingHandler())
	m.HandleFunc("/config", configHandler(cfg))
	m.HandleFunc("/db/listblocks/", db.ListBlocksHandler(d))
	m.HandleFunc("/db/listblocksvalidator/", db.ListBlocksByValidatorHandler(d))
	m.HandleFunc("/db/numblocksvalidator/", db.NumBlocksByValidatorHandler(d))
	m.HandleFunc("/db/envelope/", db.GetEnvelopeHandler(d))
	m.HandleFunc("/db/listenvelopes/", db.ListEnvelopesHandler(d))
	m.HandleFunc("/db/listenvelopesprocess/", db.ListEnvelopesByProcessHandler(d))
	m.HandleFunc("/db/envprocheight/", db.EnvelopeHeightByProcessHandler(d))
	m.HandleFunc("/db/entityprocheight/", db.ProcessHeightByEntityHandler(d))
	m.HandleFunc("/db/block/", db.GetBlockHandler(d))
	m.HandleFunc("/db/height/", db.HeightHandler(d))
	m.HandleFunc("/db/heightmap/", db.HeightMapHandler(d))
	m.HandleFunc("/db/listtxs/", db.ListTxsHandler(d))
	m.HandleFunc("/db/tx/", db.GetTxHandler(d))
	m.HandleFunc("/db/txhash/", db.TxHeightFromHashHandler(d))
	m.HandleFunc("/db/envelopenullifier/", db.EnvelopeHeightFromNullifierHandler(d))
	m.HandleFunc("/db/validator/", db.GetValidatorHandler(d))
	m.HandleFunc("/db/listvalidators/", db.ListValidatorsHandler(d))
	m.HandleFunc("/db/entity/", db.GetEntityHandler(d))
	m.HandleFunc("/db/process/", db.GetProcessHandler(d))
	m.HandleFunc("/db/listentities/", db.ListEntitiesHandler(d))
	m.HandleFunc("/db/listprocesses/", db.ListProcessesHandler(d))
	m.HandleFunc("/db/listprocessesbyentity/", db.ListProcessesByEntityHandler(d))
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

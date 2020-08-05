package router

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tendermint/go-amino"
	dvotedb "gitlab.com/vocdoni/go-dvote/db"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/db"
	"gitlab.com/vocdoni/vocexplorer/types"
)

// RegisterRoutes takes a mux and registers all the routes callbacks within this package
func RegisterRoutes(m *mux.Router, cfg *config.Cfg, d *dvotedb.BadgerDB) {
	// Init amino encoder
	var cdc = amino.NewCodec()
	cdc.RegisterConcrete(types.StoreBlock{}, "storeBlock", nil)
	cdc.RegisterConcrete(types.StoreTx{}, "storeTx", nil)

	// Page Routes
	m.HandleFunc("/", indexHandler)
	m.HandleFunc("/vocdash", indexHandler)
	m.HandleFunc("/processes/{id}", indexHandler)
	m.HandleFunc("/entities/{id}", indexHandler)
	m.HandleFunc("/blocktxs", indexHandler)
	m.HandleFunc("/blocks/{id}", indexHandler)

	// API Routes
	m.HandleFunc("/config", configHandler(cfg))
	m.HandleFunc("/db/listblocks/", db.ListBlocksHandler(d, cdc))
	m.HandleFunc("/db/height/", db.HeightHandler(d))
	m.HandleFunc("/db/hash/", db.HashHandler(d))
	m.HandleFunc("/db/listtxs/", db.ListTxsHandler(d, cdc))
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

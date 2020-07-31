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

	m.HandleFunc("/", indexHandler)
	m.HandleFunc("/vocdash", indexHandler)
	m.HandleFunc("/processes/{id}", indexHandler)
	m.HandleFunc("/entities/{id}", indexHandler)
	m.HandleFunc("/config", configHandler(cfg))
	m.HandleFunc("/db/list/", db.ListHandler(d))
	m.HandleFunc("/db/val/", db.ValHandler(d))
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

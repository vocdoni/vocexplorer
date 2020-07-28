package router

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/vocdoni/vocexplorer/config"
)

// RegisterRoutes takes a mux and registers all the routes callbacks within this package
func RegisterRoutes(m *mux.Router, cfg *config.Cfg) {

	m.HandleFunc("/", indexHandler)
	m.HandleFunc("/vocdash", indexHandler)
	m.HandleFunc("/processes/{id}", indexHandler)
	m.HandleFunc("/entities/{id}", indexHandler)
	m.HandleFunc("/config", configHandler(cfg))
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

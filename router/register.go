package router

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/vocdoni/vocexplorer/config"
)

// RegisterRoutes takes a mux and registers all the routes callbacks within this package
func RegisterRoutes(m *mux.Router, cfg *config.Cfg) {

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
	m.HandleFunc("/transaction/{block}/{index}", indexHandler)
	m.HandleFunc("/stats", indexHandler)
	m.HandleFunc("/validators", indexHandler)
	m.HandleFunc("/validator/{id}", indexHandler)

	// API Routes
	m.HandleFunc("/ping", pingHandler())
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

//PingHandler responds to a ping
func pingHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	}
}

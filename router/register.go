package router

import (
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterRoutes takes a mux and registers all the routes callbacks within this package
func RegisterRoutes(m *mux.Router) {

	m.HandleFunc("/", indexHandler)
	m.HandleFunc("/processes", indexHandler)
	m.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	m.NotFoundHandler = http.Handler(http.NotFoundHandler())
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/index.html")
}

// func redirectHandler(w http.ResponseWriter, r *http.Request) {
// 	http.Redirect(w, r, "/", http.StatusFound)
// }

// http.Handle("/", http.FileServer(http.Dir("./static")))

// http.HandleFunc("/processes", redirectHandler)
// http.HandleFunc("/blocks", redirectHandler)
// http.HandleFunc("/txs", redirectHandler)

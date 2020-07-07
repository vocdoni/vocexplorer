package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"gitlab.com/vocdoni/go-dvote/log"
)

func main() {
	if _, err := os.Stat("./static/wasm_exec.js"); os.IsNotExist(err) {
		panic("File not found ./static/wasm_exec.js : find it in $GOROOT/misc/wasm/ note it must be from the same version of go used during compiling")
	}

	urlR, err := url.Parse("http://localhost:8081")
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Infof("Server on: %v\n", urlR)
	// http.HandleFunc("/", indexHandler)
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/processes", redirectHandler)
	http.HandleFunc("/blocks", redirectHandler)
	http.HandleFunc("/txs", redirectHandler)

	err = http.ListenAndServe(urlR.Host, nil)
	if err != nil {
		log.Error(err)
	}
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusFound)
}

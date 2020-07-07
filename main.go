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
	http.Handle("/", http.FileServer(http.Dir("./static")))
	// http.Handle("/processes/", http.FileServer(http.Dir("./static")))
	// http.Handle("/blocks/", http.FileServer(http.Dir("./static")))
	// http.Handle("/txs/", http.FileServer(http.Dir("./static")))
	err = http.ListenAndServe(urlR.Host, nil)
	if err != nil {
		log.Error(err)
	}
}

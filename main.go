package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/gorilla/mux"
	"gitlab.com/vocdoni/go-dvote/log"
	"gitlab.com/vocdoni/vocexplorer/router"
)

func main() {
	// gatewayHost := flag.String("gatewayHost", "ws://0.0.0.0:9090/dvote", "gateway API host to connect to")
	// vochainHost := flag.String("vochainHost", "ws://0.0.0.0:26657/websocket", "gateway API host to connect to")
	hostURL := flag.String("hostURL", "http://localhost:8081", "url to host block explorer")
	logLevel := flag.String("logLevel", "error", "log level <debug, info, warn, error>")
	flag.Parse()
	log.Init(*logLevel, "stdout")

	if _, err := os.Stat("./static/wasm_exec.js"); os.IsNotExist(err) {
		panic("File not found ./static/wasm_exec.js : find it in $GOROOT/misc/wasm/ note it must be from the same version of go used during compiling")
	}

	urlR, err := url.Parse(*hostURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Infof("Server on: %v\n", urlR)

	r := mux.NewRouter()
	router.RegisterRoutes(r)

	err = http.ListenAndServe(urlR.Host, r)
	if err != nil {
		log.Fatal(err)
	}
}

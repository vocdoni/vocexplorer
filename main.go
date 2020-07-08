package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"gitlab.com/vocdoni/go-dvote/log"
	"gitlab.com/vocdoni/vocexplorer/client"
	"nhooyr.io/websocket"
)

func main() {
	gatewayHost := flag.String("gatewayHost", "ws://0.0.0.0:9090/dvote", "gateway API host to connect to")
	vochainHost := flag.String("vochainHost", "ws://0.0.0.0:26657/websocket", "gateway API host to connect to")
	hostURL := flag.String("hostURL", "http://localhost:8081", "url to host block explorer")
	logLevel := flag.String("logLevel", "error", "log level <debug, info, warn, error>")
	flag.Parse()
	log.Init(*logLevel, "stdout")

	log.Infof("connecting to %s", *gatewayHost)
	gw, err := client.New(*gatewayHost)
	if err != nil {
		log.Fatal(err)
	}
	defer gw.Conn.Close(websocket.StatusNormalClosure, "")

	log.Infof("connecting to %s", *vochainHost)
	vc, err := client.New(*vochainHost)
	if err != nil {
		log.Fatal(err)
	}
	defer vc.Conn.Close(websocket.StatusNormalClosure, "")

	if _, err := os.Stat("./static/wasm_exec.js"); os.IsNotExist(err) {
		panic("File not found ./static/wasm_exec.js : find it in $GOROOT/misc/wasm/ note it must be from the same version of go used during compiling")
	}

	urlR, err := url.Parse(*hostURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Infof("Server on: %v\n", urlR)

	// Home serves app which can be navigated via buttons.
	// TODO (possibly): make this stateless, sub-urls handle their own requests instead of redirecting
	// Not sure how to do this with vecty-router
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

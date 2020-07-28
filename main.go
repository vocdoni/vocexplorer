package main

import (
	"net/http"
	"net/url"
	"os"

	flag "github.com/spf13/pflag"

	"github.com/NYTimes/gziphandler"
	"github.com/gorilla/mux"
	"gitlab.com/vocdoni/go-dvote/log"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/router"
	"gitlab.com/vocdoni/vocexplorer/util"
)

func main() {
	var cfg config.Cfg

	gatewayHost := flag.String("gatewayHost", "ws://0.0.0.0:9090/dvote", "gateway API host to connect to")
	tendermintHost := flag.String("vochainHost", "http://0.0.0.0:26657", "gateway API host to connect to")
	refreshTime := flag.Int("refreshTime", 5, "Number of seconds between each content refresh")
	nozip := flag.Bool("disableGzip", false, "use to disable gzip compression on web server")
	hostURL := flag.String("hostURL", "http://localhost:8081", "url to host block explorer")
	logLevel := flag.String("logLevel", "error", "log level <debug, info, warn, error>")
	flag.Parse()

	cfg.GatewayHost = *gatewayHost
	cfg.TendermintHost = *tendermintHost
	cfg.RefreshTime = *refreshTime

	log.Init(*logLevel, "stdout")

	if _, err := os.Stat("./static/wasm_exec.js"); os.IsNotExist(err) {
		panic("File not found ./static/wasm_exec.js : find it in $GOROOT/misc/wasm/ note it must be from the same version of go used during compiling")
	}

	urlR, err := url.Parse(*hostURL)
	if util.ErrPrint(err) {
		return
	}
	log.Infof("Server on: %v\n", urlR)

	r := mux.NewRouter()
	router.RegisterRoutes(r, &cfg)

	if *nozip {
		err = http.ListenAndServe(urlR.Host, r)
		util.ErrFatal(err)
	} else {
		h := gziphandler.GzipHandler(r)
		err = http.ListenAndServe(urlR.Host, h)
		util.ErrFatal(err)
	}

}

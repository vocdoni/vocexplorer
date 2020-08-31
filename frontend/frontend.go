package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gopherjs/vecty"

	"gitlab.com/vocdoni/go-dvote/log"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/api"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/rpc/rpcinit"
	"gitlab.com/vocdoni/vocexplorer/util"
)

//go:generate env GOARCH=wasm GOOS=js go build -o ../static/main.wasm

func main() {
	initFrontend()
	// components.BeforeUnload(func() {
	// 	fmt.Println("Unloading page")
	// 	store.GatewayClient.Close()
	// })
	vecty.SetTitle("Vochain Block Explorer")
	vecty.RenderBody(&Body{})
}

func initFrontend() {
	var cfg *config.Cfg
	resp, err := http.Get("/config")
	util.ErrPrint(err)
	body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	resp.Body.Close()
	if !util.ErrPrint(err) {
		err = json.Unmarshal(body, &cfg)
		util.ErrPrint(err)
	}
	// Init clients with cfg so we don't have to wait for it to store
	initClients(cfg)
	dispatcher.Dispatch(&actions.StoreConfig{Config: *cfg})
	if cfg == nil {
		log.Fatal("Unable to get application configuraion")
	}
	// Wait for store.Config to populate
	i := 0
	for ; store.Config.GatewayHost == "" && store.Config.GatewaySocket == "" && store.Config.RefreshTime == 0; i++ {
		if i > 50 {
			log.Fatal("Config could not be stored")
		}
	}
}

func initClients(cfg *config.Cfg) {
	tm := rpcinit.StartClient(cfg.TendermintHost)
	gw, _ := api.InitGateway(cfg.GatewayHost + cfg.GatewaySocket)
	if tm == nil || gw == nil {
		log.Error("Cannot connect to blockchain clients")
	}
	dispatcher.Dispatch(&actions.TendermintClientInit{Client: tm})
	dispatcher.Dispatch(&actions.GatewayClientInit{Client: gw})
}

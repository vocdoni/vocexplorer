package main

import (
	"encoding/json"
	"net/http"

	"github.com/gopherjs/vecty"
	tmhttp "github.com/tendermint/tendermint/rpc/client/http"
	"gitlab.com/vocdoni/go-dvote/log"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/api"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/rpc/rpcinit"
)

func main() {
	initFrontend()
	vecty.SetTitle("Vochain Block Explorer")
	vecty.RenderBody(&Body{})
}

func initFrontend() {
	var cfg *config.Cfg
	resp, err := http.Get("/config")
	if err != nil {
		log.Error(err)
	}
	err = json.NewDecoder(resp.Body).Decode(&cfg)
	if err != nil {
		log.Error(err)
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
	var tm *tmhttp.HTTP
	var gw *api.GatewayClient
	for i := 0; i < 5 && tm == nil; i++ {
		tm = rpcinit.StartClient(cfg.TendermintHost)
	}
	if tm == nil {
		log.Error("Cannot connect to tendermint api")
	}
	for i := 0; i < 5 && gw == nil; i++ {
		gw, _ = api.InitGateway(cfg.GatewayHost + cfg.GatewaySocket)
	}
	if gw == nil {
		log.Error("Cannot connect to gateway api")
	}
	dispatcher.Dispatch(&actions.TendermintClientInit{Client: tm})
	dispatcher.Dispatch(&actions.GatewayClientInit{Client: gw})
}

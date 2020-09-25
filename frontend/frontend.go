package main

import (
	"encoding/json"
	"net/http"
	"syscall/js"

	"github.com/hexops/vecty"
	"gitlab.com/vocdoni/go-dvote/log"
	"gitlab.com/vocdoni/vocexplorer/api"
	"gitlab.com/vocdoni/vocexplorer/api/rpc"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
)

func main() {
	initFrontend()
	vecty.SetTitle("Vochain Block Explorer")
	vecty.RenderBody(&Body{})
	beforeUnload()
}

func initFrontend() {
	var cfg *config.Cfg
	dispatcher.Dispatch(&actions.ServerConnected{Connected: true})
	resp, err := http.Get("/config")
	if err != nil {
		log.Warn(err)
	}
	err = json.NewDecoder(resp.Body).Decode(&cfg)
	if err != nil {
		log.Error(err)
	}
	// Init clients with cfg so we don't have to wait for it to store
	go initClients(cfg)
	dispatcher.Dispatch(&actions.StoreConfig{Config: *cfg})
	if cfg == nil {
		log.Fatal("Unable to get application configuraion")
	}
	// Wait for store.Config to populate
	i := 0
	for ; store.Config.GatewayHost == "" && store.Config.RefreshTime == 0; i++ {
		if i > 5 {
			log.Fatal("Config could not be stored")
		}
	}
}

func initClients(cfg *config.Cfg) {
	var tm *rpc.TendermintRPC
	var gw *api.GatewayClient
	var err error
	for i := 0; i < 2 && tm == nil; i++ {
		tm = api.StartTendermintClient(cfg.TendermintHost, 2)
	}
	if tm == nil {
		log.Warn("Cannot connect to tendermint api")
	}
	for i := 0; i < 2 && gw == nil; i++ {
		gw, _, err = api.InitGatewayClient(cfg.GatewayHost)
		if err != nil {
			log.Warn(err)
		}
	}
	if gw == nil {
		log.Warn("Cannot connect to gateway api")
	}
	dispatcher.Dispatch(&actions.TendermintClientInit{Client: tm})
	dispatcher.Dispatch(&actions.GatewayClientInit{Client: gw})
	dispatcher.Dispatch(&actions.GatewayConnected{Connected: true})
	store.RedirectChan <- struct{}{}
}

// Beforeunload cleans up before page unload
func beforeUnload() {
	var unloadFunc js.Func
	unloadFunc = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if store.TendermintClient != nil {
			store.TendermintClient.Close()
		}
		if store.GatewayClient != nil {
			store.GatewayClient.Close()
		}
		unloadFunc.Release() // release the function if the button will not be clicked again
		return nil
	})
	js.Global().Call("addEventListener", "beforeunload", unloadFunc)
}

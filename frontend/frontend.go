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
	var unloadFunc js.Func
	unloadFunc = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		close()
		unloadFunc.Release() // release the function if the button will not be clicked again
		return nil
	})
	js.Global().Call("addEventListener", "beforeunload", unloadFunc)

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
	for ; store.Config.GatewayHost == "" && store.Config.RefreshTime == 0; i++ {
		if i > 50 {
			log.Fatal("Config could not be stored")
		}
	}
	beforeUnload()
}

func initClients(cfg *config.Cfg) {
	var tm *rpc.TendermintRPC
	var gw *api.GatewayClient
	for i := 0; i < 5 && tm == nil; i++ {
		tm = api.StartTendermintClient(cfg.TendermintHost, 5)
	}
	if tm == nil {
		log.Error("Cannot connect to tendermint api")
	}
	for i := 0; i < 5 && gw == nil; i++ {
		gw, _ = api.InitGateway(cfg.GatewayHost)
	}
	if gw == nil {
		log.Error("Cannot connect to gateway api")
	}
	dispatcher.Dispatch(&actions.TendermintClientInit{Client: tm})
	dispatcher.Dispatch(&actions.GatewayClientInit{Client: gw})
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

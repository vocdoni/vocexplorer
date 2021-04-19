package main

import (
	"encoding/json"
	"net/http"
	"syscall/js"

	"github.com/hexops/vecty"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/logger"
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
		logger.Warn(err.Error())
	}
	err = json.NewDecoder(resp.Body).Decode(&cfg)
	if err != nil {
		logger.Error(err)
	}
	dispatcher.Dispatch(&actions.StoreConfig{Config: *cfg})
	if cfg == nil {
		logger.Fatal("Unable to get application configuraion")
	}
	if cfg.Dev {
		store.ProcessDomain = "https://app.dev.vocdoni.net/processes/#/0x"
		store.EntityDomain = "https://dev.vocdoni.link/entities/0x"
	}
	// Wait for store.Config to populate
	i := 0
	for ; store.Config.RefreshTime == 0; i++ {
		if i > 5 {
			logger.Fatal("Config could not be stored")
		}
	}
}

// Beforeunload cleans up before page unload
func beforeUnload() {
	var unloadFunc js.Func
	unloadFunc = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		//TODO unload func- is needed without need to close websockets?
		unloadFunc.Release() // release the function if the button will not be clicked again
		return nil
	})
	js.Global().Call("addEventListener", "beforeunload", unloadFunc)
}

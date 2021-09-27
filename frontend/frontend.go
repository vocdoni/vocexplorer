package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"syscall/js"

	"github.com/hexops/vecty"
	"gitlab.com/vocdoni/vocexplorer/client"
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
	// dispatcher.Dispatch(&actions.ServerConnected{Connected: true})
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
	// Wait for store.Config to populate
	i := 0
	for ; store.Config.RefreshTime == 0; i++ {
		if i > 5 {
			logger.Fatal("Config could not be stored")
		}
	}
	if cfg.Network == "dev" {
		dispatcher.Dispatch(&actions.SetLinkURLs{ProcessURL: strings.ReplaceAll(config.ProcessURL, config.DomainKey, config.DevDomain), EntityURL: strings.ReplaceAll(config.EntityURL, config.DomainKey, config.DevDomain)})
	} else if cfg.Network == "stg" {
		dispatcher.Dispatch(&actions.SetLinkURLs{ProcessURL: strings.ReplaceAll(config.ProcessURL, config.DomainKey, config.StgDomain), EntityURL: strings.ReplaceAll(config.EntityURL, config.DomainKey, config.StgDomain)})
	} else {
		dispatcher.Dispatch(&actions.SetLinkURLs{ProcessURL: strings.ReplaceAll(config.ProcessURL, config.DomainKey, config.MainDomain), EntityURL: strings.ReplaceAll(config.EntityURL, config.DomainKey, config.MainDomain)})
	}
	store.Client, err = client.New(store.Config.GatewayUrl)
	if err != nil {
		// If first connection fails, record the error and return so the UI is not frozen.
		// Then attempt to connect 5 more times before giving up.
		logger.Error(err)
		dispatcher.Dispatch(&actions.GatewayConnected{GatewayErr: err})
		go func() {
			for i := 0; i < 5; i++ {
				store.Client, err = client.New(store.Config.GatewayUrl)
				if err != nil {
					logger.Error(err)
				} else {
					return
				}

			}
		}()
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

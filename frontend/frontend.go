package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gopherjs/vecty"

	"gitlab.com/vocdoni/go-dvote/log"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/components"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/util"
)

//go:generate env GOARCH=wasm GOOS=js go build -o ../static/main.wasm

func main() {
	initFrontend()
	components.BeforeUnload(func() {
		fmt.Println("Unloading page")
		store.GatewayClient.Close()
	})
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
	dispatcher.Dispatch(&actions.StoreConfig{Config: *cfg})
	initClients()
}

func initClients() {
	dispatcher.Dispatch(&actions.TendermintClientInit{})
	dispatcher.Dispatch(&actions.GatewayClientInit{})
	if store.GatewayClient == nil || store.TendermintClient == nil {
		log.Error("Cannot connect to blockchain clients")
	}
}

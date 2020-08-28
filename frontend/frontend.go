package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gopherjs/vecty"

	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/components"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/util"
)

//go:generate env GOARCH=wasm GOOS=js go build -o ../static/main.wasm

func main() {
	cfg := initFrontend()
	components.BeforeUnload(func() {
		fmt.Println("Unloading page")
		store.Vochain.Close()
	})
	vecty.SetTitle("Vochain Block Explorer")
	vecty.RenderBody(&Body{Cfg: cfg})
}

func initFrontend() *config.Cfg {
	var cfg *config.Cfg
	resp, err := http.Get("/config")
	util.ErrPrint(err)
	body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	resp.Body.Close()
	if !util.ErrPrint(err) {
		err = json.Unmarshal(body, &cfg)
		util.ErrPrint(err)
	}
	dispatcher.Dispatch(&actions.TendermintClientInit{
		Host: cfg.TendermintHost,
	})
	dispatcher.Dispatch(&actions.VochainClientInit{
		Host: cfg.GatewayHost,
	})
	store.Entities.PagChannel = make(chan int, 50)
	store.Processes.PagChannel = make(chan int, 50)
	return cfg
}

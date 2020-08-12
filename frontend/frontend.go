package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gopherjs/vecty"

	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/util"
)

//go:generate env GOARCH=wasm GOOS=js go build -o ../static/main.wasm

func main() {
	var cfg *config.Cfg
	resp, err := http.Get("/config")
	util.ErrPrint(err)
	body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	resp.Body.Close()
	if !util.ErrPrint(err) {
		err = json.Unmarshal(body, &cfg)
		util.ErrPrint(err)
	}
	vecty.SetTitle("Vochain Block Explorer")
	vecty.RenderBody(&Body{Cfg: cfg})
}

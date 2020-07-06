package main

import (
	"github.com/gopherjs/vecty"

	"gitlab.com/NateWilliams2/vocexplorer/frontend/components"
)

//go:generate env GOARCH=wasm GOOS=js go build -o ../static/main.wasm

func main() {
	vecty.SetTitle("Vochain Block Explorer")
	vecty.RenderBody(&components.Body{})
}

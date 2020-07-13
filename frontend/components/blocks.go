package components

import (
	"fmt"
	"syscall/js"
	"time"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/tendermint/tendermint/rpc/client/http"
	"gitlab.com/vocdoni/vocexplorer/rpc"
)

// BlocksView renders the blocks page
type BlocksView struct {
	vecty.Core
	t *rpc.TendermintInfo
	c *http.HTTP
}

// Render renders the BlocksView component
func (b *BlocksView) Render() vecty.ComponentOrHTML {
	return elem.Div(
		// &Header{currentPage: "blocks"},
		renderTendermintInfo(b.t),
	)
}

func renderTendermintInfo(t *rpc.TendermintInfo) vecty.ComponentOrHTML {
	if t != nil {
		return elem.Div(
			renderBlocks(t),
		)
	}
	return elem.Div(vecty.Text("No info struct rendered"))
}

func renderBlocks(t *rpc.TendermintInfo) vecty.ComponentOrHTML {
	return elem.Div(vecty.Text("Blocks pane"))
}

func initBlocksView(t *rpc.TendermintInfo) *BlocksView {
	js.Global().Set("tendermint", true)
	c, err := rpc.InitClient()
	if err != nil {
		js.Global().Get("alert").Invoke("Unable to connect to Tendermint client. Please see readme file")
		return nil
	}
	// var t *rpc.TendermintInfo
	var blocksView BlocksView
	blocksView.c = c
	blocksView.t = t
	go updateAndRenderBlocks(&blocksView)
	return &blocksView
}

func updateAndRenderBlocks(bv *BlocksView) {
	for js.Global().Get("tendermint").Bool() {
		fmt.Println("Getting tendermint info")
		rpc.UpdateTendermintInfo(bv.c, bv.t)
		time.Sleep(5 * time.Second)
		vecty.Rerender(bv)
	}
	fmt.Println("Closing tendermint updater")
}

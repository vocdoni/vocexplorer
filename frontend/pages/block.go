package pages

import (
	"strconv"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/go-dvote/log"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/components"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/rpc"
	"gitlab.com/vocdoni/vocexplorer/util"
	router "marwan.io/vecty-router"
)

// BlockView renders the Block page
type BlockView struct {
	vecty.Core
	Cfg *config.Cfg
}

// Render renders the BlockView component
func (home *BlockView) Render() vecty.ComponentOrHTML {
	height, err := strconv.ParseInt(router.GetNamedVar(home)["id"], 0, 64)
	util.ErrPrint(err)
	dispatcher.Dispatch(&actions.SetCurrentBlock{Block: rpc.GetBlock(store.TendermintClient, height)})
	if store.Blocks.CurrentBlock == nil {
		log.Errorf("Block unavailable")
		return elem.Div(
			elem.Main(vecty.Text("Block not available")),
		)
	}
	return elem.Div(&components.BlockContents{})
}

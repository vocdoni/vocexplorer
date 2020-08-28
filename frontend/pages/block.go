package pages

import (
	"strconv"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/go-dvote/log"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/components"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/rpc"
	"gitlab.com/vocdoni/vocexplorer/util"
	router "marwan.io/vecty-router"
)

// BlocksView renders the Blocks page
type BlockView struct {
	vecty.Core
	Cfg *config.Cfg
}

// Render renders the BlocksView component
func (home *BlockView) Render() vecty.ComponentOrHTML {
	height, err := strconv.ParseInt(router.GetNamedVar(home)["id"], 0, 64)
	util.ErrPrint(err)
	block := rpc.GetBlock(store.Tendermint, height)
	if block == nil {
		log.Errorf("Block unavailable")
		return elem.Div(
			elem.Main(vecty.Text("Block not available")),
		)
	}
	return &components.BlockContents{
		Block: block.Block,
		Hash:  block.BlockID.Hash,
	}
}

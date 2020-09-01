package pages

import (
	"strconv"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/go-dvote/log"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/frontend/components"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/rpc"
	"gitlab.com/vocdoni/vocexplorer/util"
	router "marwan.io/vecty-router"
)

// BlocksView is a pretty page for all our blockchain statistics
type BlocksView struct {
	vecty.Core
	Cfg *config.Cfg
}

// Render renders the Blocks component
func (home *BlocksView) Render() vecty.ComponentOrHTML {
	return home.Component()
}

// Component generates the actual BlocksView component
func (home *BlocksView) Component() vecty.ComponentOrHTML {
	id, ok := router.GetNamedVar(home)["id"]
	if !ok {
		return components.Container(
			elem.Section(
				bootstrap.Card(bootstrap.CardParams{
					Body: vecty.List{
						elem.Heading3(
							vecty.Text("Blocks"),
						),
						vecty.Text("Blocks list will be here"),
					},
				}),
			),
		)
	}
	height, err := strconv.ParseInt(id, 0, 64)
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

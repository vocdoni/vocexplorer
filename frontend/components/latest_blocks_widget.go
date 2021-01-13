package components

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/vocdoni/vocexplorer/api/dbtypes"
	"github.com/vocdoni/vocexplorer/frontend/bootstrap"
	"github.com/vocdoni/vocexplorer/frontend/store"
)

//LatestBlocksWidget is a component for a widget of recent blocks
type LatestBlocksWidget struct {
	vecty.Core
}

//Render renders the LatestBlocksWidget component
func (b *LatestBlocksWidget) Render() vecty.ComponentOrHTML {

	var blockList []vecty.MarkupOrChild

	max := 4
	for i := len(store.Blocks.Blocks) - 1; i >= len(store.Blocks.Blocks)-max; i-- {
		if dbtypes.BlockIsEmpty(store.Blocks.Blocks[i]) {
			continue
		}
		blockList = append(blockList, elem.Div(
			vecty.Markup(vecty.Class("card-deck-col")),
			BlockCard(store.Blocks.Blocks[i]),
		))
	}
	if len(blockList) == 0 {
		return elem.Div(vecty.Text("Loading Blocks..."))
	}
	blockList = append(blockList, vecty.Markup(vecty.Class("responsive-card-deck")))

	return elem.Section(
		vecty.Markup(vecty.Class("recent-blocks")),
		bootstrap.Card(bootstrap.CardParams{
			Body: vecty.List{
				elem.Heading2(vecty.Text("Latest blocks")),
				elem.Div(
					blockList...,
				),
			},
		}),
	)
}

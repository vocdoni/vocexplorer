package components

import (
	"fmt"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/rpc"
	"gitlab.com/vocdoni/vocexplorer/types"
)

type LatestBlocksWidget struct {
	vecty.Core
	T *rpc.TendermintInfo
}

func (b *LatestBlocksWidget) Render() vecty.ComponentOrHTML {

	var blockList []vecty.MarkupOrChild

	max := 4
	for i := len(b.T.BlockList) - 1; i >= len(b.T.BlockList)-max; i-- {
		if types.BlockIsEmpty(b.T.BlockList[i]) {
			continue
		}
		blockList = append(blockList, BlockCard(b.T.BlockList[i]))
	}
	if len(blockList) == 0 {
		fmt.Println("No blocks available")
		return elem.Div(vecty.Text("Loading Blocks..."))
	}
	blockList = append(blockList, vecty.Markup(vecty.Class("responsive-card-deck")))

	return elem.Section(
		vecty.Markup(vecty.Class("recent-blocks")),
		elem.Heading4(vecty.Text("Latest blocks")),
		elem.Div(
			blockList...,
		),
	)
}

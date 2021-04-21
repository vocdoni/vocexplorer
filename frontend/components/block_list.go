package components

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"

	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
)

// BlockList is the block list component
type BlockList struct {
	vecty.Core
}

// Render renders the block list component
func (b *BlockList) Render() vecty.ComponentOrHTML {
	if len(store.Blocks.Blocks) > 0 {
		p := &Pagination{
			TotalPages:      int(store.Blocks.Count) / config.ListSize,
			TotalItems:      &store.Blocks.Count,
			CurrentPage:     &store.Blocks.Pagination.CurrentPage,
			RefreshCh:       store.Blocks.Pagination.PagChannel,
			ListSize:        config.ListSize,
			DisableUpdate:   &store.Blocks.Pagination.DisableUpdate,
			SearchCh:        store.Blocks.Pagination.SearchChannel,
			Searching:       &store.Blocks.Pagination.Search,
			RenderSearchBar: true,
		}
		p.RenderFunc = func(index int) vecty.ComponentOrHTML {
			return renderBlocks(p, index)
		}

		return elem.Section(
			vecty.Markup(vecty.Class("list", "paginated")),
			bootstrap.Card(bootstrap.CardParams{
				Body: vecty.List{
					elem.Heading1(
						vecty.Text("Blocks"),
					),
					p,
				},
			}),
		)
	}
	return elem.Div(vecty.Text("Waiting for blockchain info..."))
}

func renderBlocks(p *Pagination, index int) vecty.ComponentOrHTML {
	var blockList []vecty.MarkupOrChild

	for i := len(store.Blocks.Blocks) - 1; i >= len(store.Blocks.Blocks)-p.ListSize; i-- {
		if dbtypes.BlockIsEmpty(store.Blocks.Blocks[i]) {
			continue
		}
		blockList = append(blockList, elem.Div(
			vecty.Markup(vecty.Class("paginated-card")),
			BlockCard(store.Blocks.Blocks[i]),
		))
	}
	if len(blockList) == 0 {
		if *p.Searching {
			return elem.Div(vecty.Text("No blocks found"))
		}
		return elem.Div(vecty.Text("Loading Blocks..."))
	}
	blockList = append(blockList, vecty.Markup(vecty.Class("row")))
	return elem.Div(
		blockList...,
	)
}

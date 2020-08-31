package components

import (
	"fmt"
	"strconv"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/gopherjs/vecty/prop"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/types"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// BlockList is the block list component
type BlockList struct {
	vecty.Core
	currentPage int
}

// Render renders the block list component
func (b *BlockList) Render() vecty.ComponentOrHTML {
	if len(store.Blocks.Blocks) > 0 {
		p := &Pagination{
			TotalPages:      int(store.Blocks.Count) / config.ListSize,
			TotalItems:      &store.Blocks.Count,
			CurrentPage:     &b.currentPage,
			RefreshCh:       store.Blocks.Pagination.PagChannel,
			ListSize:        config.ListSize,
			DisableUpdate:   &store.Blocks.Pagination.DisableUpdate,
			RenderSearchBar: true,
		}
		p.RenderFunc = func(index int) vecty.ComponentOrHTML {
			return renderBlocks(p, index)
		}
		p.SearchBar = func(self *Pagination) vecty.ComponentOrHTML {
			return elem.Input(vecty.Markup(
				event.Input(func(e *vecty.Event) {
					search := e.Target.Get("value").String()
					index, err := strconv.Atoi(e.Target.Get("value").String())
					if err != nil || index < 0 || index > int(*self.TotalItems) || search == "" {
						*self.CurrentPage = 0
						dispatcher.Dispatch(&actions.DisableBlockUpdate{Disabled: false})
						self.RefreshCh <- *self.CurrentPage * config.ListSize
					} else {
						*self.CurrentPage = util.Max(int(*self.TotalItems)-index-1, 0) / config.ListSize
						dispatcher.Dispatch(&actions.DisableBlockUpdate{Disabled: true})
						self.RefreshCh <- int(*self.TotalItems) - index
					}
					vecty.Rerender(self)
				}),
				prop.Placeholder("search blocks"),
			))
		}

		return elem.Section(
			vecty.Markup(vecty.Class("list", "paginated")),
			bootstrap.Card(bootstrap.CardParams{
				Body: vecty.List{
					elem.Heading3(
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
		if types.BlockIsEmpty(store.Blocks.Blocks[i]) {
			continue
		}
		blockList = append(blockList, elem.Div(
			vecty.Markup(vecty.Class("paginated-card")),
			BlockCard(store.Blocks.Blocks[i]),
		))
	}
	if len(blockList) == 0 {
		fmt.Println("No blocks available")
		return elem.Div(vecty.Text("Loading Blocks..."))
	}
	blockList = append(blockList, vecty.Markup(vecty.Class("row")))
	return elem.Div(
		blockList...,
	)
}

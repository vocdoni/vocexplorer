package components

import (
	"fmt"
	"strconv"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/gopherjs/vecty/prop"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/rpc"
	"gitlab.com/vocdoni/vocexplorer/types"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// BlockList is the block list component
type BlockList struct {
	vecty.Core
	currentPage   int
	disableUpdate *bool
	refreshCh     chan int
	t             *rpc.TendermintInfo
}

// Render renders the block list component
func (b *BlockList) Render() vecty.ComponentOrHTML {
	if b.t != nil && b.t.ResultStatus != nil {
		p := &Pagination{
			TotalPages:      int(b.t.TotalBlocks) / config.ListSize,
			TotalItems:      &b.t.TotalBlocks,
			CurrentPage:     &b.currentPage,
			RefreshCh:       b.refreshCh,
			ListSize:        config.ListSize,
			RenderSearchBar: true,
		}
		p.RenderFunc = func(index int) vecty.ComponentOrHTML {
			return renderBlocks(p, b.t, index)
		}
		p.SearchBar = func(self *Pagination) vecty.ComponentOrHTML {
			return elem.Input(vecty.Markup(
				event.Input(func(e *vecty.Event) {
					search := e.Target.Get("value").String()
					index, err := strconv.Atoi(e.Target.Get("value").String())
					if err != nil || index < 0 || index > int(*self.TotalItems) || search == "" {
						*self.CurrentPage = 0
						*b.disableUpdate = false
						self.RefreshCh <- *self.CurrentPage * config.ListSize
					} else {
						*self.CurrentPage = util.Max(int(*self.TotalItems)-index-1, 0) / config.ListSize
						*b.disableUpdate = true
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

func renderBlocks(p *Pagination, t *rpc.TendermintInfo, index int) vecty.ComponentOrHTML {
	var blockList []vecty.MarkupOrChild

	for i := len(t.BlockList) - 1; i >= len(t.BlockList)-p.ListSize; i-- {
		if types.BlockIsEmpty(t.BlockList[i]) {
			continue
		}
		blockList = append(blockList, elem.Div(
			vecty.Markup(vecty.Class("paginated-card")),
			BlockCard(t.BlockList[i]),
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

package components

import (
	"strconv"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/gopherjs/vecty/prop"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/rpc"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// BlockList is the block list component
type BlockList struct {
	vecty.Core
	t           *rpc.TendermintInfo
	currentPage int
	refreshCh   chan int
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
						self.RefreshCh <- *self.CurrentPage * config.ListSize
					} else {
						*self.CurrentPage = util.Max(int(*self.TotalItems)-index-1, 0) / config.ListSize
						self.RefreshCh <- int(*self.TotalItems) - index
					}
					vecty.Rerender(self)
				}),
				prop.Placeholder("search blocks & txs"),
			))
		}

		return elem.Div(
			vecty.Markup(vecty.Class("recent-blocks")),
			elem.Heading3(
				vecty.Text("Blocks"),
			),
			p,
		)
	}
	return elem.Div(vecty.Text("Waiting for blockchain info..."))
}

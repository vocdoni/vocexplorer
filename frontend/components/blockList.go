package components

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/gopherjs/vecty/prop"
	"github.com/xeonx/timeago"
	"gitlab.com/vocdoni/vocexplorer/config"
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

func renderBlocks(p *Pagination, t *rpc.TendermintInfo, index int) vecty.ComponentOrHTML {
	var blockList []vecty.MarkupOrChild

	empty := p.ListSize
	for i := len(t.BlockList) - 1; i >= len(t.BlockList)-p.ListSize; i-- {
		if types.BlockIsEmpty(t.BlockList[i]) {
			empty--
		}
		block := t.BlockList[i]
		// for i, block := range t.BlockList {
		blockList = append(blockList, renderBlock(block))
	}
	if empty == 0 {
		fmt.Println("No blocks available")
		return elem.Div(vecty.Text("Loading Blocks..."))
	}
	blockList = append(blockList, vecty.Markup(vecty.Class("responsive-card-deck")))
	return elem.Div(
		blockList...,
	)
}

func renderBlock(block *types.StoreBlock) vecty.ComponentOrHTML {
	var tm time.Time
	var err error
	if block.GetTime() != nil {
		tm, err = ptypes.Timestamp(block.GetTime())
		util.ErrPrint(err)
	}
	return elem.Div(vecty.Markup(vecty.Class("card-deck-col")),
		elem.Div(vecty.Markup(vecty.Class("card")),
			elem.Div(
				vecty.Markup(vecty.Class("card-header")),
				elem.Anchor(
					vecty.Markup(
						vecty.Class("nav-link"),
						vecty.Attribute("href", "/blocks/"+util.IntToString(block.GetHeight())),
					),
					vecty.Text(util.IntToString(block.GetHeight())),
				),
			),
			elem.Div(
				vecty.Markup(vecty.Class("card-body")),
				elem.Div(
					vecty.Markup(vecty.Class("block-card-heading")),
					elem.Div(
						vecty.Text(util.IntToString(block.GetNumTxs())+" transactions"),
					),
					elem.Div(
						vecty.Text(timeago.English.Format(tm)),
					),
				),
				elem.Div(
					elem.Div(
						vecty.Markup(vecty.Class("dt")),
						vecty.Text("Hash"),
					),
					elem.Div(
						vecty.Markup(vecty.Class("dd")),
						vecty.Text(util.HexToString(block.GetHash())),
					),
				),
				elem.Div(
					elem.Div(
						vecty.Markup(vecty.Class("dt")),
						vecty.Text("Proposer"),
					),
					elem.Div(
						vecty.Markup(vecty.Class("dd")),
						vecty.Text(util.HexToString(block.GetProposer())),
					),
				),
			),
		),
	)
}

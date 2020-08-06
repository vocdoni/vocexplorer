package components

import (
	"encoding/json"
	"fmt"
	"strconv"

	humanize "github.com/dustin/go-humanize"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/gopherjs/vecty/prop"
	dvotetypes "gitlab.com/vocdoni/go-dvote/types"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/rpc"
	"gitlab.com/vocdoni/vocexplorer/types"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// TxList is the tx list component
type TxList struct {
	vecty.Core
	t             *rpc.TendermintInfo
	currentPage   int
	refreshCh     chan int
	disableUpdate *bool
}

// Render renders the tx list component
func (b *TxList) Render() vecty.ComponentOrHTML {
	if b.t != nil && b.t.ResultStatus != nil {
		p := &Pagination{
			TotalPages:      int(b.t.TotalTxs) / config.ListSize,
			TotalItems:      &b.t.TotalTxs,
			CurrentPage:     &b.currentPage,
			RefreshCh:       b.refreshCh,
			ListSize:        config.ListSize,
			RenderSearchBar: true,
		}
		p.RenderFunc = func(index int) vecty.ComponentOrHTML {
			return renderTxs(p, b.t, index)
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
				prop.Placeholder("search transactions"),
			))
		}

		return elem.Div(
			vecty.Markup(vecty.Class("recent-blocks")),
			elem.Heading3(
				vecty.Text("Txs"),
			),
			p,
		)
	}
	return elem.Div(vecty.Text("Waiting for txchain info..."))
}

func renderTxs(p *Pagination, t *rpc.TendermintInfo, index int) vecty.ComponentOrHTML {
	var txList []vecty.MarkupOrChild

	empty := p.ListSize
	for i := p.ListSize - 1; i >= 0; i-- {
		if t.TxList[i].IsEmpty() {
			empty--
		}
		tx := t.TxList[i]
		// for i, tx := range t.TxList {
		txList = append(txList, renderTx(&tx))
	}
	if empty == 0 {
		fmt.Println("No txs available")
		return elem.Div(vecty.Text("Loading Txs..."))
	}
	txList = append(txList, vecty.Markup(vecty.Class("responsive-card-deck")))
	return elem.Div(
		txList...,
	)
}

func renderTx(tx *types.SendTx) vecty.ComponentOrHTML {
	var rawTx dvotetypes.Tx
	err := json.Unmarshal(tx.Store.Tx, &rawTx)
	util.ErrPrint(err)
	return elem.Div(vecty.Markup(vecty.Class("card-deck-col")),
		elem.Div(vecty.Markup(vecty.Class("card")),
			elem.Div(
				vecty.Markup(vecty.Class("card-header")),
				elem.Anchor(
					vecty.Markup(
						vecty.Class("nav-link"),
						vecty.Attribute("href", "/txs/"+util.IntToString((tx.Store.TxHeight))),
					),
					vecty.Text(util.IntToString(tx.Store.TxHeight)),
				),
			),
			elem.Div(
				vecty.Markup(vecty.Class("card-body")),
				elem.Div(
					vecty.Markup(vecty.Class("dt")),
					vecty.Text(humanize.Ordinal(int(tx.Store.Index+1))+" transaction on block "),
					elem.Anchor(
						vecty.Markup(
							vecty.Attribute("href", "/blocks/"+util.IntToString(tx.Store.Height)),
						),
						vecty.Text(util.IntToString(tx.Store.Height)),
					),
				),
				elem.Div(
					elem.Div(
						vecty.Markup(vecty.Class("dt")),
						vecty.Text("Hash"),
					),
					elem.Div(
						vecty.Markup(vecty.Class("dd")),
						vecty.Text(tx.Hash.String()),
					),
				),
				vecty.If(rawTx.Type != "",
					elem.Div(
						vecty.Text("Type: "+rawTx.Type),
					)),
			),
		),
	)
}

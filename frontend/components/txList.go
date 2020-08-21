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
	currentPage   int
	disableUpdate *bool
	refreshCh     chan int
	t             *rpc.TendermintInfo
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

		return elem.Section(
			vecty.Markup(vecty.Class("list", "paginated")),
			elem.Heading3(
				vecty.Text("Transactions"),
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
		if types.TxIsEmpty(t.TxList[i]) {
			empty--
		} else {
			tx := t.TxList[i]
			txList = append(txList, renderTx(tx))
		}
	}
	if empty == 0 {
		fmt.Println("No txs available")
		return elem.Div(vecty.Text("Loading Txs..."))
	}

	return elem.Div(
		txList...,
	)
}

func renderTx(tx *types.SendTx) vecty.ComponentOrHTML {
	var rawTx dvotetypes.Tx
	err := json.Unmarshal(tx.Store.Tx, &rawTx)
	util.ErrPrint(err)
	return elem.Div(
		vecty.Markup(vecty.Class("tile", rawTx.Type)),
		elem.Div(
			vecty.Markup(vecty.Class("tile-body")),
			elem.Div(
				vecty.Markup(vecty.Class("type")),
				elem.Div(
					elem.Span(
						vecty.Markup(vecty.Class("title")),
						vecty.Text(rawTx.Type),
					),
				),
			),
			elem.Div(
				vecty.Markup(vecty.Class("contents")),
				elem.Div(
					vecty.Markup(vecty.Class("dt")),
					vecty.Text(
						fmt.Sprintf("%s transaction on block ", humanize.Ordinal(int(tx.Store.Index+1))),
					),
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
						vecty.Text(util.HexToString(tx.GetHash())),
					),
				),
			),
		),
	)
}

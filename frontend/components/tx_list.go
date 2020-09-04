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
	"gitlab.com/vocdoni/go-dvote/log"
	dvotetypes "gitlab.com/vocdoni/go-dvote/types"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/proto"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// TxList is the tx list component
type TxList struct {
	vecty.Core
}

// Render renders the tx list component
func (b *TxList) Render() vecty.ComponentOrHTML {
	if store.Stats.ResultStatus != nil {
		p := &Pagination{
			TotalPages:      int(store.Transactions.Count) / config.ListSize,
			TotalItems:      &store.Transactions.Count,
			CurrentPage:     &store.Transactions.Pagination.CurrentPage,
			RefreshCh:       store.Transactions.Pagination.PagChannel,
			ListSize:        config.ListSize,
			DisableUpdate:   &store.Transactions.Pagination.DisableUpdate,
			RenderSearchBar: true,
		}
		p.RenderFunc = func(index int) vecty.ComponentOrHTML {
			return renderTxs(p, index)
		}
		p.SearchBar = func(self *Pagination) vecty.ComponentOrHTML {
			return elem.Input(vecty.Markup(
				event.Input(func(e *vecty.Event) {
					search := e.Target.Get("value").String()
					index, err := strconv.Atoi(e.Target.Get("value").String())
					if err != nil || index < 0 || index > int(*self.TotalItems) || search == "" {
						*self.CurrentPage = 0
						dispatcher.Dispatch(&actions.DisableTransactionsUpdate{Disabled: false})
						self.RefreshCh <- *self.CurrentPage * config.ListSize
					} else {
						*self.CurrentPage = util.Max(int(*self.TotalItems)-index-1, 0) / config.ListSize
						dispatcher.Dispatch(&actions.DisableTransactionsUpdate{Disabled: true})
						self.RefreshCh <- int(*self.TotalItems) - index
					}
					vecty.Rerender(self)
				}),
				prop.Placeholder("search transactions"),
			))
		}

		return elem.Section(
			vecty.Markup(vecty.Class("list", "paginated")),
			bootstrap.Card(bootstrap.CardParams{
				Body: vecty.List{
					elem.Heading3(
						vecty.Text("Transactions"),
					),
					p,
				},
			}),
		)
	}
	return elem.Div(vecty.Text("Waiting for blockchain info..."))
}

func renderTxs(p *Pagination, index int) vecty.ComponentOrHTML {
	var txList []vecty.MarkupOrChild

	empty := p.ListSize
	for i := p.ListSize - 1; i >= 0; i-- {
		if proto.TxIsEmpty(store.Transactions.Transactions[i]) {
			empty--
		} else {
			tx := store.Transactions.Transactions[i]
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

func renderTx(tx *proto.SendTx) vecty.ComponentOrHTML {
	var rawTx dvotetypes.Tx
	err := json.Unmarshal(tx.Store.Tx, &rawTx)
	if err != nil {
		log.Error(err)
	}
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
					vecty.Text(fmt.Sprintf("# %d", tx.Store.GetTxHeight())),
				),
				elem.Div(
					vecty.Markup(vecty.Class("dt")),
					vecty.Text(
						fmt.Sprintf("%s transaction on block ", humanize.Ordinal(int(tx.Store.Index+1))),
					),
					Link(
						"/block/"+util.IntToString(tx.Store.Height),
						util.IntToString(tx.Store.Height),
						"",
					),
				),
				elem.Div(
					elem.Div(
						vecty.Markup(vecty.Class("dt")),
						vecty.Text("Hash"),
					),
					elem.Div(
						vecty.Markup(vecty.Class("dd")),
						Link(
							"/tx/"+util.IntToString(tx.Store.TxHeight),
							util.HexToString(tx.GetHash()),
							"",
						),
					),
				),
			),
		),
	)
}

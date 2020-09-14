package components

import (
	"encoding/json"
	"fmt"

	humanize "github.com/dustin/go-humanize"
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"gitlab.com/vocdoni/go-dvote/log"
	dvotetypes "gitlab.com/vocdoni/go-dvote/types"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
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
			SearchCh:        store.Transactions.Pagination.SearchChannel,
			Searching:       &store.Transactions.Pagination.Search,
			RenderSearchBar: true,
		}
		p.RenderFunc = func(index int) vecty.ComponentOrHTML {
			return renderTxs(p, index)
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

	for i := len(store.Transactions.Transactions) - 1; i >= len(store.Transactions.Transactions)-p.ListSize; i-- {
		if proto.TxIsEmpty(store.Transactions.Transactions[i]) {
			continue
		}
		txList = append(txList, renderTx(store.Transactions.Transactions[i]))
	}
	if len(txList) == 0 {
		if *p.Searching {
			return elem.Div(vecty.Text("No Transactions Found With Given ID"))
		}
		return elem.Div(vecty.Text("Loading Transactions..."))
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
						vecty.Text(fmt.Sprintf("#%d", tx.Store.GetTxHeight())),
					),
					elem.Span(
						vecty.Markup(vecty.Class("title")),
						vecty.Text(rawTx.Type),
					),
				),
			),
			elem.Div(
				vecty.Markup(vecty.Class("contents")),
				elem.Div(
					elem.Div(
						Link(
							"/transaction/"+util.IntToString(tx.Store.TxHeight),
							util.HexToString(tx.GetHash()),
							"",
						),
					),
					vecty.Text(
						fmt.Sprintf("%s transaction on block ", humanize.Ordinal(int(tx.Store.Index+1))),
					),
					Link(
						"/block/"+util.IntToString(tx.Store.Height),
						util.IntToString(tx.Store.Height),
						"",
					),
				),
			),
		),
	)
}

package components

import (
	"fmt"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"

	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/frontend/store/storeutil"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// TransactionList is the transaction list component
type TransactionList struct {
	vecty.Core
}

// Render renders the transaction list component
func (b *TransactionList) Render() vecty.ComponentOrHTML {
	if len(store.Transactions.Transactions) > 0 {
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
			SearchPrompt:    "search by tx height",
		}
		p.RenderFunc = func(index int) vecty.ComponentOrHTML {
			return renderTransactions(p, index)
		}

		return elem.Section(
			vecty.Markup(vecty.Class("list", "paginated")),
			bootstrap.Card(bootstrap.CardParams{
				Body: vecty.List{
					elem.Heading1(
						vecty.Text("Transactions"),
					),
					p,
				},
			}),
		)
	}
	return Unavailable("Waiting for transactionchain info...", "")
}

func renderTransactions(p *Pagination, index int) vecty.ComponentOrHTML {
	var txList []vecty.MarkupOrChild

	for i := len(store.Transactions.Transactions) - 1; i >= util.Max(len(store.Transactions.Transactions)-p.ListSize, 0); i-- {
		txList = append(txList, renderTx(store.Transactions.Transactions[i]))
	}
	if len(txList) == 0 {
		if *p.Searching {
			return elem.Div(vecty.Text("No transactions found"))
		}
		return elem.Div(vecty.Text("Loading Transactions..."))
	}
	return elem.Div(
		txList...,
	)
}

func renderTx(tx *storeutil.FullTransaction) vecty.ComponentOrHTML {
	if tx.Decoded == nil || tx.Package == nil {
		return elem.Div(vecty.Text("Loading Transaction..."))
	}
	tp := util.GetTransactionType(tx.Decoded.RawTx)
	if tp == "" {
		tp = "Unknown"
	}
	return elem.Div(
		vecty.Markup(vecty.Class("tile", tp)),
		elem.Div(
			vecty.Markup(vecty.Class("tile-body")),
			elem.Div(
				vecty.Markup(vecty.Class("type")),
				elem.Div(
					elem.Span(
						vecty.Text(fmt.Sprintf("#%d", tx.Package.ID)),
					),
					elem.Span(
						vecty.Markup(vecty.Class("title")),
						vecty.Text(util.GetTransactionName(tp)),
					),
				),
			),
			elem.Div(
				vecty.Markup(vecty.Class("contents")),
				elem.Div(
					elem.Div(
						Link(
							"/transaction/"+util.IntToString(tx.Package.BlockHeight)+"/"+util.IntToString(tx.Package.Index),
							util.HexToString(tx.Package.Hash),
							"",
						),
					),
				),
			),
		),
	)
}

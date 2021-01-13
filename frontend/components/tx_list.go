package components

import (
	"fmt"

	humanize "github.com/dustin/go-humanize"
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/vocdoni/vocexplorer/api/dbtypes"
	"github.com/vocdoni/vocexplorer/config"
	"github.com/vocdoni/vocexplorer/frontend/bootstrap"
	"github.com/vocdoni/vocexplorer/frontend/store"
	"github.com/vocdoni/vocexplorer/logger"
	"github.com/vocdoni/vocexplorer/util"
	"go.vocdoni.io/proto/build/go/models"
	"google.golang.org/protobuf/proto"
)

// TxList is the tx list component
type TxList struct {
	vecty.Core
}

// Render renders the tx list component
func (b *TxList) Render() vecty.ComponentOrHTML {
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
				elem.Heading1(
					vecty.Text("Transactions"),
				),
				p,
			},
		}),
	)
}

func renderTxs(p *Pagination, index int) vecty.ComponentOrHTML {
	var txList []vecty.MarkupOrChild

	for i := len(store.Transactions.Transactions) - 1; i >= len(store.Transactions.Transactions)-p.ListSize; i-- {
		if dbtypes.TxIsEmpty(store.Transactions.Transactions[i]) {
			continue
		}
		txList = append(txList, renderTx(store.Transactions.Transactions[i]))
	}
	if len(txList) == 0 {
		if *p.Searching {
			return elem.Div(vecty.Text("No transactions found"))
		}
		return elem.Div(vecty.Text("No transactions available"))
	}

	return elem.Div(
		txList...,
	)
}

func renderTx(tx *dbtypes.Transaction) vecty.ComponentOrHTML {
	var rawTx models.Tx
	err := proto.Unmarshal(tx.Tx, &rawTx)
	if err != nil {
		logger.Error(err)
	}
	txType := util.GetTransactionType(&rawTx)
	return elem.Div(
		vecty.Markup(vecty.Class("tile", txType)),
		elem.Div(
			vecty.Markup(vecty.Class("tile-body")),
			elem.Div(
				vecty.Markup(vecty.Class("type")),
				elem.Div(
					elem.Span(
						vecty.Text(fmt.Sprintf("#%d", tx.TxHeight)),
					),
					elem.Span(
						vecty.Markup(vecty.Class("title")),
						vecty.Text(util.GetTransactionName(txType)),
					),
				),
			),
			elem.Div(
				vecty.Markup(vecty.Class("contents")),
				elem.Div(
					elem.Div(
						Link(
							"/transaction/"+util.IntToString(tx.TxHeight),
							util.HexToString(tx.Hash),
							"",
						),
					),
					vecty.Text(
						fmt.Sprintf("%s transaction on block ", humanize.Ordinal(int(tx.Index+1))),
					),
					Link(
						"/block/"+util.IntToString(tx.Height),
						util.IntToString(tx.Height),
						"",
					),
				),
			),
		),
	)
}

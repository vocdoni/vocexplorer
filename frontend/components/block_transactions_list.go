package components

import (
	"encoding/json"
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	dvotetypes "gitlab.com/vocdoni/go-dvote/types"
	"gitlab.com/vocdoni/vocexplorer/api/dbtypes"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/logger"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// BlockTransactionsListView renders the transaction pagination for a block
type BlockTransactionsListView struct {
	vecty.Core
}

// Render renders the BlockTransactionsListView component
func (b *BlockTransactionsListView) Render() vecty.ComponentOrHTML {
	numTxs := len(store.Blocks.CurrentBlock.Block.Data.Txs)
	if numTxs == 0 {
		return elem.Preformatted(
			vecty.Markup(vecty.Class("empty")),
			vecty.Text("No transactions"),
		)
	}
	p := &Pagination{
		TotalPages:      numTxs / config.ListSize,
		TotalItems:      &numTxs,
		CurrentPage:     &store.Blocks.TransactionPagination.CurrentPage,
		RefreshCh:       store.Blocks.TransactionPagination.PagChannel,
		ListSize:        config.ListSize,
		DisableUpdate:   &store.Blocks.TransactionPagination.DisableUpdate,
		SearchCh:        store.Blocks.TransactionPagination.SearchChannel,
		Searching:       &store.Blocks.TransactionPagination.Search,
		RenderSearchBar: false,
	}
	p.RenderFunc = func(index int) vecty.ComponentOrHTML {
		return renderBlockTxs(p, index)
	}

	return elem.Section(
		vecty.Markup(vecty.Class("list", "paginated")),
		bootstrap.Card(bootstrap.CardParams{
			Body: vecty.List{
				elem.Heading2(
					vecty.Text("Transactions"),
				),
				p,
			},
		}),
	)
}

func renderBlockTxs(p *Pagination, index int) vecty.ComponentOrHTML {
	var txList []vecty.MarkupOrChild

	for i := len(store.Blocks.CurrentTxs) - 1; i >= len(store.Blocks.CurrentTxs)-p.ListSize; i-- {
		if dbtypes.TxIsEmpty(store.Blocks.CurrentTxs[i]) {
			continue
		}
		txList = append(txList, renderBlockTx(store.Blocks.CurrentTxs[i]))
	}
	if len(txList) == 0 {
		if *p.Searching {
			return elem.Div(vecty.Text("No transactions found"))
		}
		return elem.Div(vecty.Text("Loading transactions..."))
	}

	return elem.Div(
		txList...,
	)
}

func renderBlockTx(tx *dbtypes.Transaction) vecty.ComponentOrHTML {
	var rawTx dvotetypes.Tx
	err := json.Unmarshal(tx.Tx, &rawTx)
	if err != nil {
		logger.Error(err)
	}
	return elem.Div(
		vecty.Markup(vecty.Class("tile", rawTx.Type)),
		elem.Div(
			vecty.Markup(vecty.Class("tile-body")),
			elem.Div(
				vecty.Markup(vecty.Class("type")),
				elem.Div(
					elem.Span(
						vecty.Text(fmt.Sprintf("#%d", tx.Index+1)),
					),
					elem.Span(
						vecty.Markup(vecty.Class("title")),
						vecty.Text(util.GetTransactionName(rawTx.Type)),
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
						fmt.Sprintf("%s transaction on the blockchain ", humanize.Ordinal(int(tx.TxHeight))),
					),
				),
			),
		),
	)
}

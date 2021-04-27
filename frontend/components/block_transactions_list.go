package components

import (
	"fmt"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/tendermint/tendermint/crypto/tmhash"

	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/logger"
	"gitlab.com/vocdoni/vocexplorer/util"
	"go.vocdoni.io/proto/build/go/models"
	"google.golang.org/protobuf/proto"
)

// BlockTransactionsListView renders the transaction pagination for a block
type BlockTransactionsListView struct {
	vecty.Core
}

// Render renders the BlockTransactionsListView component
func (b *BlockTransactionsListView) Render() vecty.ComponentOrHTML {
	numTxs := int(store.Blocks.CurrentBlock.NumTxs)
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

	for i := len(store.Blocks.CurrentTxs) - 1; i >= util.Max(len(store.Blocks.CurrentTxs)-p.ListSize, 0); i-- {
		txList = append(txList, renderBlockTx(store.Blocks.CurrentTxs[i], int(store.Blocks.CurrentBlock.NumTxs)-(p.ListSize**p.CurrentPage)-(len(store.Blocks.CurrentTxs)-i-1)))
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

func renderBlockTx(tx *models.SignedTx, index int) vecty.ComponentOrHTML {
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
						vecty.Text(fmt.Sprintf("#%d", index)),
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
							"/transaction/"+util.IntToString(store.Blocks.CurrentBlock.Height)+"/"+util.IntToString(index),
							util.HexToString(tmhash.Sum(tx.Tx)),
							"",
						),
					),
					// vecty.Text(
					// 	fmt.Sprintf("%s transaction on the blockchain ", humanize.Ordinal(int(tx.TxHeight))),
					// ),
				),
			),
		),
	)
}

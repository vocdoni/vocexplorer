package components

import (
	"fmt"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"

	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/frontend/update"
	"gitlab.com/vocdoni/vocexplorer/logger"
	"gitlab.com/vocdoni/vocexplorer/util"
	indexertypes "go.vocdoni.io/dvote/vochain/scrutinizer/indexertypes"
)

// BlockContents renders block contents
type BlockContents struct {
	vecty.Core
	vecty.Mounter
	Rendered    bool
	Unavailable bool
}

// Mount triggers when BlockContents renders
func (c *BlockContents) Mount() {
	if !c.Rendered {
		c.Rendered = true
		vecty.Rerender(c)
	}
}

// Render renders the BlockContents component
func (c *BlockContents) Render() vecty.ComponentOrHTML {
	if !c.Rendered {
		return LoadingBar()
	}
	if c.Unavailable {
		return Unavailable("Block unavailable", "")
	}
	if store.Blocks.CurrentBlock == nil {
		return Unavailable("Loading block...", "")
	}
	return Container(
		vecty.Markup(vecty.Attribute("id", "main")),
		renderServerConnectionBanner(),
		elem.Section(
			vecty.Markup(vecty.Class("details-view", "no-column")),
			elem.Div(
				vecty.Markup(vecty.Class("row")),
				elem.Div(
					vecty.Markup(vecty.Class("main-column")),
					bootstrap.Card(bootstrap.CardParams{
						Body: BlockView(),
					}),
				),
			),
		),
		elem.Section(
			vecty.Markup(vecty.Class("row")),
			elem.Div(
				vecty.Markup(vecty.Class("col-12")),
				bootstrap.Card(bootstrap.CardParams{
					Body: c.BlockDetails(),
				}),
			),
		),
	)
}

// UpdateBlockContents keeps the block contents up to date
func UpdateBlockContents(d *BlockContents) {
	// Set block to nil so previous block is not displayed
	dispatcher.Dispatch(&actions.SetCurrentBlockTransactionList{
		TransactionList: []*indexertypes.TxMetadata{},
	})
	dispatcher.Dispatch(&actions.SetCurrentBlock{Block: nil})
	dispatcher.Dispatch(&actions.EnableAllUpdates{})
	// Fetch block contents
	d.fetchBlock()
	ticker := time.NewTicker(time.Duration(store.Config.RefreshTime) * 5 * time.Second)
	if !update.CheckCurrentPage("block", ticker) {
		return
	}
	if !d.Unavailable {
		updateBlockTransactions(int(store.Blocks.CurrentBlock.NumTxs) - store.Blocks.TransactionPagination.Index - config.ListSize)
	}
	for {
		select {
		case <-store.RedirectChan:
			if !update.CheckCurrentPage("block", ticker) {
				return
			}
		case <-ticker.C:
			if !update.CheckCurrentPage("block", ticker) {
				return
			}
			// If block never loaded, load it
			if d.Unavailable {
				d.fetchBlock()
			}
		case i := <-store.Blocks.TransactionPagination.PagChannel:
		txloop:
			for {
				// If many indices waiting in buffer, scan to last one.
				select {
				case i = <-store.Blocks.TransactionPagination.PagChannel:
				default:
					break txloop
				}
			}
			if !update.CheckCurrentPage("block", ticker) {
				return
			}
			if !d.Unavailable {
				// update the current page of txs
				dispatcher.Dispatch(&actions.BlockTransactionsIndexChange{Index: i})
				updateBlockTransactions(int(store.Blocks.CurrentBlock.NumTxs) - store.Blocks.TransactionPagination.Index - config.ListSize)
			}
		}
	}
}

func (c *BlockContents) fetchBlock() {
	block, err := store.Client.GetBlock(uint32(store.Blocks.CurrentBlockHeight))
	if err != nil {
		logger.Error(err)
		c.Unavailable = true
		dispatcher.Dispatch(&actions.SetCurrentBlock{Block: nil})
	} else {
		c.Unavailable = false
		dispatcher.Dispatch(&actions.SetCurrentBlock{Block: block})
	}
}

func updateBlockTransactions(index int) {
	listSize := config.ListSize
	if index < 0 {
		listSize += index
		index = 0
	}
	logger.Info(fmt.Sprintf("Getting %d transactions from index %d\n", listSize, index))

	if store.Blocks.CurrentBlock != nil {
		txs, err := store.Client.GetTxListForBlock(uint32(store.Blocks.CurrentBlock.Height), index, listSize)
		if err != nil {
			logger.Error(err)
		}
		dispatcher.Dispatch(&actions.SetCurrentBlockTransactionList{TransactionList: txs})
	}
}

//BlockView renders a single block card
func BlockView() vecty.List {
	return vecty.List{
		elem.Heading1(
			vecty.Markup(vecty.Class("card-title")),
			vecty.Text("Block details"),
		),
		elem.Heading2(
			vecty.Text(fmt.Sprintf("Block Height: %d", store.Blocks.CurrentBlock.Height)),
		),
		elem.Div(
			vecty.Markup(vecty.Class("details")),
			elem.Span(
				vecty.If(int(store.Blocks.CurrentBlock.NumTxs) == 1,
					vecty.Text(fmt.Sprintf("%d transaction", int(store.Blocks.CurrentBlock.NumTxs)))),
				vecty.If(int(store.Blocks.CurrentBlock.NumTxs) != 1,
					vecty.Text(fmt.Sprintf("%d transactions", int(store.Blocks.CurrentBlock.NumTxs)))),
			),
			elem.Span(vecty.Text(
				humanize.Time(store.Blocks.CurrentBlock.Timestamp),
			)),
		),
		elem.HorizontalRule(),
		elem.DescriptionList(
			elem.DefinitionTerm(
				vecty.Text("Hash"),
			),
			elem.Description(
				vecty.Text(util.HexToString(store.Blocks.CurrentBlock.Hash)),
			),
			elem.DefinitionTerm(
				vecty.Text("Parent hash"),
			),
			elem.Description(
				Link(
					fmt.Sprintf("/block/%d", store.Blocks.CurrentBlock.Height-1),
					util.HexToString(store.Blocks.CurrentBlock.LastBlockHash),
					"",
				),
			),
			elem.DefinitionTerm(
				vecty.Text("Proposer Address"),
			),
			elem.Description(
				Link(
					"/validator/"+util.HexToString(store.Blocks.CurrentBlock.ProposerAddress),
					util.HexToString(store.Blocks.CurrentBlock.ProposerAddress),
					"",
				),
			),
			elem.DefinitionTerm(
				vecty.Text("Total transactions"),
			),
			elem.Description(
				vecty.Text(fmt.Sprintf("%d", int(store.Blocks.CurrentBlock.NumTxs))),
			),
			elem.DefinitionTerm(
				vecty.Text("Time"),
			),
			elem.Description(
				vecty.Text(store.Blocks.CurrentBlock.Timestamp.String()),
			),
		),
	}
}

//BlockTab component displays a tab for the block page
type BlockTab struct {
	*Tab
}

func (b *BlockTab) store() string {
	return store.Blocks.Pagination.Tab
}
func (b *BlockTab) dispatch() interface{} {
	return &actions.BlocksTabChange{
		Tab: b.alias(),
	}
}

//BlockDetails displays the details for a single block
func (c *BlockContents) BlockDetails() vecty.List {
	transactions := &BlockTab{&Tab{
		Text:  "Transactions",
		Alias: "transactions",
	}}
	return vecty.List{
		elem.Navigation(
			vecty.Markup(vecty.Class("tabs")),
			elem.UnorderedList(
				TabLink(c, transactions),
			),
		),
		elem.Div(
			vecty.Markup(vecty.Class("tabs-content")),
			TabContents(transactions, &BlockTransactionsListView{}),
		),
	}
}

func preformattedBlockHeader(block *indexertypes.BlockMetadata) vecty.ComponentOrHTML {
	return elem.Preformatted(elem.Code(vecty.Text(block.String())))
}

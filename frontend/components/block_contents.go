package components

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	dvotetypes "gitlab.com/vocdoni/go-dvote/types"
	"gitlab.com/vocdoni/vocexplorer/api"
	"gitlab.com/vocdoni/vocexplorer/api/dbtypes"
	"gitlab.com/vocdoni/vocexplorer/api/tmtypes"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/frontend/update"
	"gitlab.com/vocdoni/vocexplorer/logger"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// BlockContents renders block contents
type BlockContents struct {
	vecty.Core
	vecty.Mounter
	Rendered bool
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
	if store.Blocks.CurrentBlock == nil {
		return Container(
			vecty.Markup(vecty.Attribute("id", "main")),
			renderServerConnectionBanner(),
			elem.Section(
				bootstrap.Card(bootstrap.CardParams{
					Body: vecty.List{
						elem.Heading2(
							vecty.Text("Loading block..."),
						),
					},
				}),
			),
		)
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
						Body: BlockView(store.Blocks.CurrentBlock.Block),
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
	dispatcher.Dispatch(&actions.EnableAllUpdates{})
	// Fetch block contents
	logger.Info("getting block")
	block, ok := api.GetBlock(store.Blocks.CurrentBlockHeight)
	if block != nil && ok {
		logger.Info("got block")
		dispatcher.Dispatch(&actions.SetCurrentBlock{Block: block})
	}
	ticker := time.NewTicker(time.Duration(store.Config.RefreshTime) * time.Second)
	if !update.CheckCurrentPage("block", ticker) {
		return
	}
	updateBlockTransactions()
	for {
		select {
		case <-store.RedirectChan:
			if !update.CheckCurrentPage("block", ticker) {
				return
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
			dispatcher.Dispatch(&actions.BlockTransactionsIndexChange{Index: i})
			updateBlockTransactions()
			// update the current page of txs
		}
	}
}

func updateBlockTransactions() {
	if store.Blocks.CurrentBlock != nil {
		maxIndex := len(store.Blocks.CurrentBlock.Block.Data.Txs) - 1
		maxIndex = util.Min(util.Max(maxIndex-store.Blocks.TransactionPagination.Index, 0), maxIndex)
		var rawTx dvotetypes.Tx
		var transactions [config.ListSize]*dbtypes.Transaction
		wg := new(sync.WaitGroup)
		for i := 0; i < config.ListSize && maxIndex-i >= 0; i++ {
			err := json.Unmarshal(store.Blocks.CurrentBlock.Block.Data.Txs[maxIndex-i], &rawTx)
			if err != nil {
				logger.Error(err)
			}
			// Asynchronously fetch all txs
			wg.Add(1)
			go func(index int) {
				hashString := fmt.Sprintf("%X", store.Blocks.CurrentBlock.Block.Data.Txs[maxIndex-index].Hash())
				fullTransaction, ok := api.GetTxByHash(hashString)
				if ok {
					transactions[index] = fullTransaction
				}
				wg.Done()
			}(i)
		}
		wg.Wait()
		reverseTxList(&transactions)
		dispatcher.Dispatch(&actions.SetCurrentBlockTransactionList{TransactionList: transactions})
	}
}

//BlockView renders a single block card
func BlockView(block *tmtypes.Block) vecty.List {
	return vecty.List{
		elem.Heading1(
			vecty.Markup(vecty.Class("card-title")),
			vecty.Text("Block details"),
		),
		elem.Heading2(
			vecty.Text(fmt.Sprintf("Block Height: %d", block.Header.Height)),
		),
		elem.Div(
			vecty.Markup(vecty.Class("details")),
			elem.Span(
				vecty.Text(fmt.Sprintf("%d transactions", len(block.Data.Txs))),
			),
			elem.Span(
				vecty.Text(humanize.Bytes(uint64(block.Size()))),
			),
			elem.Span(vecty.Text(
				humanize.Time(block.Header.Time),
			)),
		),
		elem.HorizontalRule(),
		elem.DescriptionList(
			elem.DefinitionTerm(
				vecty.Text("Hash"),
			),
			elem.Description(
				vecty.Text(block.Header.Hash().String()),
			),
			elem.DefinitionTerm(
				vecty.Text("Parent hash"),
			),
			elem.Description(
				Link(
					fmt.Sprintf("/block/%d", block.Header.Height-1),
					block.Header.LastBlockID.Hash.String(),
					"",
				),
			),
			elem.DefinitionTerm(
				vecty.Text("Proposer Address"),
			),
			elem.Description(
				Link(
					"/validator/"+block.ProposerAddress.String(),
					block.ProposerAddress.String(),
					"",
				),
			),
			elem.DefinitionTerm(
				vecty.Text("Total transactions"),
			),
			elem.Description(
				vecty.Text(fmt.Sprintf("%d", (len(block.Data.Txs)))),
			),
			elem.DefinitionTerm(
				vecty.Text("Block size"),
			),
			elem.Description(
				vecty.Text(fmt.Sprintf("%d bytes", block.Size())),
			),
			elem.DefinitionTerm(
				vecty.Text("Time"),
			),
			elem.Description(
				vecty.Text(block.Header.Time.UTC().String()),
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
	header := &BlockTab{&Tab{
		Text:  "Header",
		Alias: "header",
	}}
	evidence := &BlockTab{&Tab{
		Text:  "Evidence",
		Alias: "evidence",
	}}
	lastCommit := &BlockTab{&Tab{
		Text:  "Last commit",
		Alias: "last-commit",
	}}

	return vecty.List{
		elem.Navigation(
			vecty.Markup(vecty.Class("tabs")),
			elem.UnorderedList(
				TabLink(c, transactions),
				TabLink(c, header),
				TabLink(c, evidence),
				TabLink(c, lastCommit),
			),
		),
		elem.Div(
			vecty.Markup(vecty.Class("tabs-content")),
			TabContents(transactions, &BlockTransactionsListView{}),
			TabContents(header, preformattedBlockHeader(store.Blocks.CurrentBlock.Block)),
			TabContents(evidence, preformattedBlockEvidence(store.Blocks.CurrentBlock.Block)),
			TabContents(lastCommit, preformattedBlockLastCommit(store.Blocks.CurrentBlock.Block)),
		),
	}
}

func preformattedBlockEvidence(block *tmtypes.Block) vecty.ComponentOrHTML {
	var evidence []byte
	var err error

	if len(block.Evidence.Evidence) <= 0 {
		return elem.Preformatted(
			vecty.Markup(vecty.Class("empty")),
			vecty.Text("No evidence"),
		)
	}

	evidence, err = json.MarshalIndent(block.Evidence, "", "\t")
	if err != nil {
		logger.Error(err)
	}

	return elem.Preformatted(elem.Code(vecty.Text(string(evidence))))
}

func preformattedBlockLastCommit(block *tmtypes.Block) vecty.ComponentOrHTML {
	commit, err := json.MarshalIndent(block.LastCommit, "", "\t")
	if err != nil {
		logger.Error(err)
	}

	return elem.Preformatted(elem.Code(vecty.Text(string(commit))))
}

func preformattedBlockHeader(block *tmtypes.Block) vecty.ComponentOrHTML {
	header, err := json.MarshalIndent(block.Header, "", "\t")
	if err != nil {
		logger.Error(err)
	}

	return elem.Preformatted(elem.Code(vecty.Text(string(header))))
}

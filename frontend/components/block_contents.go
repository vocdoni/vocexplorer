package components

import (
	"crypto/sha256"
	"fmt"
	"sync"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/api"
	"gitlab.com/vocdoni/vocexplorer/api/dbtypes"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/frontend/update"
	"gitlab.com/vocdoni/vocexplorer/logger"
	"gitlab.com/vocdoni/vocexplorer/util"
	"go.vocdoni.io/proto/build/go/models"
	"google.golang.org/protobuf/proto"
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
						Body: BlockView(store.Blocks.CurrentBlock),
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
	dispatcher.Dispatch(&actions.SetCurrentBlock{Block: nil})
	dispatcher.Dispatch(&actions.EnableAllUpdates{})
	// Fetch block contents
	logger.Info("getting block")
	block, ok := api.GetBlock(store.Blocks.CurrentBlockHeight)
	if block != nil && ok {
		d.Unavailable = false
		dispatcher.Dispatch(&actions.SetCurrentBlock{Block: block})
	} else {
		d.Unavailable = true
		dispatcher.Dispatch(&actions.SetCurrentBlock{Block: nil})
		return
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
		maxIndex := len(store.Blocks.CurrentBlock.Data) - 1
		maxIndex = util.Min(util.Max(maxIndex-store.Blocks.TransactionPagination.Index, 0), maxIndex)
		var rawTx models.Tx
		var transactions [config.ListSize]*dbtypes.Transaction
		wg := new(sync.WaitGroup)
		for i := 0; i < config.ListSize && maxIndex-i >= 0; i++ {
			err := proto.Unmarshal(store.Blocks.CurrentBlock.Data[maxIndex-i], &rawTx)
			if err != nil {
				logger.Error(err)
			}
			// Asynchronously fetch all txs
			wg.Add(1)
			go func(index int) {
				hashString := fmt.Sprintf("%X", txHash(store.Blocks.CurrentBlock.Data[maxIndex-index]))
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
func BlockView(block *api.Block) vecty.List {
	return vecty.List{
		elem.Heading1(
			vecty.Markup(vecty.Class("card-title")),
			vecty.Text("Block details"),
		),
		elem.Heading2(
			vecty.Text(fmt.Sprintf("Block Height: %d", block.Height)),
		),
		elem.Div(
			vecty.Markup(vecty.Class("details")),
			elem.Span(
				vecty.Text(fmt.Sprintf("%d transactions", len(block.Data))),
			),
			elem.Span(
				vecty.Text(humanize.Bytes(uint64(block.Size))),
			),
			elem.Span(vecty.Text(
				humanize.Time(block.Time),
			)),
		),
		elem.HorizontalRule(),
		elem.DescriptionList(
			elem.DefinitionTerm(
				vecty.Text("Hash"),
			),
			elem.Description(
				vecty.Text(block.Hash),
			),
			elem.DefinitionTerm(
				vecty.Text("Parent hash"),
			),
			elem.Description(
				Link(
					fmt.Sprintf("/block/%d", block.Height-1),
					block.LastBlockID,
					"",
				),
			),
			elem.DefinitionTerm(
				vecty.Text("Proposer Address"),
			),
			elem.Description(
				Link(
					"/validator/"+block.ProposerAddress,
					block.ProposerAddress,
					"",
				),
			),
			elem.DefinitionTerm(
				vecty.Text("Total transactions"),
			),
			elem.Description(
				vecty.Text(fmt.Sprintf("%d", (len(block.Data)))),
			),
			elem.DefinitionTerm(
				vecty.Text("Block size"),
			),
			elem.Description(
				vecty.Text(fmt.Sprintf("%d bytes", block.Size)),
			),
			elem.DefinitionTerm(
				vecty.Text("Time"),
			),
			elem.Description(
				vecty.Text(block.Time.UTC().String()),
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
			TabContents(header, preformattedBlockHeader(store.Blocks.CurrentBlock)),
			TabContents(evidence, preformattedBlockEvidence(store.Blocks.CurrentBlock)),
			TabContents(lastCommit, preformattedBlockLastCommit(store.Blocks.CurrentBlock)),
		),
	}
}

func preformattedBlockEvidence(block *api.Block) vecty.ComponentOrHTML {

	if len(block.Evidence) <= 0 {
		return elem.Preformatted(
			vecty.Markup(vecty.Class("empty")),
			vecty.Text("No evidence"),
		)
	}

	return elem.Preformatted(elem.Code(vecty.Text(block.Evidence)))
}

func preformattedBlockLastCommit(block *api.Block) vecty.ComponentOrHTML {
	return elem.Preformatted(elem.Code(vecty.Text(string(block.LastCommit))))
}

func preformattedBlockHeader(block *api.Block) vecty.ComponentOrHTML {
	return elem.Preformatted(elem.Code(vecty.Text(string(block.Header))))
}

func txHash(tx []byte) []byte {
	// Sum returns the SHA256 of the bz.
	h := sha256.Sum256(tx)
	return h[:]
}

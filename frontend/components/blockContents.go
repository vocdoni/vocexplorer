package components

import (
	"encoding/json"
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	tmtypes "github.com/tendermint/tendermint/types"
	dvotetypes "gitlab.com/vocdoni/go-dvote/types"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/api"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/rpc"
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
		return elem.Div(vecty.Text("Loading..."))
	}
	if store.Blocks.CurrentBlock == nil {
		return Container(
			elem.Section(
				bootstrap.Card(bootstrap.CardParams{
					Body: vecty.List{
						elem.Heading3(
							vecty.Text("Block does not exist"),
						),
					},
				}),
			),
		)
	}
	return Container(
		elem.Section(
			vecty.Markup(vecty.Class("details-view")),
			elem.Div(
				vecty.Markup(vecty.Class("row")),
				elem.Div(
					vecty.Markup(vecty.Class("main-column")),
					bootstrap.Card(bootstrap.CardParams{
						Body: BlockView(store.Blocks.CurrentBlock.Block),
					}),
				),
				elem.Div(
					vecty.Markup(vecty.Class("extra-column")),
					bootstrap.Card(bootstrap.CardParams{
						Header: elem.Heading4(vecty.Text("Validator")),
						Body: elem.Div(
							Link(
								"/validator/"+store.Blocks.CurrentBlock.Block.ValidatorsHash.String(),
								store.Blocks.CurrentBlock.Block.ValidatorsHash.String(),
								"",
							),
						),
						ClassNames: []string{"validator"},
					}),
					bootstrap.Card(bootstrap.CardParams{
						Body:       vecty.Text("card body"),
						ClassNames: []string{"flex-grow-1", "ml-0", "ml-md-5", "ml-lg-0"},
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

// UpdateAndRenderBlockContents keeps the block contents up to date
func UpdateAndRenderBlockContents(d *BlockContents) {
	actions.EnableUpdates()
	// Fetch block contents
	block := rpc.GetBlock(store.TendermintClient, store.Blocks.CurrentBlockHeight)
	dispatcher.Dispatch(&actions.SetCurrentBlock{Block: block})
	var rawTx dvotetypes.Tx
	var txHeights []int64
	for _, tx := range store.Blocks.CurrentBlock.Block.Data.Txs {
		err := json.Unmarshal(tx, &rawTx)
		util.ErrPrint(err)
		hashString := fmt.Sprintf("%X", tx.Hash())
		txHeight, _ := api.GetTxHeightFromHash(hashString)
		txHeights = append(txHeights, txHeight)
	}
	dispatcher.Dispatch(&actions.SetCurrentBlockTxHeights{Heights: txHeights})
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
				vecty.Text(fmt.Sprintf("%d bytes", block.Size())),
			),
			elem.Span(vecty.Text(fmt.Sprintf(
				"%s (%s)",
				humanize.Time(block.Header.Time),
				block.Header.Time.Local().String(),
			))),
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
				vecty.Markup(vecty.Class("dt")),
				vecty.Text("Proposer Address"),
			),
			elem.Description(
				Link(
					"/validator/"+block.ProposerAddress.String(),
					block.ProposerAddress.String(),
					"",
				),
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
			TabContents(transactions, preformattedBlockTransactions(store.Blocks.CurrentBlock.Block)),
			TabContents(header, preformattedBlockHeader(store.Blocks.CurrentBlock.Block)),
			TabContents(evidence, preformattedBlockEvidence(store.Blocks.CurrentBlock.Block)),
			TabContents(lastCommit, preformattedBlockLastCommit(store.Blocks.CurrentBlock.Block)),
		),
	}
}

func preformattedBlockTransactions(block *tmtypes.Block) vecty.ComponentOrHTML {
	var rawTx dvotetypes.Tx
	numTx := 0
	var txHeight int64
	data := []vecty.MarkupOrChild{vecty.Text("Transactions: [\n")}
	for i, tx := range block.Data.Txs {
		numTx++
		err := json.Unmarshal(tx, &rawTx)
		util.ErrPrint(err)
		hashString := fmt.Sprintf("%X", tx.Hash())
		var hashElement vecty.ComponentOrHTML
		if len(store.Blocks.CurrentBlockTxHeights) > i {
			txHeight = store.Blocks.CurrentBlockTxHeights[i]
			hashElement = Link(
				"/tx/"+util.IntToString(txHeight),
				hashString,
				"",
			)
		} else {
			hashElement = elem.Div(vecty.Markup(vecty.Class("nav-link")), vecty.Text(hashString))
		}
		data = append(
			data,
			elem.Div(
				vecty.Text("\tHash: "),
				hashElement,
				vecty.Text(fmt.Sprintf(" (%d bytes) Type: %s, \n", len(tx), rawTx.Type)),
			),
		)
	}
	data = append(data, vecty.Text("]"))
	if numTx == 0 {
		return elem.Preformatted(
			vecty.Markup(vecty.Class("empty")),
			vecty.Text("No transactions"),
		)
	}

	return elem.Preformatted(elem.Code(data...))
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
	util.ErrPrint(err)

	return elem.Preformatted(elem.Code(vecty.Text(string(evidence))))
}

func preformattedBlockLastCommit(block *tmtypes.Block) vecty.ComponentOrHTML {
	commit, err := json.MarshalIndent(block.LastCommit, "", "\t")
	util.ErrPrint(err)

	return elem.Preformatted(elem.Code(vecty.Text(string(commit))))
}

func preformattedBlockHeader(block *tmtypes.Block) vecty.ComponentOrHTML {
	header, err := json.MarshalIndent(block.Header, "", "\t")
	util.ErrPrint(err)

	return elem.Preformatted(elem.Code(vecty.Text(string(header))))
}

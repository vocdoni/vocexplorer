package components

import (
	"encoding/json"
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	tmtypes "github.com/tendermint/tendermint/types"
	dvotetypes "gitlab.com/vocdoni/go-dvote/types"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// BlockContents renders block contents
type BlockContents struct {
	vecty.Core
	Block *tmtypes.Block
	Hash  tmbytes.HexBytes
	// BlockDetails vecty.ComponentOrHTML
}

// Render renders the BlockContents component
func (c *BlockContents) Render() vecty.ComponentOrHTML {
	return Container(
		elem.Section(
			vecty.Markup(vecty.Class("details-view")),
			elem.Div(
				vecty.Markup(vecty.Class("row")),
				elem.Div(
					vecty.Markup(vecty.Class("main-column")),
					bootstrap.Card(bootstrap.CardParams{
						Body: BlockView(c.Block),
					}),
				),
				elem.Div(
					vecty.Markup(vecty.Class("extra-column")),
					bootstrap.Card(bootstrap.CardParams{
						Header: elem.Heading4(vecty.Text("Validator")),
						Body: elem.Div(
							elem.Anchor(
								vecty.Markup(vecty.Attribute("href", "/validators/"+c.Block.ValidatorsHash.String())),
								vecty.Text(c.Block.ValidatorsHash.String()),
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
				elem.Anchor(
					vecty.Markup(
						vecty.Attribute("href", fmt.Sprintf("/blocks/%d", block.Header.Height-1)),
					),
					vecty.Text(block.Header.LastBlockID.Hash.String()),
				),
			),
			elem.DefinitionTerm(
				vecty.Markup(vecty.Class("dt")),
				vecty.Text("Proposer Address"),
			),
			elem.Description(
				elem.Anchor(
					vecty.Markup(
						vecty.Attribute("href", "/validators/"+block.ProposerAddress.String()),
					),
					vecty.Text(block.ProposerAddress.String()),
				),
			),
		),
	}
}

type BlockTab struct {
	*Tab
}

func (b *BlockTab) store() string {
	return store.BlockTabActive
}
func (b *BlockTab) dispatch() interface{} {
	return &actions.BlocksTabChange{
		Tab: b.alias(),
	}
}

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
			TabContents(transactions, preformattedBlockTransactions(c.Block)),
			TabContents(header, preformattedBlockHeader(c.Block)),
			TabContents(evidence, preformattedBlockEvidence(c.Block)),
			TabContents(lastCommit, preformattedBlockLastCommit(c.Block)),
		),
	}
}

func preformattedBlockTransactions(block *tmtypes.Block) vecty.ComponentOrHTML {
	var rawTx dvotetypes.Tx
	numTx := 0
	data := []vecty.MarkupOrChild{vecty.Text("Transactions: [\n")}
	for _, tx := range block.Data.Txs {
		numTx++
		err := json.Unmarshal(tx, &rawTx)
		util.ErrPrint(err)
		data = append(
			data,
			elem.Div(
				vecty.Text("\tHash: "),
				elem.Anchor(
					vecty.Markup(
						vecty.Attribute("href", fmt.Sprintf("/db/txhash/?hash=%X", tx.Hash())),
					),
					vecty.Text(fmt.Sprintf("%X", tx.Hash())),
				),
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

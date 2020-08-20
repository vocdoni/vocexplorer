package components

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/prop"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	tmtypes "github.com/tendermint/tendermint/types"
	dvotetypes "gitlab.com/vocdoni/go-dvote/types"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// BlockContents renders block contents
type BlockContents struct {
	vecty.Core
	Block *tmtypes.Block
	Hash  tmbytes.HexBytes
}

// Render renders the BlockContents component
func (contents *BlockContents) Render() vecty.ComponentOrHTML {
	return Container(
		elem.Section(
			vecty.Markup(vecty.Class("details-view")),
			elem.Div(
				vecty.Markup(vecty.Class("row")),
				elem.Div(
					vecty.Markup(vecty.Class("main-column")),
					bootstrap.Card(bootstrap.CardParams{
						Body: BlockDetails(contents.Block),
					}),
				),
				elem.Div(
					vecty.Markup(vecty.Class("extra-column")),
					bootstrap.Card(bootstrap.CardParams{
						Body:       vecty.Text("card body"),
						ClassNames: []string{"validators"},
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
					Header: elem.Heading2(
						vecty.Text("2do"),
					),
					Body: renderBlockContents(contents.Block),
				}),
			),
		),
	)
}

func BlockDetails(block *tmtypes.Block) vecty.List {
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
			elem.Span(
				vecty.Text(
					fmt.Sprintf("%s (%s)", humanize.Time(block.Header.Time), block.Header.Time.Local().String()),
				),
			),
		),
		elem.HorizontalRule(),
		elem.DescriptionList(
			elem.Definition(
				vecty.Text("Hash"),
			),
			elem.Description(
				vecty.Text(block.Header.Hash().String()),
			),
			elem.Definition(
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
		),
	}
}

func renderBlockHeader(numTxs int, hash tmbytes.HexBytes, height int64, tm time.Time) vecty.ComponentOrHTML {
	return elem.Div(vecty.Markup(vecty.Class("card-deck-col")),
		elem.Div(vecty.Markup(vecty.Class("card")),
			elem.Div(
				elem.Heading2(
					vecty.Markup(vecty.Class("card-header")),
					vecty.Text("Block "+util.IntToString(height)),
				),
			),
			elem.Div(
				vecty.Markup(vecty.Class("card-body")),
				elem.Div(
					vecty.Markup(vecty.Class("block-card-heading")),
					elem.Div(
						vecty.Text(util.IntToString(numTxs)+" transactions"),
					),
					elem.Div(
						vecty.Text(humanize.Time(tm)),
					),
				),
				elem.Div(
					elem.Div(
						vecty.Markup(vecty.Class("dt")),
						vecty.Text("Hash"),
					),
					elem.Div(
						vecty.Markup(vecty.Class("dd")),
						vecty.Text(hash.String()),
					),
				),
			),
		),
	)
}

func renderBlockContents(block *tmtypes.Block) vecty.ComponentOrHTML {
	header, err := json.MarshalIndent(block.Header, "", "\t")
	util.ErrPrint(err)
	var rawTx dvotetypes.Tx
	numTx := 0
	data := []vecty.MarkupOrChild{vecty.Text("Transactions: [\n")}
	for _, tx := range block.Data.Txs {
		numTx++
		err = json.Unmarshal(tx, &rawTx)
		util.ErrPrint(err)
		data = append(
			data,
			elem.Div(
				vecty.Text("    Hash: "),
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
		data = []vecty.MarkupOrChild{vecty.Text("No transactions")}
	}
	transactions := elem.Preformatted(elem.Code(data...))
	var evidence []byte
	if len(block.Evidence.Evidence) > 0 {
		evidence, err = json.MarshalIndent(block.Evidence, "", "\t")
		util.ErrPrint(err)
	} else {
		evidence = []byte("No evidence")
	}
	commit, err := json.MarshalIndent(block.LastCommit, "", "\t")
	util.ErrPrint(err)
	accordionName := "accordionBlock"
	return elem.Div(
		vecty.Markup(vecty.Class("accordion"), prop.ID(accordionName)),
		renderCollapsible("Block Header", accordionName, "One", elem.Preformatted(elem.Code(vecty.Text(string(header))))),
		renderCollapsible("Data", accordionName, "Two", transactions),
		renderCollapsible("Evidence", accordionName, "Three", elem.Preformatted(elem.Code(vecty.Text(string(evidence))))),
		renderCollapsible("Last Commit", accordionName, "Four", elem.Preformatted(elem.Code(vecty.Text(string(commit))))),
	)
}

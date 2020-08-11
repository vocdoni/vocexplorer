package components

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/prop"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	tmtypes "github.com/tendermint/tendermint/types"
	"github.com/xeonx/timeago"
	dvotetypes "gitlab.com/vocdoni/go-dvote/types"
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
	return elem.Main(
		renderBlockHeader(len(contents.Block.Data.Txs), contents.Hash, contents.Block.ProposerAddress, contents.Block.Header.Height, contents.Block.Header.Time),
		renderBlockContents(contents.Block),
	)
}

func renderBlockHeader(numTxs int, hash, proposer tmbytes.HexBytes, height int64, tm time.Time) vecty.ComponentOrHTML {
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
						vecty.Text(timeago.English.Format(tm)),
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
				elem.Div(
					elem.Div(
						vecty.Markup(vecty.Class("dt")),
						vecty.Text("Proposer Address"),
					),
					elem.Div(
						vecty.Markup(vecty.Class("dd")),
						vecty.Text(proposer.String()),
					),
				),
			),
		),
	)
}

func renderBlockContents(block *tmtypes.Block) vecty.ComponentOrHTML {
	header, err := json.MarshalIndent(block.Header, "", "    ")
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
	transactions := elem.Preformatted(data...)
	var evidence []byte
	if len(block.Evidence.Evidence) > 0 {
		evidence, err = json.MarshalIndent(block.Evidence, "", "    ")
		util.ErrPrint(err)
	} else {
		evidence = []byte("No evidence")
	}
	commit, err := json.MarshalIndent(block.LastCommit, "", "    ")
	util.ErrPrint(err)
	accordionName := "accordionBlock"
	return elem.Div(
		vecty.Markup(vecty.Class("accordion"), prop.ID(accordionName)),
		renderCollapsible("Block Header", accordionName, "One", elem.Preformatted(vecty.Text(string(header)))),
		renderCollapsible("Data", accordionName, "Two", transactions),
		renderCollapsible("Evidence", accordionName, "Three", elem.Preformatted(vecty.Text(string(evidence)))),
		renderCollapsible("Last Commit", accordionName, "Four", elem.Preformatted(vecty.Text(string(commit)))),
	)
}

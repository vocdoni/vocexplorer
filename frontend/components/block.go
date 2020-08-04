package components

import (
	"encoding/json"
	"time"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/prop"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	tmtypes "github.com/tendermint/tendermint/types"
	"github.com/xeonx/timeago"
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
		renderBlockHeader(len(contents.Block.Data.Txs), contents.Hash, contents.Block.Header.Height, contents.Block.Header.Time),
		renderBlockContents(contents.Block),
	)
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
			),
		),
	)
}

func renderBlockContents(block *tmtypes.Block) vecty.ComponentOrHTML {
	header, err := json.MarshalIndent(block.Header, "", "    ")
	util.ErrPrint(err)
	data := block.Data.StringIndented("    ")
	evidence, err := json.MarshalIndent(block.Evidence, "", "    ")
	util.ErrPrint(err)
	commit, err := json.MarshalIndent(block.LastCommit, "", "    ")
	util.ErrPrint(err)
	accordionName := "accordionBlock"
	return elem.Div(
		vecty.Markup(vecty.Class("accordion"), prop.ID(accordionName)),
		renderCollapsible("Block Header", string(header), accordionName, "One"),
		renderCollapsible("Data", string(data), accordionName, "Two"),
		renderCollapsible("Evidence", string(evidence), accordionName, "Three"),
		renderCollapsible("Last Commit", string(commit), accordionName, "Four"),
	)
}

package components

import (
	"fmt"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/util"
)

//BlockchainInfo is the component to display blockchain information
type BlockchainInfo struct {
	vecty.Core
	header bool
}

//Render renders the BlockchainInfo component
func (b *BlockchainInfo) Render() vecty.ComponentOrHTML {
	syncing := store.Stats.Syncing

	rows := []vecty.MarkupOrChild{vecty.Markup(vecty.Class("stats"))}
	rows = append(rows, row(
		head(vecty.Text("ID")),
		data(vecty.Text(store.Stats.ChainID)),
		head(vecty.Text("Blockchain genesis timestamp")),
		data(vecty.Text(
			fmt.Sprintf(store.Stats.GenesisTimeStamp.Format("Mon Jan _2 15:04:05 UTC 2006")),
		)),
	))
	rows = append(rows, row(
		head(vecty.Text("Block height")),
		data(vecty.Text(
			humanize.Comma(int64(store.Blocks.Count)),
		)),
		head(vecty.Text("Latest block timestamp")),
		data(vecty.Text(
			fmt.Sprintf(time.Unix(int64(store.Stats.BlockTimeStamp), 0).Format("Mon Jan _2 15:04:05 UTC 2006")),
		)),
	))
	rows = append(rows, row())
	rows = append(rows, row(
		head(vecty.Text("Total entities")),
		data(vecty.Text(humanize.Comma(int64(store.Entities.Count)))),
		head(vecty.Text("Total processes")),
		data(vecty.Text(humanize.Comma(int64(store.Processes.Count)))),
	))
	rows = append(rows, row(
		head(vecty.Text("Number of validators")),
		data(vecty.Text(humanize.Comma(int64(store.Stats.ValidatorCount)))),
		head(vecty.Text("Total vote envelopes")),
		data(vecty.Text(humanize.Comma(int64(store.Envelopes.Count)))),
	))
	rows = append(rows, row(
		spanHead(vecty.Text("Sync status")),
		spanData(
			vecty.If(syncing, elem.Span(
				vecty.Markup(vecty.Class("badge", "badge-warning")),
				vecty.Markup(
					vecty.UnsafeHTML("Syncing ("+util.IntToString(store.Blocks.Count+1)+" blocks stored)"),
				),
			)),
			vecty.If(!syncing, &bootstrap.Badge{
				Contents: "In sync",
				Type:     "success",
			}),
		),
	))

	var header *vecty.HTML
	if b.header {
		header = elem.Heading1(vecty.Text("Blockchain information"))
	} else {
		header = elem.Heading2(vecty.Text("Blockchain information"))
	}

	return elem.Section(
		vecty.Markup(vecty.Class("blockchain-info")),
		bootstrap.Card(
			bootstrap.CardParams{
				Body: vecty.List{
					header,
					Container(
						rows...,
					),
				},
			},
		),
	)
}

func row(markup ...vecty.MarkupOrChild) vecty.ComponentOrHTML {
	markup = append(
		markup,
		vecty.Markup(vecty.Class("row")),
	)
	return elem.Div(markup...)
}

func data(markup ...vecty.MarkupOrChild) vecty.ComponentOrHTML {
	markup = append(
		markup,
		vecty.Markup(vecty.Class("col-6", "col-md-3", "data")),
	)
	return elem.Div(markup...)
}

func head(markup ...vecty.MarkupOrChild) vecty.ComponentOrHTML {
	markup = append(
		markup,
		vecty.Markup(vecty.Class("col-6", "col-md-3", "head")),
	)
	return elem.Div(markup...)
}

func spanHead(markup ...vecty.MarkupOrChild) vecty.ComponentOrHTML {
	markup = append(
		markup,
		vecty.Markup(vecty.Class("col-6", "col-md-6", "head")),
	)
	return elem.Div(markup...)
}

func spanData(markup ...vecty.MarkupOrChild) vecty.ComponentOrHTML {
	markup = append(
		markup,
		vecty.Markup(vecty.Class("col-6", "col-md-6", "data")),
	)
	return elem.Div(markup...)
}

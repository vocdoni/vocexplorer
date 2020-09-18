package components

import (
	"time"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/util"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

//BlockchainInfo is the component to display blockchain information
type BlockchainInfo struct {
	vecty.Core
}

//Render renders the BlockchainInfo component
func (b *BlockchainInfo) Render() vecty.ComponentOrHTML {

	if store.Stats.ResultStatus == nil || store.Stats.Genesis == nil {
		return &bootstrap.Alert{
			Type:     "warning",
			Contents: "Waiting for blocks data",
		}
	}

	// Buffer of +- 1 block so syncing does not flash back/forth
	syncing := int(store.Stats.ResultStatus.SyncInfo.LatestBlockHeight)-store.Blocks.Count > 2
	p := message.NewPrinter(language.English)

	return elem.Section(
		vecty.Markup(vecty.Class("blockchain-info")),
		bootstrap.Card(
			bootstrap.CardParams{
				Body: vecty.List{
					elem.Heading4(
						vecty.Text("Blockchain information"),
					),
					Container(
						vecty.Markup(vecty.Class("stats")),
						row(
							head(vecty.Text("ID")),
							data(vecty.Text(store.Stats.Genesis.ChainID)),
							head(vecty.Text("Version")),
							data(vecty.Text(store.Stats.ResultStatus.NodeInfo.Version)),
						),
						row(
							head(vecty.Text("Max block size")),
							data(vecty.Text(p.Sprintf("%d", store.Stats.Genesis.ConsensusParams.Block.MaxBytes))),
							head(vecty.Text("Latest block timestamp")),
							data(vecty.Text(
								p.Sprintf(time.Unix(int64(store.Stats.BlockTimeStamp), 0).Format("Mon Jan _2 15:04:05 UTC 2006")),
							)),
						),
						row(
							head(vecty.Text("Block height")),
							data(vecty.Text(p.Sprintf("%d", store.Stats.ResultStatus.SyncInfo.LatestBlockHeight))),
							head(vecty.Text("Total transactions")),
							data(vecty.Text(p.Sprintf("%d", store.Transactions.Count))),
						),
						row(
							head(vecty.Text("Total entities")),
							data(vecty.Text(p.Sprintf("%d", store.Entities.Count))),
							head(vecty.Text("Total processes")),
							data(vecty.Text(p.Sprintf("%d", store.Processes.Count))),
						),

						row(
							head(vecty.Text("Number of validators")),
							data(vecty.Text(p.Sprintf("%d", len(store.Stats.Genesis.Validators)))),
							head(vecty.Text("Total vote envelopes")),
							data(vecty.Text(p.Sprintf("%d", store.Envelopes.Count))),
						),
						row(
							head(vecty.Text("Average transactions per block")),
							data(vecty.Text(p.Sprintf("%.2f", store.Stats.AvgTxsPerBlock))),

							head(vecty.Text("Max transactions on one block")),
							data(vecty.Text(p.Sprintf("%d", store.Stats.MaxTxsPerBlock))),
						),
						row(
							head(vecty.Text("Average transactions per minute")),
							data(vecty.Text(p.Sprintf("%.2f", store.Stats.AvgTxsPerMinute))),
							head(vecty.Text("Max transactions in one minute")),
							data(vecty.Text(p.Sprintf("%d", store.Stats.MaxTxsPerMinute))),
						),
						row(

							head(vecty.Text("Block with the most transactions")),
							data(
								vecty.Markup(vecty.Class("text-truncate")),
								Link(
									"/block/"+util.IntToString(store.Stats.MaxTxsBlockHeight),
									store.Stats.MaxTxsBlockHash[:util.Min(10, len(store.Stats.MaxTxsBlockHash))]+"...",
									""),
							),

							head(vecty.Text("Minute with the most transactions")),
							data(vecty.Text(p.Sprintf(store.Stats.MaxTxsMinute.Format("Mon Jan _2 15:04 UTC 2006")))),
						),
						row(
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
						),
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

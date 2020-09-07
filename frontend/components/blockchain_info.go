package components

import (
	"time"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
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

	syncing := int(store.Stats.ResultStatus.SyncInfo.LatestBlockHeight)-store.Blocks.Count > 1
	p := message.NewPrinter(language.English)

	return elem.Section(
		vecty.Markup(vecty.Class("blockchain-info")),
		bootstrap.Card(
			bootstrap.CardParams{
				Body: vecty.List{
					elem.Heading4(
						vecty.Text("Blockchain information"),
					),
					elem.Table(
						vecty.Markup(vecty.Class("table")),
						elem.TableRow(
							elem.TableHeader(vecty.Text("ID")),
							elem.TableData(vecty.Text(store.Stats.Genesis.ChainID)),
						),
						elem.TableRow(
							elem.TableHeader(vecty.Text("Version")),
							elem.TableData(vecty.Text(store.Stats.ResultStatus.NodeInfo.Version)),
						),
						elem.TableRow(
							elem.TableHeader(vecty.Text("Max block size")),
							elem.TableData(vecty.Text(
								p.Sprintf("%d", store.Stats.Genesis.ConsensusParams.Block.MaxBytes),
							)),
						),
						elem.TableRow(
							elem.TableHeader(vecty.Text("Latest block timestamp")),
							elem.TableData(vecty.Text(
								p.Sprintf(time.Unix(int64(store.Stats.BlockTimeStamp), 0).Format("Mon Jan _2 15:04:05 UTC 2006")),
							)),
						),
						elem.TableRow(
							elem.TableHeader(vecty.Text("Block height")),
							elem.TableData(vecty.Text(
								p.Sprintf("%d", store.Stats.ResultStatus.SyncInfo.LatestBlockHeight),
							)),
						),
						elem.TableRow(
							elem.TableHeader(vecty.Text("Total transactions")),
							elem.TableData(vecty.Text(
								p.Sprintf("%d", store.Transactions.Count),
							)),
						),
						elem.TableRow(
							elem.TableHeader(vecty.Text("Total entities")),
							elem.TableData(vecty.Text(
								p.Sprintf("%d", store.Entities.Count),
							)),
						),
						elem.TableRow(
							elem.TableHeader(vecty.Text("Total processes")),
							elem.TableData(vecty.Text(
								p.Sprintf("%d", store.Processes.Count),
							)),
						),
						elem.TableRow(
							elem.TableHeader(vecty.Text("Total vote envelopes")),
							elem.TableData(vecty.Text(
								p.Sprintf("%d", store.Envelopes.Count),
							)),
						),
						elem.TableRow(
							elem.TableHeader(vecty.Text("Number of validators")),
							elem.TableData(vecty.Text(
								p.Sprintf("%d", len(store.Stats.Genesis.Validators)),
							)),
						),
						elem.TableRow(
							elem.TableHeader(vecty.Text("Sync status")),
							elem.TableData(
								vecty.If(syncing, elem.Span(
									vecty.Markup(vecty.Class("badge", "badge-warning")),
									vecty.Markup(
										vecty.UnsafeHTML("Syncing ("+util.IntToString(store.Blocks.Count)+" blocks stored)"),
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

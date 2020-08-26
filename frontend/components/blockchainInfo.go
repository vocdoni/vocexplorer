package components

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/rpc"
	"gitlab.com/vocdoni/vocexplorer/util"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type BlockchainInfo struct {
	vecty.Core
	T *rpc.TendermintInfo
}

func (b *BlockchainInfo) Render() vecty.ComponentOrHTML {

	if b.T.ResultStatus == nil {
		return &bootstrap.Alert{
			Type:     "warning",
			Contents: "Waiting for blocks data",
		}
	}

	syncing := int(b.T.ResultStatus.SyncInfo.LatestBlockHeight)-b.T.TotalBlocks > 1
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
							elem.TableData(vecty.Text(b.T.Genesis.ChainID)),
						),
						elem.TableRow(
							elem.TableHeader(vecty.Text("Version")),
							elem.TableData(vecty.Text(b.T.ResultStatus.NodeInfo.Version)),
						),
						elem.TableRow(
							elem.TableHeader(vecty.Text("Block Height")),
							elem.TableData(vecty.Text(
								p.Sprintf("%d", store.CurrentBlockHeight),
							)),
						),
						elem.TableRow(
							elem.TableHeader(vecty.Text("Max block size")),
							elem.TableData(vecty.Text(
								p.Sprintf("%d", b.T.Genesis.ConsensusParams.Block.MaxBytes),
							)),
						),
						elem.TableRow(
							elem.TableHeader(vecty.Text("Total transactions")),
							elem.TableData(vecty.Text(
								p.Sprintf("%d", b.T.TotalTxs),
							)),
						),
						elem.TableRow(
							elem.TableHeader(vecty.Text("Total entities")),
							elem.TableData(vecty.Text(
								p.Sprintf("%d", b.T.TotalEntities),
							)),
						),
						elem.TableRow(
							elem.TableHeader(vecty.Text("Total processes")),
							elem.TableData(vecty.Text(
								p.Sprintf("%d", b.T.TotalProcesses),
							)),
						),
						elem.TableRow(
							elem.TableHeader(vecty.Text("Total vote envelopes")),
							elem.TableData(vecty.Text(
								p.Sprintf("%d", b.T.TotalEnvelopes),
							)),
						),
						elem.TableRow(
							elem.TableHeader(vecty.Text("Number of validators")),
							elem.TableData(vecty.Text(
								p.Sprintf("%d", len(b.T.Genesis.Validators)),
							)),
						),
						elem.TableRow(
							elem.TableHeader(vecty.Text("Sync status")),
							elem.TableData(
								// This badge component does not rerender when it should. Not sure why
								// vecty.If(syncing, &bootstrap.Badge{
								// 	Contents: p.Sprintf("Syncing (%d blocks stored)", b.T.TotalBlocks),
								// 	Type:     "warning",
								// }),
								vecty.If(syncing, elem.Span(
									vecty.Markup(vecty.Class("badge", "badge-warning")),
									vecty.Markup(
										vecty.UnsafeHTML("Syncing ("+util.IntToString(b.T.TotalBlocks)+" blocks stored)"),
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

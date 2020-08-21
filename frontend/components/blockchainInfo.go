package components

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/rpc"
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
								p.Sprintf("%d", b.T.ResultStatus.SyncInfo.LatestBlockHeight),
							)),
						),
						elem.TableRow(
							elem.TableHeader(vecty.Text("Total transactions")),
							elem.TableData(vecty.Text(
								p.Sprintf("%d", b.T.TotalTxs),
							)),
						),
						elem.TableRow(
							elem.TableHeader(vecty.Text("Sync status")),
							elem.TableData(
								vecty.If(syncing, &bootstrap.Badge{
									Contents: p.Sprintf("Syncing (%d blocks stored)", +b.T.TotalBlocks),
									Type:     "warning",
								}),
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

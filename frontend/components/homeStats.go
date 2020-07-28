package components

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/xeonx/timeago"
	"gitlab.com/vocdoni/vocexplorer/client"
	"gitlab.com/vocdoni/vocexplorer/rpc"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// StatsView renders the stats pane
type StatsView struct {
	vecty.Core
	t  *rpc.TendermintInfo
	vc *client.VochainInfo
}

// Render renders the StatsView component
func (b *StatsView) Render() vecty.ComponentOrHTML {
	if b.t != nil && b.vc != nil {
		return elem.Section(
			renderBlockchainStats(b.t, b.vc),
			renderRecentBlocks(b.t),
			// renderTimeStats(b.t),
			// renderStatus(b.t),
		)
	}
	return elem.Div(vecty.Text("Waiting for blockchain statistics..."))
}

func renderBlockchainStats(t *rpc.TendermintInfo, vc *client.VochainInfo) vecty.ComponentOrHTML {
	if vc.BlockTime != nil && t.ResultStatus != nil {
		return elem.Div(
			vecty.Markup(vecty.Class("bc-stats")),
			elem.Div(
				elem.Div(vecty.Markup(vecty.Class("card-col-3")),
					vecty.If(vc.BlockTime != nil, elem.Table(
						elem.Caption(elem.Heading2(vecty.Text("Average block times: "))),
						elem.TableHead(
							elem.TableRow(elem.TableHeader(vecty.Text("Time period")), elem.TableHeader(vecty.Text("Avg time"))),
						),
						elem.TableBody(
							vecty.If(vc.BlockTime[0] > 0, elem.TableRow(elem.TableData(vecty.Text("Last 1m")), elem.TableData(vecty.Text(util.MsToString(vc.BlockTime[0]))))),
							vecty.If(vc.BlockTime[1] > 0, elem.TableRow(elem.TableData(vecty.Text("Last 10m")), elem.TableData(vecty.Text(util.MsToString(vc.BlockTime[1]))))),
							vecty.If(vc.BlockTime[2] > 0, elem.TableRow(elem.TableData(vecty.Text("Last 1h")), elem.TableData(vecty.Text(util.MsToString(vc.BlockTime[2]))))),
							vecty.If(vc.BlockTime[3] > 0, elem.TableRow(elem.TableData(vecty.Text("Last 6h")), elem.TableData(vecty.Text(util.MsToString(vc.BlockTime[3]))))),
							vecty.If(vc.BlockTime[4] > 0, elem.TableRow(elem.TableData(vecty.Text("Last 24h")), elem.TableData(vecty.Text(util.MsToString(vc.BlockTime[4]))))),
						),
					)),
				),
				elem.Div(vecty.Markup(vecty.Class("card-col-3")),
					vecty.Text("Current Block Height: "+util.IntToString(t.ResultStatus.SyncInfo.LatestBlockHeight)),
				),
				elem.Div(vecty.Markup(vecty.Class("card-col-3")),
					vecty.Text("Total Txs: "+util.IntToString(t.TxCount)),
				),
			),
			elem.Div(
				vecty.Markup(vecty.Class("card-col-1")),
				elem.Table(
					elem.TableHead(
						elem.TableRow(elem.TableHeader(vecty.Text("Chain ID")), elem.TableHeader(vecty.Text("App Version")), elem.TableHeader(vecty.Text("Max Block Size")), elem.TableHeader(vecty.Text("Num Validators"))),
					),
					elem.TableBody(
						elem.TableRow(
							elem.TableData(vecty.Text(t.Genesis.ChainID)),
							elem.TableData(vecty.Text(t.ResultStatus.NodeInfo.Version)),
							elem.TableData(vecty.Text(util.IntToString(t.Genesis.ConsensusParams.Block.MaxBytes))),
							elem.TableData(vecty.Text(util.IntToString(len(t.Genesis.Validators)))),
						),
					),
				),
			),
		)
	}
	return elem.Div(vecty.Text("Waiting for transaction & block data..."))
}

func renderTimeStats(t *rpc.TendermintInfo) vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(vecty.Class("tx-stats")),
		vecty.Text("Txs/hr: "),
		vecty.Text("Txs/day: "),
	)
}

func renderRecentBlocks(t *rpc.TendermintInfo) vecty.ComponentOrHTML {
	if t.RecentBlocks != nil {
		return elem.Div(
			vecty.Markup(vecty.Class("recent-blocks")),
			elem.Heading3(
				vecty.Text("Blocks"),
			),
			elem.Div(
				vecty.Markup(vecty.Class("card-deck")),
				renderBlock(t.RecentBlocks, t.RecentBlockResults, 3),
				renderBlock(t.RecentBlocks, t.RecentBlockResults, 2),
				renderBlock(t.RecentBlocks, t.RecentBlockResults, 1),
				renderBlock(t.RecentBlocks, t.RecentBlockResults, 0),
			),
		)
	}
	return elem.Div(vecty.Text("No updated list of recent blocks"))
}

func renderBlock(recentBlocks []coretypes.ResultBlock, recentBlockResults []coretypes.ResultBlockResults, index int) vecty.ComponentOrHTML {
	if len(recentBlocks) > index && len(recentBlockResults) > index {
		return elem.Div(vecty.Markup(vecty.Class("card")),
			elem.Div(
				vecty.Markup(vecty.Class("card-header")),
				vecty.Text(util.IntToString(recentBlocks[index].Block.Header.Height)),
			),
			elem.Div(
				vecty.Markup(vecty.Class("card-body")),
				elem.Div(
					vecty.Markup(vecty.Class("block-card-heading")),
					elem.Div(
						vecty.Text(util.IntToString(len(recentBlockResults[index].TxsResults))+" transactions"),
					),
					elem.Div(
						vecty.Text(timeago.English.Format(recentBlocks[index].Block.Header.Time)),
					),
				),
				elem.Div(
					elem.Div(
						vecty.Markup(vecty.Class("dt")),
						vecty.Text("Hash"),
					),
					elem.Div(
						vecty.Markup(vecty.Class("dd")),
						vecty.Text(recentBlocks[index].BlockID.Hash.String()),
					),
				),
			),
		)
	}
	return vecty.Text("No block available ")
}

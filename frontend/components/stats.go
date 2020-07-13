package components

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
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
							elem.TableRow(elem.TableData(vecty.Text("Last 1m")), elem.TableData(vecty.Text(util.MsToString(vc.BlockTime[0])))),
							elem.TableRow(elem.TableData(vecty.Text("Last 10m")), elem.TableData(vecty.Text(util.MsToString(vc.BlockTime[1])))),
							elem.TableRow(elem.TableData(vecty.Text("Last 1h")), elem.TableData(vecty.Text(util.MsToString(vc.BlockTime[2])))),
							elem.TableRow(elem.TableData(vecty.Text("Last 6h")), elem.TableData(vecty.Text(util.MsToString(vc.BlockTime[3])))),
							elem.TableRow(elem.TableData(vecty.Text("Last 24h")), elem.TableData(vecty.Text(util.MsToString(vc.BlockTime[4])))),
						),
					)),
				),
				elem.Div(vecty.Markup(vecty.Class("card-col-3")),
					vecty.Text("Current Block Num: "+util.IntToString(t.ResultStatus.SyncInfo.LatestBlockHeight)),
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
			vecty.If(len(*t.RecentBlocks) >= 0, renderBlock(&(*t.RecentBlocks)[0])),
			vecty.If(len(*t.RecentBlocks) >= 1, renderBlock(&(*t.RecentBlocks)[1])),
			vecty.If(len(*t.RecentBlocks) >= 2, renderBlock(&(*t.RecentBlocks)[2])),
			vecty.If(len(*t.RecentBlocks) >= 3, renderBlock(&(*t.RecentBlocks)[3])),
			elem.Div(
				vecty.Markup(vecty.Class("card-col-3")),
				vecty.Text("Total Txs: "+util.IntToString(t.TxCount)),
			),
		)
	}
	return elem.Div(vecty.Text("No updated list of recent blocks"))
}

func renderBlock(result *coretypes.ResultBlock) vecty.ComponentOrHTML {
	return elem.Div(vecty.Markup(vecty.Class("card-col-3")),
		vecty.Text("Block Num: "+util.IntToString(result.Block.Header.Height)),
		vecty.Text("Block Hash: "+result.BlockID.Hash.String()),
		// vecty.Text("Total Transactions: "+util.IntToString(t.RecentBlocks[0].Block.Header.Height)),
	)
}

// func renderStatus(t *rpc.TendermintInfo) vecty.ComponentOrHTML {
// 	if t.ResultStatus != nil {
// 		sync := t.ResultStatus.SyncInfo
// 		return elem.Div(
// 			elem.Div(vecty.Text("Latest Block Hash: "+sync.LatestBlockHash.String())),
// 			elem.Div(vecty.Text("Latest App Hash: "+sync.LatestAppHash.String())),
// 			elem.Div(vecty.Text("Latest Block Height: "+strconv.Itoa(int(sync.LatestBlockHeight)))),
// 			elem.Div(vecty.Text("Latest Block Time: "+sync.LatestBlockTime.String())),
// 		)
// 	}
// 	return elem.Div(vecty.Text("Waiting for tendermint blockchain info..."))
// }

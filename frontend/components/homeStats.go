package components

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/gopherjs/vecty/prop"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	"gitlab.com/vocdoni/vocexplorer/client"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/rpc"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// StatsView renders the stats pane
type StatsView struct {
	vecty.Core
	t          *rpc.TendermintInfo
	vc         *client.VochainInfo
	blockIndex int
	refreshCh  chan int
}

// Render renders the StatsView component
func (b *StatsView) Render() vecty.ComponentOrHTML {
	if b.t != nil && b.vc != nil {
		return elem.Section(
			renderBlockchainStats(b.t, b.vc),
			renderBlockList(b),
		)
	}
	return elem.Div(vecty.Text("Waiting for blockchain statistics..."))
}

func renderBlockList(b *StatsView) vecty.ComponentOrHTML {
	if b.t != nil && b.t.ResultStatus != nil {
		return elem.Div(
			vecty.Text("Page "+util.IntToString(b.blockIndex/10+1)),
			elem.Button(
				vecty.Text("back to top"),
				vecty.Markup(
					event.Click(func(e *vecty.Event) {
						b.blockIndex = 0
						b.refreshCh <- b.blockIndex
						vecty.Rerender(b)
					}),
					vecty.MarkupIf(
						b.blockIndex != 0,
						prop.Disabled(false),
					),
					vecty.MarkupIf(
						b.blockIndex == 0,
						prop.Disabled(true),
					),
				),
			),
			elem.Button(
				vecty.Text("prev"),
				vecty.Markup(
					event.Click(func(e *vecty.Event) {
						b.blockIndex = util.Max(b.blockIndex-config.SearchPageSmall, 0)
						b.refreshCh <- b.blockIndex
						vecty.Rerender(b)
					}),
					vecty.MarkupIf(
						b.blockIndex > 0,
						prop.Disabled(false),
					),
					vecty.MarkupIf(
						b.blockIndex < 1,
						prop.Disabled(true),
					),
				),
			),
			elem.Button(vecty.Text("next"),
				vecty.Markup(
					event.Click(func(e *vecty.Event) {
						b.blockIndex = util.Min(b.blockIndex+config.SearchPageSmall, int(b.t.ResultStatus.SyncInfo.LatestBlockHeight))
						b.refreshCh <- b.blockIndex
						vecty.Rerender(b)
					}),
					vecty.MarkupIf(
						b.blockIndex < int(b.t.ResultStatus.SyncInfo.LatestBlockHeight),
						prop.Disabled(false),
					),
					vecty.MarkupIf(
						b.blockIndex >= int(b.t.ResultStatus.SyncInfo.LatestBlockHeight),
						prop.Disabled(true),
					),
				),
			),
			elem.Heading4(vecty.Text("Block ID list: ")),
			// vecty.If(len(b.vc.BlockSearchList) < b.numBlocks, vecty.Text("Loading block info...")),
			elem.UnorderedList(
				renderBlocks(b.t, b.blockIndex)...,
			),
		)
	}
	return elem.Div(vecty.Text("Waiting for blockchain info..."))
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

func renderBlocks(t *rpc.TendermintInfo, index int) []vecty.MarkupOrChild {
	var blockList []vecty.MarkupOrChild
	if t.BlockList[config.SearchPageSmall-1].Block == nil {
		return []vecty.MarkupOrChild{elem.Div(vecty.Text("No blocks available"))}
	}
	if t.BlockList[config.SearchPageSmall-1].Block.Height != t.ResultStatus.SyncInfo.LatestBlockHeight-1-int64(index) {
		return []vecty.MarkupOrChild{elem.Div(vecty.Text("Loading blocks........"))}
	}
	for i := config.SearchPageSmall - 1; i >= 0; i-- {
		block := t.BlockList[i]
		// for i, block := range t.BlockList {
		blockList = append(blockList, renderBlock(block, t.BlockListResults[i]))
	}
	return blockList
}

func renderBlock(block coretypes.ResultBlock, blockResults coretypes.ResultBlockResults) vecty.ComponentOrHTML {
	return elem.ListItem(
		vecty.Markup(vecty.Class("card-col-3")),
		vecty.Text("Block Num: "+util.IntToString(block.Block.Header.Height)),
		vecty.Text(" Block Hash: "+block.BlockID.Hash.String()),
		// vecty.Text(" Total Transactions: "+util.IntToString(len(blockResults.TxsResults))),
	)
}

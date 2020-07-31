package components

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/xeonx/timeago"
	"gitlab.com/vocdoni/vocexplorer/client"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/rpc"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// StatsView renders the stats pane
type StatsView struct {
	vecty.Core
	t           *rpc.TendermintInfo
	vc          *client.VochainInfo
	currentPage int
	refreshCh   chan int
}

// Render renders the StatsView component
func (b *StatsView) Render() vecty.ComponentOrHTML {
	if b.t != nil && b.vc != nil {
		return elem.Section(
			&Jumbotron{
				vc: b.vc,
				t:  b.t,
			},
			elem.Div(
				vecty.Markup(vecty.Class("container")),
				renderBlockchainStats(b.t, b.vc),
				renderBlockList(b),
			),
		)
	}
	return elem.Div(vecty.Text("Waiting for blockchain statistics..."))
}

func renderBlockList(b *StatsView) vecty.ComponentOrHTML {
	if b.t != nil && b.t.ResultStatus != nil {
		p := &Pagination{
			TotalPages:      int(b.t.ResultStatus.SyncInfo.LatestBlockHeight-1) / config.HomeWidgetBlocksListSize,
			TotalItems:      &b.t.TotalBlocks,
			CurrentPage:     &b.currentPage,
			RefreshCh:       b.refreshCh,
			ListSize:        config.HomeWidgetBlocksListSize,
			RenderSearchBar: true,
		}
		p.RenderFunc = func(index int) vecty.ComponentOrHTML {
			return renderBlocks(p, b.t, index)
		}
		return elem.Div(
			vecty.Markup(vecty.Class("recent-blocks")),
			elem.Heading3(
				vecty.Text("Blocks"),
			),
			p,
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
					elem.Table(
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
					),
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

func renderBlocks(p *Pagination, t *rpc.TendermintInfo, index int) vecty.ComponentOrHTML {
	var blockList []vecty.MarkupOrChild
	if t.BlockList[p.ListSize-1].Block == nil {
		return elem.Div(vecty.Text("No blocks available"))
	}
	// if t.BlockList[p.ListSize-1].Block.Height != t.ResultStatus.SyncInfo.LatestBlockHeight-1-int64(index*p.ListSize) {
	// 	return elem.Div(vecty.Text("Loading blocks........"))
	// }
	for i := p.ListSize - 1; i >= 0; i-- {
		block := t.BlockList[i]
		// for i, block := range t.BlockList {
		blockList = append(blockList, renderBlock(block, t.BlockListResults[i]))
	}
	blockList = append(blockList, vecty.Markup(vecty.Class("responsive-card-deck")))
	return elem.Div(
		blockList...,
	)
}

func renderBlock(block coretypes.ResultBlock, blockResults coretypes.ResultBlockResults) vecty.ComponentOrHTML {
	return elem.Div(vecty.Markup(vecty.Class("card-deck-col")),
		elem.Div(vecty.Markup(vecty.Class("card")),
			elem.Div(
				vecty.Markup(vecty.Class("card-header")),
				vecty.Text(util.IntToString(block.Block.Header.Height)),
			),
			elem.Div(
				vecty.Markup(vecty.Class("card-body")),
				elem.Div(
					vecty.Markup(vecty.Class("block-card-heading")),
					elem.Div(
						vecty.Text(util.IntToString(len(blockResults.TxsResults))+" transactions"),
					),
					elem.Div(
						vecty.Text(timeago.English.Format(block.Block.Header.Time)),
					),
				),
				elem.Div(
					elem.Div(
						vecty.Markup(vecty.Class("dt")),
						vecty.Text("Hash"),
					),
					elem.Div(
						vecty.Markup(vecty.Class("dd")),
						vecty.Text(block.BlockID.Hash.String()),
					),
				),
			),
		),
	)
}

// elem.Navigation(
// 	elem.Span(
// 		vecty.Text("Page "+util.IntToString(b.blockIndex/10+1)),
// 	),
// 	elem.UnorderedList(
// 		vecty.Markup(vecty.Class("pagination")),
// 		elem.ListItem(
// 			vecty.Markup(
// 				vecty.MarkupIf(
// 					b.blockIndex != 0,
// 					vecty.Class("page-item"),
// 				),
// 				vecty.MarkupIf(
// 					b.blockIndex == 0,
// 					vecty.Class("page-item", "disabled"),
// 				),
// 			),
// 			elem.Button(
// 				vecty.Markup(
// 					vecty.Class("page-link"),
// 					event.Click(func(e *vecty.Event) {
// 						b.blockIndex = 0
// 						b.refreshCh <- b.blockIndex
// 						vecty.Rerender(b)
// 					}),
// 					vecty.MarkupIf(
// 						b.blockIndex != 0,
// 						prop.Disabled(false),
// 					),
// 					vecty.MarkupIf(
// 						b.blockIndex == 0,
// 						prop.Disabled(true),
// 					),
// 				),
// 				elem.Span(
// 					vecty.Text("«"),
// 				),
// 				elem.Span(
// 					vecty.Markup(vecty.Class("sr-only")),
// 					vecty.Text("Back to top"),
// 				),
// 			),
// 		),
// 		elem.ListItem(
// 			vecty.Markup(
// 				vecty.MarkupIf(
// 					b.blockIndex > 0,
// 					vecty.Class("page-item"),
// 				),
// 				vecty.MarkupIf(
// 					b.blockIndex <= 0,
// 					vecty.Class("page-item", "disabled"),
// 				),
// 			),
// 			elem.Button(
// 				vecty.Text("prev"),
// 				vecty.Markup(
// 					vecty.Class("page-link"),
// 					event.Click(func(e *vecty.Event) {
// 						b.blockIndex = util.Max(b.blockIndex-config.HomeWidgetBlocksListSize, 0)
// 						b.refreshCh <- b.blockIndex
// 						vecty.Rerender(b)
// 					}),
// 					vecty.MarkupIf(
// 						b.blockIndex > 0,
// 						prop.Disabled(false),
// 					),
// 					vecty.MarkupIf(
// 						b.blockIndex < 1,
// 						prop.Disabled(true),
// 					),
// 				),
// 			),
// 		),
// 		elem.ListItem(
// 			vecty.Markup(
// 				vecty.MarkupIf(
// 					b.blockIndex < int(b.t.ResultStatus.SyncInfo.LatestBlockHeight),
// 					vecty.Class("page-item"),
// 				),
// 				vecty.MarkupIf(
// 					b.blockIndex >= int(b.t.ResultStatus.SyncInfo.LatestBlockHeight),
// 					vecty.Class("page-item", "disabled"),
// 				),
// 			),
// 			elem.Button(vecty.Text("next"),
// 				vecty.Markup(
// 					vecty.Class("page-link"),
// 					event.Click(func(e *vecty.Event) {
// 						b.blockIndex = util.Min(b.blockIndex+config.HomeWidgetBlocksListSize, int(b.t.ResultStatus.SyncInfo.LatestBlockHeight))
// 						b.refreshCh <- b.blockIndex
// 						vecty.Rerender(b)
// 					}),
// 					vecty.MarkupIf(
// 						b.blockIndex < int(b.t.ResultStatus.SyncInfo.LatestBlockHeight),
// 						prop.Disabled(false),
// 					),
// 					vecty.MarkupIf(
// 						b.blockIndex >= int(b.t.ResultStatus.SyncInfo.LatestBlockHeight),
// 						prop.Disabled(true),
// 					),
// 				),
// 			),
// 		),
// 	),
// ),
// renderBlocks(b.t, b.blockIndex),

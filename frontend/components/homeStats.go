package components

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"gitlab.com/vocdoni/vocexplorer/client"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/rpc"
)

// StatsView renders the stats pane
type StatsView struct {
	vecty.Core
	t           *rpc.TendermintInfo
	vc          *client.VochainInfo
	currentPage int
	refreshCh   chan int
	gwClient    *client.Client
}

// Render renders the StatsView component
func (b *StatsView) Render() vecty.ComponentOrHTML {
	if b.t != nil && b.vc != nil {
		return elem.Section(
			vecty.Markup(
				event.BeforeUnload(func(i *vecty.Event) {
					b.gwClient.Close()
				}),
			),
			&Jumbotron{
				vc: b.vc,
				t:  b.t,
			},
			Container(
				&LatestBlocksWidget{
					T: b.t,
				},
				&BlockchainInfo{
					T: b.t,
				},
				&AverageBlockTimes{
					VC: b.vc,
				},
			),
		)
	}
	return elem.Div(vecty.Text("Waiting for blockchain statistics..."))
}

func renderHomeBlockList(b *StatsView) vecty.ComponentOrHTML {
	if b.t != nil && b.t.ResultStatus != nil {
		p := &Pagination{
			TotalPages:      int(b.t.TotalBlocks) / config.HomeWidgetBlocksListSize,
			TotalItems:      &b.t.TotalBlocks,
			CurrentPage:     &b.currentPage,
			RefreshCh:       b.refreshCh,
			ListSize:        config.HomeWidgetBlocksListSize,
			RenderSearchBar: false,
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

func renderTimeStats(t *rpc.TendermintInfo) vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(vecty.Class("tx-stats")),
		vecty.Text("Txs/hr: "),
		vecty.Text("Txs/day: "),
	)
}

package components

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/client"
	"gitlab.com/vocdoni/vocexplorer/rpc"
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

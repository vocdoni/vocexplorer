package components

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
)

// StatsView renders the stats pane
type StatsView struct {
	vecty.Core
}

// Render renders the StatsView component
func (b *StatsView) Render() vecty.ComponentOrHTML {
	if store.Stats == nil {
		return LoadingBar()
	}
	return elem.Section(
		&Jumbotron{},
		Container(
			elem.Div(
				vecty.Markup(vecty.Class("dash-heading")),
				elem.Heading1(
					vecty.Text("Vochain Explorer: "+store.Stats.ChainID),
				),
			),
			&LatestBlocksWidget{},
			&BlockchainInfo{
				header: false,
			},
			&AverageBlockTimes{},
		),
	)
}

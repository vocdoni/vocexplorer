package components

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

// StatsView renders the stats pane
type StatsView struct {
	vecty.Core
}

// Render renders the StatsView component
func (b *StatsView) Render() vecty.ComponentOrHTML {
	// if store.Stats.ResultStatus != nil || len(store.Stats.APIList) > 0 || store.Stats.Genesis != nil {
	return elem.Section(
		&Jumbotron{},
		Container(
			&LatestBlocksWidget{},
			&BlockchainInfo{},
			&AverageBlockTimes{},
		),
	)
	// }
	// return elem.Div(vecty.Text("Waiting for blockchain statistics..."))
}

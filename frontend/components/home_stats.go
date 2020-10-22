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
	return elem.Section(
		&Jumbotron{},
		Container(
			&LatestBlocksWidget{},
			&BlockchainInfo{},
			&AverageBlockTimes{},
		),
	)
}

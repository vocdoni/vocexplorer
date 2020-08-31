package pages

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/frontend/components"
)

// Stats is a pretty page for all our blockchain statistics
type Stats struct {
	vecty.Core
	Cfg *config.Cfg
}

// Render renders the Stats component
func (stats *Stats) Render() vecty.ComponentOrHTML {
	return stats.Component()
}

// Component returns the stats component
func (stats *Stats) Component() vecty.ComponentOrHTML {
	return components.Container(
		elem.Section(
			bootstrap.Card(bootstrap.CardParams{
				Body: vecty.List{
					elem.Heading3(
						vecty.Text("Stats"),
					),
					&components.BlockchainInfo{},
					vecty.Text("Most stats will be moved here :)"),
				},
			}),
		),
	)
}

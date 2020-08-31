package pages

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/frontend/components"
)

// BlocksView is a pretty page for all our blockchain statistics
type BlocksView struct {
	vecty.Core
	Cfg *config.Cfg
}

// Render renders the Blocks component
func (stats *BlocksView) Render() vecty.ComponentOrHTML {
	return stats.Component()
}

// Component generates the actual BlocksView component
func (stats *BlocksView) Component() vecty.ComponentOrHTML {
	return components.Container(
		elem.Section(
			bootstrap.Card(bootstrap.CardParams{
				Body: vecty.List{
					elem.Heading3(
						vecty.Text("Blocks"),
					),
					vecty.Text("Blocks list will be here"),
				},
			}),
		),
	)
}

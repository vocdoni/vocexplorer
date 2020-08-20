package pages

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/components"
)

// Stats is a pretty page for all our blockchain statistics
type Stats struct {
	vecty.Core
	Cfg *config.Cfg
}

// Render renders the Stats component
func (stats *Stats) Render() vecty.ComponentOrHTML {
	return components.Container(
		elem.Section(
			elem.Heading3(
				vecty.Text("Stats"),
			),
			vecty.Text("Most stats will be moved here :)"),
		),
	)
}

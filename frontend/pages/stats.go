package pages

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/components"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
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
	dash := new(components.StatsDashboardView)
	dash.Rendered = false
	// Ensure component rerender is only triggered once component has been rendered
	if !store.Listeners.Has(dash) {
		store.Listeners.Add(dash, func() {
			if dash.Rendered {
				vecty.Rerender(dash)
			}
		})
	}
	go components.UpdateAndRenderStatsDashboard(dash)
	return elem.Div(dash)
}

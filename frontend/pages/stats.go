package pages

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/vocdoni/vocexplorer/config"
	"github.com/vocdoni/vocexplorer/frontend/actions"
	"github.com/vocdoni/vocexplorer/frontend/components"
	"github.com/vocdoni/vocexplorer/frontend/dispatcher"
	"github.com/vocdoni/vocexplorer/frontend/store"
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
	dispatcher.Dispatch(&actions.SetCurrentPage{Page: "stats"})
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
	go components.UpdateStatsDashboard(dash)
	return elem.Div(
		&components.Header{},
		dash,
	)
}

package pages

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/components"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
)

// ParticipationView renders the processes page
type ParticipationView struct {
	vecty.Core
	Cfg *config.Cfg
}

// Render renders the ParticipationView component
func (home *ParticipationView) Render() vecty.ComponentOrHTML {
	dash := new(components.ParticipationDashboardView)
	dash.Rendered = false
	// Ensure component rerender is only triggered once component has been rendered
	if !store.Listeners.Has(dash) {
		store.Listeners.Add(dash, func() {
			if dash.Rendered {
				vecty.Rerender(dash)
			}
		})
	}
	go components.UpdateAndRenderParticipationDashboard(dash)
	return elem.Div(dash)
}

package pages

import (
	"github.com/gopherjs/vecty"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/components"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
)

// HomeView renders the Home landing page
type HomeView struct {
	vecty.Core
	Cfg *config.Cfg
}

// Render renders the HomeView component
func (home *HomeView) Render() vecty.ComponentOrHTML {
	dash := new(components.DashboardView)
	dash.Rendered = false
	// Ensure component rerender is only triggered once component has been rendered
	store.Listeners.Add(dash, func() {
		if dash.Rendered {
			vecty.Rerender(dash)
		}
	})
	go components.UpdateAndRenderHomeDashboard(dash)
	return dash
}

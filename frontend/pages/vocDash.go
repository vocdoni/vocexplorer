package pages

import (
	"github.com/gopherjs/vecty"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/components"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
)

// VocDashView renders the processes page
type VocDashView struct {
	vecty.Core
	Cfg *config.Cfg
}

// Render renders the VocDashView component
func (home *VocDashView) Render() vecty.ComponentOrHTML {
	dash := new(components.VocDashDashboardView)
	dash.Rendered = false
	// Ensure component rerender is only triggered once component has been rendered
	store.Listeners.Add(dash, func() {
		if dash.Rendered {
			vecty.Rerender(dash)
		}
	})
	go components.UpdateAndRenderVocDashDashboard(dash)
	return dash
}

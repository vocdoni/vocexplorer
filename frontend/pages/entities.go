package pages

import (
	"github.com/gopherjs/vecty"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/components"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	router "marwan.io/vecty-router"
)

// EntitiesView renders the Entities page
type EntitiesView struct {
	vecty.Core
	Cfg *config.Cfg
}

// Render renders the EntitiesView component
func (home *EntitiesView) Render() vecty.ComponentOrHTML {
	dash := new(components.EntitiesDashboardView)
	dispatcher.Dispatch(&actions.SetCurrentEntityID{EntityID: router.GetNamedVar(home)["id"]})
	store.Listeners.Add(dash, func() {
		vecty.Rerender(dash)
	})
	go components.UpdateAndRenderEntitiesDashboard(dash)
	return dash
}

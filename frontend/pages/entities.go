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

// EntitiesView renders the processes page
type EntitiesView struct {
	vecty.Core
	Cfg *config.Cfg
}

// Render renders the EntitiesView component
func (home *EntitiesView) Render() vecty.ComponentOrHTML {
	dispatcher.Dispatch(&actions.SetCurrentPage{Page: "entities"})
	dash := new(components.EntitiesDashboardView)
	dash.Rendered = false
	// Ensure component rerender is only triggered once component has been rendered
	if !store.Listeners.Has(dash) {
		store.Listeners.Add(dash, func() {
			if dash.Rendered {
				vecty.Rerender(dash)
			}
		})
	}
	go components.UpdateEntitiesDashboard(dash)
	return elem.Div(
		&components.Header{},
		dash,
	)
}

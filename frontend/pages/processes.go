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

// ProcessesView renders the processes page
type ProcessesView struct {
	vecty.Core
	Cfg *config.Cfg
}

// Render renders the ProcessesView component
func (home *ProcessesView) Render() vecty.ComponentOrHTML {
	dash := new(components.ProcessesDashboardView)
	dispatcher.Dispatch(&actions.SetCurrentProcessID{ID: router.GetNamedVar(home)["id"]})
	dash.Rendered = false
	// Ensure component rerender is only triggered once component has been rendered
	store.Listeners.Add(dash, func() {
		if dash.Rendered {
			vecty.Rerender(dash)
		}
	})
	go components.UpdateAndRenderProcessesDashboard(dash)
	return dash
}

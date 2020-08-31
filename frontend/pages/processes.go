package pages

import (
	"github.com/gopherjs/vecty"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/components"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
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
	go components.UpdateAndRenderProcessesDashboard(dash)
	return dash
}

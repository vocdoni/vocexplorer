package pages

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/components"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
)

// ProcessesView renders the processes page
type ProcessesView struct {
	vecty.Core
	Cfg *config.Cfg
}

// Render renders the ProcessesView component
func (home *ProcessesView) Render() vecty.ComponentOrHTML {
	dash := new(components.ProcessesDashboardView)
	dash.Rendered = false
	// Ensure component rerender is only triggered once component has been rendered
	if !store.Listeners.Has(dash) {
		store.Listeners.Add(dash, func() {
			if dash.Rendered {
				vecty.Rerender(dash)
			}
		})
	}
	go components.UpdateProcessesDashboard(dash)
	return elem.Div(dash)
}

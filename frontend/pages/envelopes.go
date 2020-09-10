package pages

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/components"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
)

// EnvelopesView renders the processes page
type EnvelopesView struct {
	vecty.Core
	Cfg *config.Cfg
}

// Render renders the EnvelopesView component
func (home *EnvelopesView) Render() vecty.ComponentOrHTML {
	dash := new(components.EnvelopesDashboardView)
	dash.Rendered = false
	// Ensure component rerender is only triggered once component has been rendered
	if !store.Listeners.Has(dash) {
		store.Listeners.Add(dash, func() {
			if dash.Rendered {
				vecty.Rerender(dash)
			}
		})
	}
	go components.UpdateEnvelopesDashboard(dash)
	return elem.Div(
		&components.Header{},
		dash,
	)
}

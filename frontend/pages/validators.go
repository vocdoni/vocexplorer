package pages

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/frontend/components"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
)

// ValidatorsView renders the Validators page
type ValidatorsView struct {
	vecty.Core
}

// Render renders the ValidatorsView component
func (home *ValidatorsView) Render() vecty.ComponentOrHTML {
	dash := new(components.ValidatorsDashboardView)
	dash.Rendered = false
	// Ensure component rerender is only triggered once component has been rendered
	if !store.Listeners.Has(dash) {
		store.Listeners.Add(dash, func() {
			if dash.Rendered {
				vecty.Rerender(dash)
			}
		})
	}
	go components.UpdateValidatorsDashboard(dash)
	return elem.Div(
		&components.Header{},
		dash,
	)
}

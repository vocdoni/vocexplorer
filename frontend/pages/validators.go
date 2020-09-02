package pages

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/components"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	router "marwan.io/vecty-router"
)

// ValidatorsView renders the Validators page
type ValidatorsView struct {
	vecty.Core
	Cfg *config.Cfg
}

// Render renders the ValidatorsView component
func (home *ValidatorsView) Render() vecty.ComponentOrHTML {
	address, ok := router.GetNamedVar(home)["id"]
	// If there is an ID to look for, render individual validator page
	if ok && address != "" {
		dispatcher.Dispatch(&actions.SetCurrentValidatorID{ID: address})
		dash := new(components.ValidatorContents)
		dash.Rendered = false
		// Ensure component rerender is only triggered once component has been rendered
		if !store.Listeners.Has(dash) {
			store.Listeners.Add(dash, func() {
				if dash.Rendered {
					vecty.Rerender(dash)
				}
			})
		}
		go dash.UpdateValidatorContents()
		return elem.Div(dash)
	}
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
	go components.UpdateAndRenderValidatorsDashboard(dash)
	return elem.Div(dash)
}

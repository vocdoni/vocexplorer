package pages

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/vocdoni/vocexplorer/frontend/actions"
	"github.com/vocdoni/vocexplorer/frontend/components"
	"github.com/vocdoni/vocexplorer/frontend/dispatcher"
	"github.com/vocdoni/vocexplorer/frontend/store"
	router "marwan.io/vecty-router"
)

// ValidatorView renders the Validator page
type ValidatorView struct {
	vecty.Core
}

// Render renders the ValidatorView component
func (home *ValidatorView) Render() vecty.ComponentOrHTML {
	dispatcher.Dispatch(&actions.SetCurrentPage{Page: "validator"})
	dispatcher.Dispatch(&actions.SetCurrentValidatorID{ID: router.GetNamedVar(home)["id"]})
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
	return elem.Div(
		&components.Header{},
		dash,
	)
}

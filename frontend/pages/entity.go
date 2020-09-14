package pages

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/components"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	router "marwan.io/vecty-router"
)

// EntityView renders the Entity page
type EntityView struct {
	vecty.Core
}

// Render renders the EntityView component
func (home *EntityView) Render() vecty.ComponentOrHTML {
	dash := new(components.EntityContentsView)
	dispatcher.Dispatch(&actions.SetCurrentEntityID{EntityID: router.GetNamedVar(home)["id"]})
	dash.Rendered = false
	// Ensure component rerender is only triggered once component has been rendered
	if !store.Listeners.Has(dash) {
		store.Listeners.Add(dash, func() {
			if dash.Rendered {
				vecty.Rerender(dash)
			}
		})
	}
	go components.UpdateEntityContents(dash)
	return elem.Div(
		&components.Header{},
		dash,
	)
}

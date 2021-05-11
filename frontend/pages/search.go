package pages

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	router "marwan.io/vecty-router"

	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/components"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
)

// SearchView renders the Process page
type SearchView struct {
	vecty.Core
}

// Render renders the SearchView component
func (home *SearchView) Render() vecty.ComponentOrHTML {
	dispatcher.Dispatch(&actions.SetCurrentPage{Page: "search"})
	searchTerm := router.GetNamedVar(home)["searchTerm"]
	dash := new(components.SearchItemsView)
	dash.Rendered = false
	// Ensure component rerender is only triggered once component has been rendered
	if !store.Listeners.Has(dash) {
		store.Listeners.Add(dash, func() {
			if dash.Rendered {
				vecty.Rerender(dash)
			}
		})
	}
	go dash.UpdateSearchItems(searchTerm)
	return elem.Div(
		&components.Header{},
		dash,
	)
}

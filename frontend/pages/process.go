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

// ProcessView renders the Process page
type ProcessView struct {
	vecty.Core
}

// Render renders the ProcessView component
func (home *ProcessView) Render() vecty.ComponentOrHTML {
	dispatcher.Dispatch(&actions.SetCurrentPage{Page: "process"})
	dash := new(components.ProcessContentsView)
	dispatcher.Dispatch(&actions.SetCurrentProcessStruct{Process: &dbtypes.Process{ID: router.GetNamedVar(home)["id"]}})
	dash.Rendered = false
	// Ensure component rerender is only triggered once component has been rendered
	if !store.Listeners.Has(dash) {
		store.Listeners.Add(dash, func() {
			if dash.Rendered {
				vecty.Rerender(dash)
			}
		})
	}
	go components.UpdateProcessContents(dash)
	return elem.Div(
		&components.Header{},
		dash,
	)
}

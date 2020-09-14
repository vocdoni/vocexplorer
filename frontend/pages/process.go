package pages

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/components"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/proto"
	router "marwan.io/vecty-router"
)

// ProcessView renders the Process page
type ProcessView struct {
	vecty.Core
}

// Render renders the ProcessView component
func (home *ProcessView) Render() vecty.ComponentOrHTML {
	dash := new(components.ProcessContentsView)
	dispatcher.Dispatch(&actions.SetCurrentProcessStruct{Process: &proto.Process{ID: router.GetNamedVar(home)["id"]}})
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

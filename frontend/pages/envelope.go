package pages

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/components"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/util"
	router "marwan.io/vecty-router"
)

const alpha = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// EnvelopeView renders the Envelope page
type EnvelopeView struct {
	vecty.Core
}

// Render renders the EnvelopeView component
func (home *EnvelopeView) Render() vecty.ComponentOrHTML {
	dispatcher.Dispatch(&actions.SetCurrentPage{Page: "envelope"})
	id := util.StringToHex(router.GetNamedVar(home)["id"])
	dispatcher.Dispatch(&actions.SetCurrentEnvelopeNullifier{
		Nullifier: id,
	})
	dash := new(components.EnvelopeContents)
	dash.Rendered = false
	// Ensure component rerender is only triggered once component has been rendered
	if !store.Listeners.Has(dash) {
		store.Listeners.Add(dash, func() {
			if dash.Rendered {
				vecty.Rerender(dash)
			}
		})
	}
	go components.UpdateEnvelopeContents(dash)
	return elem.Div(
		&components.Header{},
		dash,
	)
}

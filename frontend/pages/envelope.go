package pages

import (
	"strconv"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/go-dvote/log"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/components"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	router "marwan.io/vecty-router"
)

// EnvelopeView renders the Envelope page
type EnvelopeView struct {
	vecty.Core
}

// Render renders the EnvelopeView component
func (home *EnvelopeView) Render() vecty.ComponentOrHTML {
	height, err := strconv.ParseInt(router.GetNamedVar(home)["id"], 0, 64)
	if err != nil {
		log.Error(err)
	}
	dispatcher.Dispatch(&actions.SetCurrentEnvelopeHeight{Height: height})
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
	return elem.Div(dash)
}

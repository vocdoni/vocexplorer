package pages

import (
	"strconv"
	"strings"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/components"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/logger"
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
	id := router.GetNamedVar(home)["id"]
	// If id contains letters, treat it as nullifier rather than height
	if strings.ContainsAny(id, alpha) {
		if strings.HasPrefix(id, "0x") {
			id = id[2:]
		}
		dash := new(components.EnvelopeNullifier)
		dash.Nullifier = id
		dash.Rendered = false
		dash.Unavailable = false
		go dash.LoadEnvelopeHeight()
		return elem.Div(&components.Header{}, dash)
	}
	height, err := strconv.ParseInt(id, 0, 64)
	if err != nil {
		logger.Error(err)
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
	return elem.Div(
		&components.Header{},
		dash,
	)
}

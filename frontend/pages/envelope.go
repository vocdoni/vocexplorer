package pages

import (
	"strconv"

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
	blockHeight, err := strconv.ParseInt(router.GetNamedVar(home)["block"], 0, 64)
	if err != nil {
		logger.Error(err)
	}
	txIndex, err := strconv.ParseInt(router.GetNamedVar(home)["index"], 0, 32)
	if err != nil {
		logger.Error(err)
	}
	dispatcher.Dispatch(&actions.SetCurrentEnvelopeReference{
		BlockHeight: uint32(blockHeight),
		TxIndex:     int32(txIndex),
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

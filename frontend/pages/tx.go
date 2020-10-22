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

// TxView renders the individual tx view
type TxView struct {
	vecty.Core
}

// Render renders the TxView component
func (home *TxView) Render() vecty.ComponentOrHTML {
	dispatcher.Dispatch(&actions.SetCurrentPage{Page: "tx"})
	height, err := strconv.ParseInt(router.GetNamedVar(home)["id"], 0, 64)
	if err != nil {
		logger.Error(err)
	}
	dispatcher.Dispatch(&actions.SetCurrentTransactionHeight{Height: height})
	dash := new(components.TxContents)
	dash.Rendered = false
	// Ensure component rerender is only triggered once component has been rendered
	if !store.Listeners.Has(dash) {
		store.Listeners.Add(dash, func() {
			if dash.Rendered {
				vecty.Rerender(dash)
			}
		})
	}
	go components.UpdateTxContents(dash)
	return elem.Div(
		&components.Header{},
		dash,
	)
}

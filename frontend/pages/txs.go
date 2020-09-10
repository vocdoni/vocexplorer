package pages

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/components"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
)

// TxsView renders the transactions search page
type TxsView struct {
	vecty.Core
	Cfg *config.Cfg
}

// Render renders the TxsView component
func (bv *TxsView) Render() vecty.ComponentOrHTML {
	dash := new(components.TxsDashboardView)
	dash.Rendered = false
	// Ensure component rerender is only triggered once component has been rendered
	if !store.Listeners.Has(dash) {
		store.Listeners.Add(dash, func() {
			if dash.Rendered {
				vecty.Rerender(dash)
			}
		})
	}
	go components.UpdateTxsDashboard(dash)
	return elem.Div(dash)
}

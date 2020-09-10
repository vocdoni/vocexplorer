package pages

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/components"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
)

// BlocksView renders the Blocks search page
type BlocksView struct {
	vecty.Core
	Cfg *config.Cfg
}

// Render renders the BlocksView component
func (bv *BlocksView) Render() vecty.ComponentOrHTML {
	dash := new(components.BlocksDashboardView)
	dash.Rendered = false
	// Ensure component rerender is only triggered once component has been rendered
	if !store.Listeners.Has(dash) {
		store.Listeners.Add(dash, func() {
			if dash.Rendered {
				vecty.Rerender(dash)
			}
		})
	}
	go components.UpdateBlocksDashboard(dash)
	return elem.Div(dash)
}

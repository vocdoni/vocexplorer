package pages

import (
	"github.com/gopherjs/vecty"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/components"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
)

// BlockTxsView renders the BlockTxs search page
type BlockTxsView struct {
	vecty.Core
	Cfg *config.Cfg
}

// Render renders the BlockTxsView component
func (bv *BlockTxsView) Render() vecty.ComponentOrHTML {
	dash := new(components.BlockTxsDashboardView)
	dash.Rendered = false
	// Ensure component rerender is only triggered once component has been rendered
	store.Listeners.Add(dash, func() {
		if dash.Rendered {
			vecty.Rerender(dash)
		}
	})
	go components.UpdateAndRenderBlockTxsDashboard(dash)
	return dash
}

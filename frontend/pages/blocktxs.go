package pages

import (
	"github.com/gopherjs/vecty"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/components"
)

// BlockTxsView renders the BlockTxs search page
type BlockTxsView struct {
	vecty.Core
	Cfg *config.Cfg
}

// Render renders the BlockTxsView component
func (bv *BlockTxsView) Render() vecty.ComponentOrHTML {
	dash := new(components.BlockTxsDashboardView)
	return dash
}

package components

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/rpc"
)

// BlockTxsView renders the BlockTxs search page
type BlockTxsView struct {
	vecty.Core
	cfg *config.Cfg
}

// Render renders the BlockTxsView component
func (bv *BlockTxsView) Render() vecty.ComponentOrHTML {
	var t rpc.TendermintInfo
	var dash BlockTxsDashboardView
	return elem.Div(
		&Header{},
		initBlockTxsDashboardView(&t, &dash, bv.cfg),
	)
}

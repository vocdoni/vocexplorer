package components

import (
	"syscall/js"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/rpc"
)

// DashboardView renders the dashboard landing page
type DashboardView struct {
	vecty.Core
}

// Render renders the DashboardView component
func (dash *DashboardView) Render() vecty.ComponentOrHTML {
	js.Global().Set("page", "dashboard")
	js.Global().Set("gateway", false)
	js.Global().Set("tendermint", false)
	var t rpc.TendermintInfo
	var gw client.GatewayInfo
	return elem.Div(
		&Header{currentPage: "dashboard"},
		elem.Main(
			vecty.Markup(vecty.Class("info-pane")),
			initStatsView(&t, &gw)
			initBlocksView(&t),
			initGatewayView(&gw),
		))
}

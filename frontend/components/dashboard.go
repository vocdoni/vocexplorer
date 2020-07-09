package components

import (
	"syscall/js"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/client"
)

// DashboardView renders the dashboard landing page
type DashboardView struct {
	vecty.Core
	val int
}

// Render renders the DashboardView component
func (dash *DashboardView) Render() vecty.ComponentOrHTML {
	js.Global().Set("page", "dashboard")
	js.Global().Set("gateway", false)
	var d client.MetaResponse
	return elem.Div(
		&Header{currentPage: "dashboard"},
		elem.Main(
			vecty.Markup(vecty.Class("info-pane")),
			initGatewayView(&d),
		))
}

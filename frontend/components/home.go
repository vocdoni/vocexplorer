package components

import (
	"syscall/js"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/client"
	"gitlab.com/vocdoni/vocexplorer/rpc"
)

// HomeView renders the Home landing page
type HomeView struct {
	vecty.Core
}

// Render renders the HomeView component
func (home *HomeView) Render() vecty.ComponentOrHTML {
	js.Global().Set("page", "home")
	var t rpc.TendermintInfo
	var vc client.VochainInfo
	var dash DashboardView
	return elem.Div(
		&Header{},
		elem.Main(
			initDashboardView(&t, &vc, &dash),
		),
	)
}

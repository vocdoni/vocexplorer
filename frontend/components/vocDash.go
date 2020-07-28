package components

import (
	"syscall/js"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/client"
	"gitlab.com/vocdoni/vocexplorer/config"
)

// VocDashView renders the processes page
type VocDashView struct {
	vecty.Core
	cfg *config.Cfg
}

// Render renders the VocDashView component
func (home *VocDashView) Render() vecty.ComponentOrHTML {
	js.Global().Set("page", "Voting Processes & Entities")
	var vc client.VochainInfo
	var dash VocDashDashboardView
	return elem.Div(
		&Header{},
		elem.Main(
			initVocDashDashboardView(&vc, &dash, home.cfg),
		),
	)
}

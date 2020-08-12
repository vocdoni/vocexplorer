package pages

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/client"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/components"
)

// VocDashView renders the processes page
type VocDashView struct {
	vecty.Core
	Cfg *config.Cfg
}

// Render renders the VocDashView component
func (home *VocDashView) Render() vecty.ComponentOrHTML {
	var vc client.VochainInfo
	var dash components.VocDashDashboardView
	return elem.Div(
		&components.Header{},
		elem.Main(
			components.InitVocDashDashboardView(&vc, &dash, home.Cfg),
		),
	)
}

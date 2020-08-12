package pages

import (
	"github.com/gopherjs/vecty"
	"gitlab.com/vocdoni/vocexplorer/client"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/components"
	"gitlab.com/vocdoni/vocexplorer/rpc"
)

// HomeView renders the Home landing page
type HomeView struct {
	vecty.Core
	Cfg *config.Cfg
}

// Render renders the HomeView component
func (home *HomeView) Render() vecty.ComponentOrHTML {
	var t rpc.TendermintInfo
	var vc client.VochainInfo
	var dash components.DashboardView
	return components.InitDashboardView(&t, &vc, &dash, home.Cfg)
}

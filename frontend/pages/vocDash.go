package pages

import (
	"github.com/gopherjs/vecty"
	"gitlab.com/vocdoni/vocexplorer/client"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/components"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
)

// VocDashView renders the processes page
type VocDashView struct {
	vecty.Core
	Cfg *config.Cfg
}

// Render renders the VocDashView component
func (home *VocDashView) Render() vecty.ComponentOrHTML {
	dash := new(components.VocDashDashboardView)
	dash.Vc = new(client.VochainInfo)
	dash.QuitCh = make(chan struct{})
	dash.RefreshEnvelopes = make(chan int, 50)
	dash.RefreshProcesses = make(chan int, 50)
	dash.RefreshEntities = make(chan int, 50)
	dash.DisableEnvelopesUpdate = false
	dash.RefreshEntities = store.Entities.PagChannel
	dash.RefreshProcesses = store.Processes.PagChannel
	dash.ServerConnected = true
	dash.GatewayConnected = true
	rendered := false
	dash.Rendered = &rendered
	go components.UpdateAndRenderVocDashDashboard(dash, home.Cfg)
	return dash
	// return components.InitVocDashDashboardView(vc, dash, home.Cfg)
}

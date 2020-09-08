package components

import (
	"fmt"
	"time"

	"github.com/gopherjs/vecty"
	"gitlab.com/vocdoni/vocexplorer/api"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/frontend/update"
)

// StatsDashboardView renders the dashboard landing page
type StatsDashboardView struct {
	vecty.Core
	vecty.Mounter
	Rendered bool
}

// Mount is called after the component renders to signal that it can be rerendered safely
func (dash *StatsDashboardView) Mount() {
	if !dash.Rendered {
		dash.Rendered = true
		vecty.Rerender(dash)
	}
}

// Render renders the StatsDashboardView component
func (dash *StatsDashboardView) Render() vecty.ComponentOrHTML {
	if !dash.Rendered {
		return LoadingBar()
	}
	if dash != nil {
		return Container(
			renderGatewayConnectionBanner(),
			renderServerConnectionBanner(),
			&BlockchainInfo{},
		)
	}
	return &bootstrap.Alert{
		Contents: "Connecting to blockchain clients",
		Type:     "warning",
	}
}

// UpdateAndRenderStatsDashboard keeps the stats dashboard updated
func UpdateAndRenderStatsDashboard(d *StatsDashboardView) {
	actions.EnableUpdates()
	ticker := time.NewTicker(time.Duration(store.Config.RefreshTime) * time.Second)
	updateStatsDashboard(d)
	for {
		select {
		case <-store.RedirectChan:
			fmt.Println("Redirecting...")
			ticker.Stop()
			return
		case <-ticker.C:
			updateStatsDashboard(d)
		}
	}
}

func updateStatsDashboard(d *StatsDashboardView) {
	dispatcher.Dispatch(&actions.GatewayConnected{Connected: store.GatewayClient.Ping()})
	dispatcher.Dispatch(&actions.ServerConnected{Connected: api.PingServer()})

	actions.UpdateCounts()
	update.DashboardInfo(store.GatewayClient)
	update.BlockchainStatus(store.TendermintClient)
}

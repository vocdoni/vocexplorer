package components

import (
	"fmt"
	"time"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/api"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/rpc"
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
		return elem.Div(vecty.Text("Loading..."))
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
	dispatcher.Dispatch(&actions.GatewayConnected{Connected: api.PingGateway(store.Config.GatewayHost)})
	dispatcher.Dispatch(&actions.ServerConnected{Connected: api.Ping()})

	actions.UpdateCounts()
	rpc.UpdateBlockchainStatus(store.TendermintClient)
}

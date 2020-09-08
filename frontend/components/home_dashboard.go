package components

import (
	"fmt"
	"time"

	"github.com/gopherjs/vecty"
	"gitlab.com/vocdoni/vocexplorer/api"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/frontend/update"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// DashboardView renders the dashboard landing page
type DashboardView struct {
	vecty.Core
	vecty.Mounter
	Rendered bool
}

// Mount is called after the component renders to signal that it can be rerendered safely
func (dash *DashboardView) Mount() {
	if !dash.Rendered {
		dash.Rendered = true
		vecty.Rerender(dash)
	}
}

// Render renders the DashboardView component
func (dash *DashboardView) Render() vecty.ComponentOrHTML {
	if !dash.Rendered {
		return LoadingBar()
	}
	if dash != nil && store.GatewayClient != nil && store.TendermintClient != nil {
		return Container(
			renderGatewayConnectionBanner(),
			renderServerConnectionBanner(),
			&StatsView{},
		)
	}
	return &bootstrap.Alert{
		Contents: "Connecting to blockchain clients",
		Type:     "warning",
	}
}

// UpdateAndRenderHomeDashboard keeps the home dashboard data up to date
func UpdateAndRenderHomeDashboard(d *DashboardView) {
	actions.EnableUpdates()
	ticker := time.NewTicker(time.Duration(util.Max(store.Config.RefreshTime, 1)) * time.Second)
	updateHomeDashboardInfo(d)
	for {
		select {
		case <-store.RedirectChan:
			fmt.Println("Redirecting...")
			ticker.Stop()
			return
		case <-ticker.C:
			updateHomeDashboardInfo(d)
		}
	}
}

func updateHomeDashboardInfo(d *DashboardView) {
	dispatcher.Dispatch(&actions.GatewayConnected{Connected: store.GatewayClient.Ping()})
	dispatcher.Dispatch(&actions.ServerConnected{Connected: api.PingServer()})

	update.BlockchainStatus(store.TendermintClient)
	update.DashboardInfo(store.GatewayClient)
	actions.UpdateCounts()
	updateHomeBlocks(d, util.Max(store.Blocks.Count, config.HomeWidgetBlocksListSize))
}

func updateHomeBlocks(d *DashboardView, index int) {
	fmt.Println("Getting blocks from index " + util.IntToString(index))
	list, ok := api.GetBlockList(index)
	if ok {
		reverseBlockList(&list)
		dispatcher.Dispatch(&actions.SetBlockList{BlockList: list})
	}
}

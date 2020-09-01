package components

import (
	"fmt"
	"time"

	"github.com/gopherjs/vecty"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/api"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/rpc"
	"gitlab.com/vocdoni/vocexplorer/update"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// DashboardView renders the dashboard landing page
type DashboardView struct {
	vecty.Core
	blockIndex int
}

// Render renders the DashboardView component
func (dash *DashboardView) Render() vecty.ComponentOrHTML {
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
	vecty.Rerender(d)
	for {
		select {
		case <-store.RedirectChan:
			fmt.Println("Redirecting...")
			ticker.Stop()
			return
		case <-ticker.C:
			updateHomeDashboardInfo(d)
			vecty.Rerender(d)
		}
	}
}

func updateHomeDashboardInfo(d *DashboardView) {
	dispatcher.Dispatch(&actions.GatewayConnected{Connected: api.PingGateway(store.Config.GatewayHost)})
	dispatcher.Dispatch(&actions.ServerConnected{Connected: api.Ping()})

	rpc.UpdateBlockchainStatus(store.TendermintClient)
	update.DashboardInfo(store.GatewayClient)
	actions.UpdateCounts()
	updateHomeBlocks(d, util.Max(store.Blocks.Count-d.blockIndex-1, config.HomeWidgetBlocksListSize))
}

func updateHomeBlocks(d *DashboardView, index int) {
	fmt.Println("Getting blocks from index " + util.IntToString(index))
	list, ok := api.GetBlockList(index)
	if ok {
		reverseBlockList(&list)
		dispatcher.Dispatch(&actions.SetBlockList{BlockList: list})
	}
}

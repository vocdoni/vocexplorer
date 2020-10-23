package components

import (
	"time"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/api"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/frontend/update"
	"gitlab.com/vocdoni/vocexplorer/logger"
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
	return elem.Div(
		vecty.Markup(vecty.Attribute("id", "main")),
		renderServerConnectionBanner(),
		&StatsView{},
	)
}

// UpdateHomeDashboard keeps the home dashboard data up to date
func UpdateHomeDashboard(d *DashboardView) {
	dispatcher.Dispatch(&actions.EnableAllUpdates{})
	ticker := time.NewTicker(time.Duration(util.Max(store.Config.RefreshTime, 1)) * time.Second)
	if !update.CheckCurrentPage("home", ticker) {
		return
	}
	updateHomeDashboardInfo(d)
	for {
		select {
		case <-store.RedirectChan:
			if !update.CheckCurrentPage("home", ticker) {
				return
			}
		case <-ticker.C:
			if !update.CheckCurrentPage("home", ticker) {
				return
			}
			updateHomeDashboardInfo(d)
		}
	}
}

func updateHomeDashboardInfo(d *DashboardView) {
	dispatcher.Dispatch(&actions.ServerConnected{Connected: api.PingServer()})
	actions.UpdateCounts()
	updateHomeBlocks(d, util.Max(store.Blocks.Count, config.HomeWidgetBlocksListSize))
}

func updateHomeBlocks(d *DashboardView, index int) {
	logger.Info("Getting blocks from index " + util.IntToString(index))
	list, ok := api.GetBlockList(index)
	if ok {
		reverseBlockList(&list)
		dispatcher.Dispatch(&actions.SetBlockList{BlockList: list})
	}
}

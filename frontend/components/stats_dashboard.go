package components

import (
	"time"

	"github.com/hexops/vecty"
	"github.com/vocdoni/vocexplorer/api"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/frontend/update"
	"gitlab.com/vocdoni/vocexplorer/logger"
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
	return Container(
		vecty.Markup(vecty.Attribute("id", "main")),
		renderServerConnectionBanner(),
		&BlockchainInfo{header: true},
	)
}

// UpdateStatsDashboard keeps the stats dashboard updated
func UpdateStatsDashboard(d *StatsDashboardView) {
	dispatcher.Dispatch(&actions.EnableAllUpdates{})
	ticker := time.NewTicker(time.Duration(store.Config.RefreshTime) * time.Second)
	if !update.CheckCurrentPage("stats", ticker) {
		return
	}
	updateStatsDashboard(d)
	for {
		select {
		case <-store.RedirectChan:
			if !update.CheckCurrentPage("stats", ticker) {
				return
			}
		case <-ticker.C:
			if !update.CheckCurrentPage("stats", ticker) {
				return
			}
			updateStatsDashboard(d)
		}
	}
}

func updateStatsDashboard(d *StatsDashboardView) {
	dispatcher.Dispatch(&actions.ServerConnected{Connected: api.PingServer()})

	stats, err := store.Client.GetStats()
	if err != nil {
		logger.Error(err)
		return
	}
	actions.UpdateCounts(stats)
}

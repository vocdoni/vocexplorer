package components

import (
	"time"

	"github.com/hexops/vecty"
	"go.vocdoni.io/proto/build/go/models"

	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/frontend/update"
	"gitlab.com/vocdoni/vocexplorer/logger"
)

// ValidatorsDashboardView renders the validators list dashboard
type ValidatorsDashboardView struct {
	vecty.Core
	vecty.Mounter
	Rendered bool
}

// Mount is called after the component renders to signal that it can be rerendered safely
func (dash *ValidatorsDashboardView) Mount() {
	if !dash.Rendered {
		dash.Rendered = true
		vecty.Rerender(dash)
	}
}

// Render renders the ValidatorsDashboardView component
func (dash *ValidatorsDashboardView) Render() vecty.ComponentOrHTML {
	if !dash.Rendered {
		return LoadingBar()
	}
	return Container(
		vecty.Markup(vecty.Attribute("id", "main")),
		renderServerConnectionBanner(),
		&ValidatorListView{},
	)
}

// UpdateValidatorsDashboard keeps the validators data up to date
func UpdateValidatorsDashboard(d *ValidatorsDashboardView) {
	dispatcher.Dispatch(&actions.EnableAllUpdates{})

	ticker := time.NewTicker(time.Duration(store.Config.RefreshTime) * time.Second)
	if !update.CheckCurrentPage("validators", ticker) {
		return
	}
	updateValidatorsDashboard(d)
	for {
		select {
		case <-store.RedirectChan:
			if !update.CheckCurrentPage("validators", ticker) {
				return
			}
		case <-ticker.C:
			if !update.CheckCurrentPage("validators", ticker) {
				return
			}
			updateValidatorsDashboard(d)
		}
	}
}

func updateValidatorsDashboard(d *ValidatorsDashboardView) {
	dispatcher.Dispatch(&actions.GatewayConnected{GatewayErr: store.Client.GetGatewayInfo()})
	stats, err := store.Client.GetStats()
	if err != nil {
		logger.Error(err)
		return
	}
	actions.UpdateCounts(stats)
	updateValidators(d)
}

func updateValidators(d *ValidatorsDashboardView) {
	list, err := store.Client.GetValidatorList()
	if err != nil {
		logger.Error(err)
	} else {
		reverseValidatorList(list)
		dispatcher.Dispatch(&actions.SetValidatorList{List: list})
	}
}

func reverseValidatorList(list []*models.Validator) {
	for i := len(list)/2 - 1; i >= 0; i-- {
		opp := len(list) - 1 - i
		list[i], list[opp] = list[opp], list[i]
	}
}

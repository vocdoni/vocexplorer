package components

import (
	"fmt"
	"time"

	"github.com/hexops/vecty"
	"gitlab.com/vocdoni/vocexplorer/api"
	"gitlab.com/vocdoni/vocexplorer/api/dbtypes"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/frontend/update"
	"gitlab.com/vocdoni/vocexplorer/logger"
	"gitlab.com/vocdoni/vocexplorer/util"
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
		case i := <-store.Validators.Pagination.PagChannel:
			if !update.CheckCurrentPage("validators", ticker) {
				return
			}
		loop:
			for {
				// If many indices waiting in buffer, scan to last one.
				select {
				case i = <-store.Validators.Pagination.PagChannel:
				default:
					break loop
				}
			}
			dispatcher.Dispatch(&actions.ValidatorsIndexChange{Index: i})
			if i < 1 {
				newHeight, _ := api.GetValidatorCount()
				dispatcher.Dispatch(&actions.SetValidatorCount{Count: int(newHeight)})
			}
			updateValidators(d, util.Max(store.Validators.Count-store.Validators.Pagination.Index, 1))
		case search := <-store.Validators.Pagination.SearchChannel:
			if !update.CheckCurrentPage("validators", ticker) {
				return
			}
		validatorSearch:
			for {
				// If many indices waiting in buffer, scan to last one.
				select {
				case search = <-store.Validators.Pagination.SearchChannel:
				default:
					break validatorSearch
				}
			}
			logger.Info("search: " + search)
			dispatcher.Dispatch(&actions.ValidatorsIndexChange{Index: 0})
			list, ok := api.GetValidatorSearch(search)
			if ok {
				reverseValidatorList(&list)
				dispatcher.Dispatch(&actions.SetValidatorList{List: list})
			} else {
				dispatcher.Dispatch(&actions.SetValidatorList{List: [config.ListSize]*dbtypes.Validator{}})
			}
		}
	}
}

func updateValidatorsDashboard(d *ValidatorsDashboardView) {
	dispatcher.Dispatch(&actions.ServerConnected{Connected: api.PingServer()})
	if !store.Validators.Pagination.DisableUpdate {
		actions.UpdateCounts()
		updateValidators(d, util.Max(store.Validators.Count-store.Validators.Pagination.Index, 1))
	}
}

func updateValidators(d *ValidatorsDashboardView, index int) {
	logger.Info(fmt.Sprintf("Getting Validators from index %d\n", index))
	list, ok := api.GetValidatorList(index)
	if ok {
		reverseValidatorList(&list)
		dispatcher.Dispatch(&actions.SetValidatorList{List: list})
	}
	blockHeights, ok := api.GetValidatorBlockHeightMap()
	if ok {
		dispatcher.Dispatch(&actions.SetValidatorBlockHeightMap{HeightMap: blockHeights})
	}
}

func reverseValidatorList(list *[config.ListSize]*dbtypes.Validator) {
	for i := len(list)/2 - 1; i >= 0; i-- {
		opp := len(list) - 1 - i
		list[i], list[opp] = list[opp], list[i]
	}
}

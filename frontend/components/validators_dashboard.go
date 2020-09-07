package components

import (
	"fmt"
	"time"

	"github.com/gopherjs/vecty"
	"gitlab.com/vocdoni/vocexplorer/api"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/frontend/update"
	"gitlab.com/vocdoni/vocexplorer/proto"
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
		renderGatewayConnectionBanner(),
		renderServerConnectionBanner(),
		&ValidatorListView{},
	)
}

// UpdateAndRenderValidatorsDashboard keeps the validators data up to date
func UpdateAndRenderValidatorsDashboard(d *ValidatorsDashboardView) {
	dispatcher.Dispatch(&actions.EnableAllUpdates{})

	ticker := time.NewTicker(time.Duration(store.Config.RefreshTime) * time.Second)
	updateValidatorsDashboard(d)
	for {
		select {
		case <-store.RedirectChan:
			fmt.Println("Redirecting...")
			ticker.Stop()
			return
		case <-ticker.C:
			updateValidatorsDashboard(d)
		case i := <-store.Validators.Pagination.PagChannel:
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
			oldValidators := store.Validators.Count
			newHeight, _ := api.GetValidatorCount()
			dispatcher.Dispatch(&actions.SetValidatorCount{Count: int(newHeight)})
			if i < 1 {
				oldValidators = store.Validators.Count
			}
			updateValidators(d, util.Max(oldValidators-store.Validators.Pagination.Index, 1))
		}
	}
}

func updateValidatorsDashboard(d *ValidatorsDashboardView) {
	dispatcher.Dispatch(&actions.GatewayConnected{Connected: store.GatewayClient.Ping()})
	dispatcher.Dispatch(&actions.ServerConnected{Connected: api.PingServer()})
	actions.UpdateCounts()
	update.BlockchainStatus(store.TendermintClient)
	if !store.Validators.Pagination.DisableUpdate {
		updateValidators(d, util.Max(store.Validators.Count-store.Validators.Pagination.Index, 1))
	}
}

func updateValidators(d *ValidatorsDashboardView, index int) {
	fmt.Printf("Getting Validators from index %d\n", index)
	list, ok := api.GetValidatorList(index)
	if ok {
		reverseValidatorList(&list)
		dispatcher.Dispatch(&actions.SetValidatorList{List: list})
	}
}

func reverseValidatorList(list *[config.ListSize]*proto.Validator) {
	for i := len(list)/2 - 1; i >= 0; i-- {
		opp := len(list) - 1 - i
		list[i], list[opp] = list[opp], list[i]
	}
}

package components

import (
	"fmt"
	"time"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/api"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/rpc"
	"gitlab.com/vocdoni/vocexplorer/types"
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
		return elem.Div(vecty.Text("Loading..."))
	}
	return Container(
		renderGatewayConnectionBanner(),
		renderServerConnectionBanner(),
		&ValidatorListView{},
		&BlockchainInfo{},
	)
}

// UpdateAndRenderValidatorsDashboard keeps the validators data up to date
func UpdateAndRenderValidatorsDashboard(d *ValidatorsDashboardView) {
	actions.EnableUpdates()
	actions.ResetIndexes()
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
			dispatcher.Dispatch(&actions.SetValidatorCount{Count: int(newHeight) - 1})
			if i < 1 {
				oldValidators = store.Validators.Count
			}
			updateValidators(d, util.Max(oldValidators-store.Validators.Pagination.Index, config.ListSize))
		}
	}
}

func updateValidatorsDashboard(d *ValidatorsDashboardView) {
	dispatcher.Dispatch(&actions.GatewayConnected{Connected: api.PingGateway(store.Config.GatewayHost)})
	dispatcher.Dispatch(&actions.ServerConnected{Connected: api.Ping()})
	actions.UpdateCounts()
	rpc.UpdateBlockchainStatus(store.TendermintClient)
	if !store.Validators.Pagination.DisableUpdate {
		updateValidators(d, util.Max(store.Validators.Count-store.Validators.Pagination.Index, config.ListSize))
	}
}

func updateValidators(d *ValidatorsDashboardView, index int) {
	fmt.Printf("Getting Blocks from index %d\n", index)
	list, ok := api.GetValidatorList(index)
	if ok {
		reverseValidatorList(&list)
		dispatcher.Dispatch(&actions.SetValidatorList{List: list})
	}
}

func reverseValidatorList(list *[config.ListSize]*types.Validator) {
	for i := len(list)/2 - 1; i >= 0; i-- {
		opp := len(list) - 1 - i
		list[i], list[opp] = list[opp], list[i]
	}
}

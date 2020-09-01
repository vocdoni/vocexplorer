package components

import (
	"fmt"
	"time"

	"github.com/gopherjs/vecty"
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
	Rendered       bool
	ValidatorIndex int
}

// Mount is called after the component renders to signal that it can be rerendered safely
func (dash *ValidatorsDashboardView) Mount() {
	dash.Rendered = true
}

// Render renders the ValidatorsDashboardView component
func (dash *ValidatorsDashboardView) Render() vecty.ComponentOrHTML {
	return Container(
		renderGatewayConnectionBanner(),
		renderServerConnectionBanner(),
		&ValidatorListView{
			refreshCh:     store.Validators.Pagination.PagChannel,
			disableUpdate: &store.Validators.Pagination.DisableUpdate,
		},
		&BlockchainInfo{},
	)
}

// UpdateAndRenderValidatorsDashboard keeps the validators data up to date
func UpdateAndRenderValidatorsDashboard(d *ValidatorsDashboardView) {
	actions.EnableUpdates()
	ticker := time.NewTicker(time.Duration(store.Config.RefreshTime) * time.Second)
	updateValidatorsDashboard(d)
	vecty.Rerender(d)
	for {
		select {
		case <-store.RedirectChan:
			fmt.Println("Redirecting...")
			ticker.Stop()
			return
		case <-ticker.C:
			updateValidatorsDashboard(d)
			vecty.Rerender(d)
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
			d.ValidatorIndex = i
			oldValidators := store.Validators.Count
			newHeight, _ := api.GetValidatorCount()
			dispatcher.Dispatch(&actions.SetValidatorCount{Count: int(newHeight) - 1})
			if i < 1 {
				oldValidators = store.Validators.Count
			}
			updateValidators(d, util.Max(oldValidators-d.ValidatorIndex, config.ListSize))

			vecty.Rerender(d)
		}
	}
}

func updateValidatorsDashboard(d *ValidatorsDashboardView) {
	dispatcher.Dispatch(&actions.GatewayConnected{Connected: api.PingGateway(store.Config.GatewayHost)})
	dispatcher.Dispatch(&actions.ServerConnected{Connected: api.Ping()})
	actions.UpdateCounts()
	rpc.UpdateBlockchainStatus(store.TendermintClient)
	if !store.Validators.Pagination.DisableUpdate {
		updateValidators(d, util.Max(store.Validators.Count-d.ValidatorIndex, config.ListSize))
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

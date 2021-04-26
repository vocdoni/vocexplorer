package components

import (
	"fmt"
	"time"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/frontend/update"
	"gitlab.com/vocdoni/vocexplorer/logger"
)

// EntitiesDashboardView renders the entities dashboard page
type EntitiesDashboardView struct {
	vecty.Core
	vecty.Mounter
	Rendered bool
}

// Mount is called after the component renders to signal that it can be rerendered safely
func (dash *EntitiesDashboardView) Mount() {
	if !dash.Rendered {
		dash.Rendered = true
		vecty.Rerender(dash)
	}
}

// Render renders the EntitiesDashboardView component
func (dash *EntitiesDashboardView) Render() vecty.ComponentOrHTML {
	if !dash.Rendered {
		return LoadingBar()
	}
	return Container(
		vecty.Markup(vecty.Attribute("id", "main")),
		renderServerConnectionBanner(),
		elem.Section(
			bootstrap.Card(bootstrap.CardParams{
				Body: vecty.List{
					elem.Heading1(vecty.Text("Entities")),
					&EntityListView{},
				},
			}),
		),
	)
}

// UpdateEntitiesDashboard continuously updates the information needed by the Entities dashboard
func UpdateEntitiesDashboard(d *EntitiesDashboardView) {
	dispatcher.Dispatch(&actions.EnableAllUpdates{})

	ticker := time.NewTicker(time.Duration(store.Config.RefreshTime) * time.Second)
	if !update.CheckCurrentPage("entities", ticker) {
		return
	}
	updateEntities(d)
	for {
		select {
		case <-store.RedirectChan:
			if !update.CheckCurrentPage("entities", ticker) {
				return
			}
		case <-ticker.C:
			if !update.CheckCurrentPage("entities", ticker) {
				return
			}
			updateEntities(d)
		case i := <-store.Entities.Pagination.PagChannel:
			if !update.CheckCurrentPage("entities", ticker) {
				return
			}
		entityLoop:
			for {
				// If many indices waiting in buffer, scan to last one.
				select {
				case i = <-store.Entities.Pagination.PagChannel:
				default:
					break entityLoop
				}
			}
			dispatcher.Dispatch(&actions.EntitiesIndexChange{Index: i})
			if i < 1 {
				newVal, err := store.Client.GetEntityCount()
				if err != nil {
					logger.Error(err)
				}
				dispatcher.Dispatch(&actions.SetEntityCount{Count: int(newVal)})
			}
			if store.Entities.Count > 0 {
				getEntities(d, store.Entities.Count-store.Entities.Pagination.Index-config.ListSize)
			}
		case search := <-store.Entities.Pagination.SearchChannel:
			if !update.CheckCurrentPage("entities", ticker) {
				return
			}
		entitySearch:
			for {
				// If many indices waiting in buffer, scan to last one.
				select {
				case search = <-store.Entities.Pagination.SearchChannel:
				default:
					break entitySearch
				}
			}
			logger.Info("search: " + search)
			dispatcher.Dispatch(&actions.EntitiesIndexChange{Index: 0})
			list, err := store.Client.GetEntityList(search, config.ListSize, 0)
			if err != nil {
				dispatcher.Dispatch(&actions.SetEntityIDs{EntityIDs: []string{}})
				logger.Error(err)
			} else {
				dispatcher.Dispatch(&actions.SetEntityIDs{EntityIDs: list})
			}
		}
	}
}

func updateEntities(d *EntitiesDashboardView) {
	dispatcher.Dispatch(&actions.GatewayConnected{GatewayErr: store.Client.GetGatewayInfo()})
	if !store.Entities.Pagination.DisableUpdate {
		stats, err := store.Client.GetStats()
		if err != nil {
			logger.Error(err)
			return
		}
		actions.UpdateCounts(stats)
		getEntities(d, store.Entities.Count-store.Entities.Pagination.Index-config.ListSize)
	}
}

func getEntities(d *EntitiesDashboardView, index int) {
	listSize := config.ListSize
	logger.Info(fmt.Sprintf("Getting entities from index %d\n", index))
	if index < 0 {
		listSize += index
		index = 0
	}
	logger.Info(fmt.Sprintf("Getting %d entities from index %d\n", listSize, index))
	list, err := store.Client.GetEntityList("", listSize, index)
	if err != nil {
		dispatcher.Dispatch(&actions.SetEntityIDs{EntityIDs: []string{}})
		logger.Error(err)
	} else {
		reverseIDList(list)
		dispatcher.Dispatch(&actions.SetEntityIDs{EntityIDs: list})
	}
	// TODO get entity process heights
}

func reverseIDList(list []string) {
	for i := len(list)/2 - 1; i >= 0; i-- {
		opp := len(list) - 1 - i
		list[i], list[opp] = list[opp], list[i]
	}
}

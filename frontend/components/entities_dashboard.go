package components

import (
	"fmt"
	"time"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/api"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/frontend/update"
	"gitlab.com/vocdoni/vocexplorer/logger"
	"gitlab.com/vocdoni/vocexplorer/util"
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
		renderServerConnectionBanner(),
		elem.Section(
			bootstrap.Card(bootstrap.CardParams{
				Body: vecty.List{
					elem.Heading2(vecty.Text("Entities")),
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
				newVal, _ := api.GetEntityCount()
				dispatcher.Dispatch(&actions.SetEntityCount{Count: int(newVal)})
			}
			if store.Entities.Count > 0 {
				getEntities(d, util.Max(store.Entities.Count-store.Entities.Pagination.Index, 1))
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
			list, ok := api.GetEntitySearch(search)
			if ok {
				dispatcher.Dispatch(&actions.SetEntityIDs{EntityIDs: list})
			} else {
				dispatcher.Dispatch(&actions.SetEntityIDs{EntityIDs: [config.ListSize]string{}})
			}
		}
	}
}

func updateEntities(d *EntitiesDashboardView) {
	dispatcher.Dispatch(&actions.ServerConnected{Connected: api.PingServer()})
	if !store.Entities.Pagination.DisableUpdate {
		actions.UpdateCounts()
		getEntities(d, util.Max(store.Entities.Count-store.Entities.Pagination.Index, 1))
	}
}

func getEntities(d *EntitiesDashboardView, index int) {
	index--
	logger.Info(fmt.Sprintf("Getting entities from index %d\n", index))
	list, ok := api.GetEntityList(index)
	if ok {
		dispatcher.Dispatch(&actions.SetEntityIDs{EntityIDs: list})
	}
	newVal, ok := api.GetEntityProcessCountMap()
	if ok {
		dispatcher.Dispatch(&actions.SetProcessHeights{ProcessHeights: newVal})
	}
}

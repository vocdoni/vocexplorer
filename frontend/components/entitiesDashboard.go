package components

import (
	"fmt"
	"time"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/api"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/update"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// EntitiesDashboardView renders the entities dashboard page
type EntitiesDashboardView struct {
	vecty.Core
	vecty.Mounter
	Rendered bool
}

//EntitiesTab is the tab component for entities
type EntitiesTab struct {
	*Tab
}

// Mount is called after the component renders to signal that it can be rerendered safely
func (dash *EntitiesDashboardView) Mount() {
	if !dash.Rendered {
		dash.Rendered = true
		vecty.Rerender(dash)
	}
}

func (e *EntitiesTab) dispatch() interface{} {
	return &actions.EntitiesTabChange{
		Tab: e.alias(),
	}
}

func (e *EntitiesTab) store() string {
	return store.Entities.Pagination.Tab
}

// Render renders the EntitiesDashboardView component
func (dash *EntitiesDashboardView) Render() vecty.ComponentOrHTML {
	if !dash.Rendered {
		return LoadingBar()
	}
	if dash == nil || store.GatewayClient == nil {
		return Container(&bootstrap.Alert{
			Type:     "warning",
			Contents: "Connecting to blockchain client",
		})
	}

	return Container(
		renderGatewayConnectionBanner(),
		renderServerConnectionBanner(),
		elem.Section(
			vecty.Markup(vecty.Class("details-view", "no-column")),
			elem.Div(
				vecty.Markup(vecty.Class("row")),
				elem.Div(
					vecty.Markup(vecty.Class("main-column")),
					bootstrap.Card(bootstrap.CardParams{
						Body: dash.EntityDetails(),
					}),
				),
			),
		),
		elem.Section(
			vecty.Markup(vecty.Class("row")),
			elem.Div(
				vecty.Markup(vecty.Class("col-12")),
				bootstrap.Card(bootstrap.CardParams{
					Body: &EntityProcessListView{},
				}),
			),
		),
	)
}

//EntityDetails renders the details of a single entity
func (dash *EntitiesDashboardView) EntityDetails() vecty.List {
	return vecty.List{
		elem.Heading1(
			vecty.Text("Entity details"),
		),
		elem.Heading2(vecty.Text(store.Entities.CurrentEntityID)),
		elem.Anchor(
			vecty.Markup(vecty.Class("hash")),
			vecty.Markup(vecty.Attribute("href", "https://manage.vocdoni.net/entities/#/0x"+store.Entities.CurrentEntityID)),
			vecty.Text("Entity Manager Page"),
		),
	}
}

// UpdateAndRenderEntitiesDashboard keeps the dashboard data up to date
func UpdateAndRenderEntitiesDashboard(d *EntitiesDashboardView) {
	actions.EnableUpdates()
	dispatcher.Dispatch(&actions.EntityProcessesIndexChange{Index: 0})
	dispatcher.Dispatch(&actions.EntityProcessesPageChange{Index: 0})
	ticker := time.NewTicker(time.Duration(store.Config.RefreshTime) * time.Second)
	updateEntityProcesses(d, util.Max(store.Entities.Count-store.Entities.ProcessesIndex, config.ListSize))
	for {
		select {
		case <-store.RedirectChan:
			fmt.Println("Redirecting...")
			ticker.Stop()
			return
		case <-ticker.C:
			updateEntityProcesses(d, util.Max(store.Entities.Count-store.Entities.ProcessesIndex, config.ListSize))
		case i := <-store.Entities.Pagination.PagChannel:
		loop:
			for {
				// If many indices waiting in buffer, scan to last one.
				select {
				case i = <-store.Entities.Pagination.PagChannel:
				default:
					break loop
				}
			}
			dispatcher.Dispatch(&actions.EntityProcessesIndexChange{Index: i})
			oldProcesses := store.Entities.Count
			newHeight, _ := api.GetEntityProcessHeight(store.Entities.CurrentEntityID)
			dispatcher.Dispatch(&actions.SetEntityCount{Count: int(newHeight)})
			if i < 1 {
				oldProcesses = store.Entities.Count
			}
			index := util.Max(oldProcesses-store.Entities.ProcessesIndex, config.ListSize)
			fmt.Printf("Getting processes from entity %s, index %d\n", store.Entities.CurrentEntityID, index)
			list, ok := api.GetProcessListByEntity(index, store.Entities.CurrentEntityID)
			if ok {
				reverseIDList(&list)
				dispatcher.Dispatch(&actions.SetEntityProcessList{ProcessList: list})
			}
			newMap, ok := api.GetProcessEnvelopeHeightMap()
			if ok {
				dispatcher.Dispatch(&actions.SetEnvelopeHeights{EnvelopeHeights: newMap})
			}
			update.EntityProcessResults()
		}
	}
}

func updateEntityProcesses(d *EntitiesDashboardView, index int) {
	dispatcher.Dispatch(&actions.GatewayConnected{Connected: api.PingGateway(store.Config.GatewayHost)})
	dispatcher.Dispatch(&actions.ServerConnected{Connected: api.Ping()})

	newCount, ok := api.GetEntityProcessHeight(store.Entities.CurrentEntityID)
	if ok {
		dispatcher.Dispatch(&actions.SetEntityProcessCount{Count: int(newCount)})
	}
	if store.Entities.CurrentEntity.ProcessCount > 0 && !store.Entities.Pagination.DisableUpdate {
		fmt.Printf("Getting processes from entity %s, index %d\n", store.Entities.CurrentEntityID, index)
		list, ok := api.GetProcessListByEntity(index, store.Entities.CurrentEntityID)
		if ok {
			reverseIDList(&list)
			dispatcher.Dispatch(&actions.SetEntityProcessList{ProcessList: list})
		}
		newMap, ok := api.GetProcessEnvelopeHeightMap()
		if ok {
			dispatcher.Dispatch(&actions.SetEnvelopeHeights{EnvelopeHeights: newMap})
		}
		update.EntityProcessResults()
	}
}

package components

import (
	"fmt"
	"time"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/api"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/frontend/update"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// EntityContentsView renders the entities dashboard page
type EntityContentsView struct {
	vecty.Core
	vecty.Mounter
	Rendered bool
}

//EntityTab is the tab component for entity
type EntityTab struct {
	*Tab
}

// Mount is called after the component renders to signal that it can be rerendered safely
func (dash *EntityContentsView) Mount() {
	if !dash.Rendered {
		dash.Rendered = true
		vecty.Rerender(dash)
	}
}

func (e *EntityTab) dispatch() interface{} {
	return &actions.EntityTabChange{
		Tab: e.alias(),
	}
}

func (e *EntityTab) store() string {
	return store.Entities.Pagination.Tab
}

// Render renders the EntityContentsView component
func (dash *EntityContentsView) Render() vecty.ComponentOrHTML {
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
func (dash *EntityContentsView) EntityDetails() vecty.List {
	return vecty.List{
		elem.Heading1(
			vecty.Text("Entity details"),
		),
		elem.Heading2(vecty.Text(store.Entities.CurrentEntityID)),
		elem.Anchor(
			vecty.Markup(vecty.Class("hash")),
			vecty.Markup(vecty.Attribute("href", "https://vocdoni.link/entities/0x"+store.Entities.CurrentEntityID)),
			vecty.Text("Entity Profile"),
		),
	}
}

// UpdateEntityContents keeps the dashboard data up to date
func UpdateEntityContents(d *EntityContentsView) {
	dispatcher.Dispatch(&actions.EnableAllUpdates{})
	ticker := time.NewTicker(time.Duration(store.Config.RefreshTime) * time.Second)
	dispatcher.Dispatch(&actions.GatewayConnected{Connected: store.GatewayClient.Ping()})
	dispatcher.Dispatch(&actions.ServerConnected{Connected: api.PingServer()})

	newCount, ok := api.GetEntityProcessCount(store.Entities.CurrentEntityID)
	if ok {
		dispatcher.Dispatch(&actions.SetEntityProcessCount{Count: int(newCount)})
	}
	updateEntityProcesses(d, util.Max(store.Entities.CurrentEntity.ProcessCount-store.Entities.ProcessPagination.Index, 1))
	for {
		select {
		case <-store.RedirectChan:
			fmt.Println("Redirecting...")
			ticker.Stop()
			return
		case <-ticker.C:
			updateEntityProcesses(d, util.Max(store.Entities.CurrentEntity.ProcessCount-store.Entities.ProcessPagination.Index, 1))
		case i := <-store.Entities.ProcessPagination.PagChannel:
		loop:
			for {
				// If many indices waiting in buffer, scan to last one.
				select {
				case i = <-store.Entities.ProcessPagination.PagChannel:
				default:
					break loop
				}
			}
			dispatcher.Dispatch(&actions.EntityProcessesIndexChange{Index: i})
			oldProcesses := store.Entities.CurrentEntity.ProcessCount
			newCount, ok := api.GetEntityProcessCount(store.Entities.CurrentEntityID)
			if ok {
				dispatcher.Dispatch(&actions.SetEntityProcessCount{Count: int(newCount)})
			}
			if i < 1 {
				oldProcesses = store.Entities.CurrentEntity.ProcessCount
			}
			index := util.Max(oldProcesses-store.Entities.ProcessPagination.Index, 1)
			fmt.Printf("Getting processes from entity %s, index %d\n", store.Entities.CurrentEntityID, index)
			list, ok := api.GetProcessListByEntity(index-1, store.Entities.CurrentEntityID)
			if ok {
				reverseIDList(&list)
				dispatcher.Dispatch(&actions.SetEntityProcessList{ProcessList: list})
			}
			newMap, ok := api.GetProcessEnvelopeCountMap()
			if ok {
				dispatcher.Dispatch(&actions.SetEnvelopeHeights{EnvelopeHeights: newMap})
			}
			update.EntityProcessResults()
		}
	}
}

func updateEntityProcesses(d *EntityContentsView, index int) {
	index--
	dispatcher.Dispatch(&actions.GatewayConnected{Connected: store.GatewayClient.Ping()})
	dispatcher.Dispatch(&actions.ServerConnected{Connected: api.PingServer()})

	newCount, ok := api.GetEntityProcessCount(store.Entities.CurrentEntityID)
	if ok {
		dispatcher.Dispatch(&actions.SetEntityProcessCount{Count: int(newCount)})
	}
	if store.Entities.CurrentEntity.ProcessCount > 0 && !store.Entities.ProcessPagination.DisableUpdate {
		fmt.Printf("Getting processes from entity %s, index %d\n", store.Entities.CurrentEntityID, index)
		list, ok := api.GetProcessListByEntity(index, store.Entities.CurrentEntityID)
		if ok {
			reverseIDList(&list)
			dispatcher.Dispatch(&actions.SetEntityProcessList{ProcessList: list})
		}
		newMap, ok := api.GetProcessEnvelopeCountMap()
		if ok {
			dispatcher.Dispatch(&actions.SetEnvelopeHeights{EnvelopeHeights: newMap})
		}
		update.EntityProcessResults()
	}
}

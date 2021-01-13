package components

import (
	"fmt"
	"time"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/vocdoni/vocexplorer/api"
	"github.com/vocdoni/vocexplorer/api/dbtypes"
	"github.com/vocdoni/vocexplorer/frontend/actions"
	"github.com/vocdoni/vocexplorer/frontend/bootstrap"
	"github.com/vocdoni/vocexplorer/frontend/dispatcher"
	"github.com/vocdoni/vocexplorer/frontend/store"
	"github.com/vocdoni/vocexplorer/frontend/update"
	"github.com/vocdoni/vocexplorer/logger"
	"github.com/vocdoni/vocexplorer/util"
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

	return Container(
		vecty.Markup(vecty.Attribute("id", "main")),
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
				vecty.If(store.Entities.CurrentEntity.ProcessCount > 0,
					bootstrap.Card(bootstrap.CardParams{
						Body: &EntityProcessListView{}})),
				vecty.If(store.Entities.CurrentEntity.ProcessCount == 0, bootstrap.Card(bootstrap.CardParams{
					Body: vecty.Text("This entity has no processes")})),
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
			vecty.Markup(
				vecty.Attribute("href", "https://vocdoni.link/entities/0x"+store.Entities.CurrentEntityID),
				vecty.Property("target", store.Entities.CurrentEntityID),
			),
			vecty.Markup(vecty.Attribute("aria-label", "Link to entity "+store.Entities.CurrentEntityID+"'s profile page")),
			vecty.Text("Entity Profile"),
		),
	}
}

// UpdateEntityContents keeps the dashboard data up to date
func UpdateEntityContents(d *EntityContentsView) {
	// Set entity process list to nil so previous list is not displayed
	dispatcher.Dispatch(&actions.SetEntityProcessList{ProcessList: [10]*dbtypes.Process{}})
	dispatcher.Dispatch(&actions.EnableAllUpdates{})
	ticker := time.NewTicker(time.Duration(store.Config.RefreshTime) * time.Second)
	dispatcher.Dispatch(&actions.ServerConnected{Connected: api.PingServer()})

	newCount, ok := api.GetEntityProcessCount(store.Entities.CurrentEntityID)
	if ok {
		dispatcher.Dispatch(&actions.SetEntityProcessCount{Count: int(newCount)})
	}
	if !update.CheckCurrentPage("entity", ticker) {
		return
	}
	updateEntityProcesses(d, util.Max(store.Entities.CurrentEntity.ProcessCount-store.Entities.ProcessPagination.Index, 1))
	for {
		select {
		case <-store.RedirectChan:
			if !update.CheckCurrentPage("entity", ticker) {
				return
			}
		case <-ticker.C:
			if !update.CheckCurrentPage("entity", ticker) {
				return
			}
			updateEntityProcesses(d, util.Max(store.Entities.CurrentEntity.ProcessCount-store.Entities.ProcessPagination.Index, 1))
		case i := <-store.Entities.ProcessPagination.PagChannel:
			if !update.CheckCurrentPage("entity", ticker) {
				return
			}
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
			if i < 1 {
				newCount, _ := api.GetEntityProcessCount(store.Entities.CurrentEntityID)
				dispatcher.Dispatch(&actions.SetEntityProcessCount{Count: int(newCount)})
			}
			index := util.Max(store.Entities.CurrentEntity.ProcessCount-store.Entities.ProcessPagination.Index, 1)
			logger.Info(fmt.Sprintf("Getting processes from entity %s, index %d\n", store.Entities.CurrentEntityID, index))
			list, ok := api.GetProcessListByEntity(index-1, store.Entities.CurrentEntityID)
			if ok {
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
	dispatcher.Dispatch(&actions.ServerConnected{Connected: api.PingServer()})

	if store.Entities.CurrentEntity.ProcessCount > 0 && !store.Entities.ProcessPagination.DisableUpdate {
		newCount, ok := api.GetEntityProcessCount(store.Entities.CurrentEntityID)
		if ok {
			dispatcher.Dispatch(&actions.SetEntityProcessCount{Count: int(newCount)})
		}
		logger.Info(fmt.Sprintf("Getting processes from entity %s, index %d\n", store.Entities.CurrentEntityID, index))
		list, ok := api.GetProcessListByEntity(index, store.Entities.CurrentEntityID)
		if ok {
			dispatcher.Dispatch(&actions.SetEntityProcessList{ProcessList: list})
		}
		newMap, ok := api.GetProcessEnvelopeCountMap()
		if ok {
			dispatcher.Dispatch(&actions.SetEnvelopeHeights{EnvelopeHeights: newMap})
		}
		update.EntityProcessResults()
	}
}

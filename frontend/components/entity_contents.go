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
	"gitlab.com/vocdoni/vocexplorer/frontend/store/storeutil"
	"gitlab.com/vocdoni/vocexplorer/frontend/update"
	"gitlab.com/vocdoni/vocexplorer/logger"
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
		elem.Div(
			elem.Span(
				vecty.Markup(vecty.Class("title")),
				elem.Anchor(
					vecty.Markup(
						vecty.Attribute("href", store.EntityURL+store.Entities.CurrentEntityID),
						vecty.Property("target", store.EntityURL+store.Entities.CurrentEntityID),
					),
					vecty.Markup(vecty.Attribute("aria-label", "Link to entity "+store.EntityURL+store.Entities.CurrentEntityID+"'s public profile")),
					vecty.Text("Entity Profile"),
				),
			),
		),
	}
}

// UpdateEntityContents keeps the dashboard data up to date
func UpdateEntityContents(d *EntityContentsView) {
	// Set entity process list to nil so previous list is not displayed
	dispatcher.Dispatch(&actions.SetEntityProcessIds{ProcessList: []string{}})
	dispatcher.Dispatch(&actions.EnableAllUpdates{})
	ticker := time.NewTicker(time.Duration(store.Config.RefreshTime) * 5 * time.Second)
	dispatcher.Dispatch(&actions.GatewayConnected{GatewayErr: store.Client.GetGatewayInfo()})

	newCount, err := store.Client.GetProcessCount(util.StringToHex(store.Entities.CurrentEntityID))
	if err != nil {
		logger.Error(err)
	} else {
		dispatcher.Dispatch(&actions.SetCurrentEntityProcessCount{Count: int(newCount)})
	}
	if !update.CheckCurrentPage("entity", ticker) {
		return
	}
	updateEntityProcesses(d, store.Entities.CurrentEntity.ProcessCount-store.Entities.ProcessPagination.Index-config.ListSize)
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
			updateEntityProcesses(d, store.Entities.CurrentEntity.ProcessCount-store.Entities.ProcessPagination.Index-config.ListSize)
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
			eid := util.StringToHex(store.Entities.CurrentEntityID)
			index := util.Max(store.Entities.CurrentEntity.ProcessCount-store.Entities.ProcessPagination.Index, 1)
			logger.Info(fmt.Sprintf("Getting processes from entity %s, index %d\n", store.Entities.CurrentEntityID, index))
			list, err := store.Client.GetProcessList(eid, "", 0, "", false, "", index-1, config.ListSize)
			if err != nil {
				logger.Error(err)
			} else {
				dispatcher.Dispatch(&actions.SetEntityProcessIds{ProcessList: list})
			}
		}
	}
}

func updateEntityProcesses(d *EntityContentsView, index int) {
	newCount, err := store.Client.GetProcessCount(util.StringToHex(store.Entities.CurrentEntityID))
	if err != nil {
		logger.Error(err)
	} else {
		dispatcher.Dispatch(&actions.SetCurrentEntityProcessCount{Count: int(newCount)})
	}

	if store.Entities.CurrentEntity.ProcessCount > 0 && !store.Entities.ProcessPagination.DisableUpdate {
		listSize := config.ListSize
		if index < 0 {
			listSize += index
			index = 0
		}
		logger.Info(fmt.Sprintf("Getting %d processes from index %d\n", listSize, index))
		list, err := store.Client.GetProcessList(util.StringToHex(store.Entities.CurrentEntityID), "", 0, "", false, "", index, listSize)
		if err != nil {
			logger.Error(err)
			return
		}
		reverseIDList(list)
		dispatcher.Dispatch(&actions.SetEntityProcessIds{ProcessList: list})
		for _, processId := range store.Entities.CurrentEntity.ProcessIds {
			summary, err := store.Client.GetProcessSummary(util.StringToHex(processId))
			if err != nil {
				logger.Error(err)
			}
			if summary == nil {
				return
			}
			dispatcher.Dispatch(&actions.SetProcess{
				PID: processId,
				Process: &storeutil.Process{
					EnvelopeCount:  int(*summary.EnvelopeHeight),
					ProcessSummary: summary,
					ProcessID:      processId,
				},
			})
		}
	}
	dispatcher.Dispatch(&actions.GatewayConnected{GatewayErr: store.Client.GetGatewayInfo()})
}

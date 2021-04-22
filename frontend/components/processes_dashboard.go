package components

import (
	"fmt"
	"time"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/vocdoni/vocexplorer/api/dbtypes"

	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/frontend/update"
	"gitlab.com/vocdoni/vocexplorer/logger"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// ProcessesDashboardView renders the processes dashboard page
type ProcessesDashboardView struct {
	vecty.Core
	vecty.Mounter
	Rendered bool
}

// Mount is called after the component renders to signal that it can be rerendered safely
func (dash *ProcessesDashboardView) Mount() {
	if !dash.Rendered {
		dash.Rendered = true
		vecty.Rerender(dash)
	}
}

// Render renders the ProcessesDashboardView component
func (dash *ProcessesDashboardView) Render() vecty.ComponentOrHTML {
	if !dash.Rendered {
		return LoadingBar()
	}
	return Container(
		vecty.Markup(vecty.Attribute("id", "main")),
		renderServerConnectionBanner(),
		elem.Section(
			bootstrap.Card(bootstrap.CardParams{
				Body: vecty.List{
					elem.Heading1(vecty.Text("Processes")),
					&ProcessListView{},
				},
			}),
		),
	)
}

// UpdateProcessesDashboard continuously updates the information needed by the Processes dashboard
func UpdateProcessesDashboard(d *ProcessesDashboardView) {
	dispatcher.Dispatch(&actions.EnableAllUpdates{})

	ticker := time.NewTicker(time.Duration(store.Config.RefreshTime) * time.Second)
	updateProcesses(d)
	for {

		select {
		case <-store.RedirectChan:
			if !update.CheckCurrentPage("processes", ticker) {
				return
			}
		case <-ticker.C:
			if !update.CheckCurrentPage("processes", ticker) {
				return
			}
			updateProcesses(d)
		case i := <-store.Processes.Pagination.PagChannel:
			if !update.CheckCurrentPage("processes", ticker) {
				return
			}
		processLoop:
			for {
				// If many indices waiting in buffer, scan to last one.
				select {
				case i = <-store.Processes.Pagination.PagChannel:
				default:
					break processLoop

				}
			}
			dispatcher.Dispatch(&actions.ProcessesIndexChange{Index: i})
			if i < 1 {
				newVal, ok := store.Client.GetProcessCount()
				if ok {
					dispatcher.Dispatch(&actions.SetProcessCount{Count: int(newVal)})
				}
			}
			if store.Processes.Count > 0 {
				getProcesses(d, util.Max(store.Processes.Count-store.Processes.Pagination.Index, 1))
				update.ProcessResults()
			}
		case search := <-store.Processes.Pagination.SearchChannel:
			if !update.CheckCurrentPage("processes", ticker) {
				return
			}
		processSearch:
			for {
				// If many indices waiting in buffer, scan to last one.
				select {
				case search = <-store.Processes.Pagination.SearchChannel:
				default:
					break processSearch
				}
			}
			logger.Info("search: " + search)
			dispatcher.Dispatch(&actions.ProcessesIndexChange{Index: 0})
			list, ok := store.Client.GetProcessList(search)
			if ok {
				dispatcher.Dispatch(&actions.SetProcessList{Processes: list})
			} else {
				dispatcher.Dispatch(&actions.SetProcessList{Processes: [config.ListSize]*dbtypes.Process{}})
			}
			update.ProcessResults()
		}
	}
}

func updateProcesses(d *ProcessesDashboardView) {
	dispatcher.Dispatch(&actions.GatewayConnected{GatewayErr: store.Client.GetGatewayInfo()})
	if !store.Processes.Pagination.DisableUpdate {
		stats, err := store.Client.GetStats()
		if err != nil {
			logger.Error(err)
			return
		}
		actions.UpdateCounts(stats)
		getProcesses(d, util.Max(store.Processes.Count-store.Processes.Pagination.Index, 1))
		update.ProcessResults()
	}
}

func getProcesses(d *ProcessesDashboardView, index int) {
	// index--
	logger.Info(fmt.Sprintf("Getting processes from index %d\n", index))
	list, ok := store.Client.GetProcessList(index)
	if ok {
		dispatcher.Dispatch(&actions.SetProcessList{Processes: list})
	}
	newVal, ok := store.Client.GetProcessEnvelopeCountMap()
	if ok {
		dispatcher.Dispatch(&actions.SetEnvelopeHeights{EnvelopeHeights: newVal})
	}
	newVal, ok = store.Client.GetEntityProcessCountMap()
	if ok {
		dispatcher.Dispatch(&actions.SetProcessHeights{ProcessHeights: newVal})
	}
}

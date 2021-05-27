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

	ticker := time.NewTicker(time.Duration(store.Config.RefreshTime) * 5 * time.Second)
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
			if store.Processes.Count > 0 {
				getProcesses(d, store.Processes.Count-store.Processes.Pagination.Index-config.ListSize)
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
			list, err := store.Client.GetProcessList([]byte{}, search, 0, "", false, 0, config.ListSize)
			if err != nil {
				dispatcher.Dispatch(&actions.SetProcessIds{Processes: []string{}})
				logger.Error(err)
			} else {
				fetchProcessMetas(list)
			}
		}
	}
}

func updateProcesses(d *ProcessesDashboardView) {
	if !store.Processes.Pagination.DisableUpdate {
		stats, err := store.Client.GetStats()
		if err != nil {
			logger.Error(err)
			return
		}
		actions.UpdateCounts(stats)
		getProcesses(d, store.Processes.Count-store.Processes.Pagination.Index-config.ListSize)
	}
	dispatcher.Dispatch(&actions.GatewayConnected{GatewayErr: store.Client.GetGatewayInfo()})
}

func getProcesses(d *ProcessesDashboardView, index int) {
	listSize := config.ListSize
	if index < 0 {
		listSize += index
		index = 0
	}
	logger.Info(fmt.Sprintf("Getting %d processes from index %d\n", listSize, index))
	list, err := store.Client.GetProcessList([]byte{}, "", 0, "", false, index, listSize)
	if err != nil {
		logger.Error(err)
		return
	}
	fetchProcessMetas(list)
}

func fetchProcessMetas(list []string) {
	reverseIDList(list)
	dispatcher.Dispatch(&actions.SetProcessIds{Processes: list})
	for _, processId := range store.Processes.ProcessIds {
		if processId == "" {
			break
		}
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

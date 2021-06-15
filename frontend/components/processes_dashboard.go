package components

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/event"
	"github.com/hexops/vecty/prop"

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
	dispatcher.Dispatch(&actions.SetProcessStatusFilter{})
	dispatcher.Dispatch(&actions.SetProcessSrcNetworkIDFilter{})
	dispatcher.Dispatch(&actions.SetProcessResultsFilter{})
	dispatcher.Dispatch(&actions.SetProcessNamespaceFilter{})

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
			logger.Info(fmt.Sprintf("search: %s, %d, %s, %s, %t", search, store.Processes.NamespaceFilter, store.Processes.StatusFilter, store.Processes.SrcNetworkIDFilter, store.Processes.ResultsFilter))
			dispatcher.Dispatch(&actions.ProcessesIndexChange{Index: 0})
			list, err := store.Client.GetProcessList([]byte{}, search, uint32(store.Processes.NamespaceFilter), store.Processes.StatusFilter, store.Processes.ResultsFilter, store.Processes.SrcNetworkIDFilter, 0, config.ListSize)
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
	list, err := store.Client.GetProcessList([]byte{}, "", 0, "", false, "", index, listSize)
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

// dropdown menus for pagination

func generateNamespaceDropdown() vecty.ComponentOrHTML {
	options := []vecty.MarkupOrChild{}
	options = append(options, elem.Option(vecty.Text("")))
	for i := 0; i < 500; i++ {
		options = append(options, elem.Option(vecty.Text(strconv.Itoa(i))))
	}
	return elem.Div(
		vecty.Markup(vecty.Class("dropdown")),
		elem.Div(
			vecty.Markup(vecty.Class("description")),
			vecty.Text("namespace"),
		),
		elem.Div(
			vecty.Markup(
				event.Change(
					func(e *vecty.Event) {
						filter, err := strconv.Atoi(e.Target.Get("value").String())
						if err != nil {
							logger.Error(fmt.Errorf("cannot parse namespace filter" + e.Target.Get("value").String()))
						}
						dispatcher.Dispatch(&actions.SetProcessNamespaceFilter{NamespaceFilter: filter})
					},
				),
			),
			vecty.Markup(vecty.Class("contents")),
			elem.Select(
				options...,
			),
		),
	)
}

func generateStatusDropdown() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(vecty.Class("dropdown")),
		elem.Div(
			vecty.Markup(vecty.Class("description")),
			vecty.Text("status"),
		),
		elem.Div(
			vecty.Markup(
				event.Change(
					func(e *vecty.Event) {
						filter := e.Target.Get("value").String()
						if filter == "unknown" {
							filter = "process_unknown"
						}
						dispatcher.Dispatch(&actions.SetProcessStatusFilter{StatusFilter: filter})
					},
				),
			),
			vecty.Markup(vecty.Class("contents")),
			elem.Select(
				elem.Option(vecty.Text("")),
				elem.Option(vecty.Text("ready")),
				elem.Option(vecty.Text("ended")),
				elem.Option(vecty.Text("canceled")),
				elem.Option(vecty.Text("paused")),
				elem.Option(vecty.Text("results")),
				elem.Option(vecty.Text("unknown")),
			),
		),
	)
}

func generateSourceNetworkIDDropdown() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(vecty.Class("dropdown")),
		elem.Div(
			vecty.Markup(vecty.Class("description")),
			vecty.Text("source network id"),
		),
		elem.Div(
			vecty.Markup(
				event.Change(
					func(e *vecty.Event) {
						filter := e.Target.Get("value").String()
						dispatcher.Dispatch(&actions.SetProcessSrcNetworkIDFilter{SrcNetworkIDFilter: filter})
					},
				),
			),
			vecty.Markup(vecty.Class("contents")),
			elem.Select(
				elem.Option(vecty.Text("")),
				elem.Option(vecty.Text("unknown")),
				elem.Option(vecty.Text("eth mainnet")),
				elem.Option(vecty.Text("eth rinkeby")),
				elem.Option(vecty.Text("eth goerli")),
				elem.Option(vecty.Text("poa xdai")),
				elem.Option(vecty.Text("poa sokol")),
				elem.Option(vecty.Text("polygon")),
				elem.Option(vecty.Text("bcd")),
				elem.Option(vecty.Text("eth mainnet signaling")),
				elem.Option(vecty.Text("eth rinkeby signaling")),
			),
		),
	)
}

func generateResultsCheckbox() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(vecty.Class("dropdown")),
		elem.Div(
			vecty.Markup(vecty.Class("description")),
			vecty.Text("has results"),
		),
		elem.Div(
			vecty.Markup(
				event.Change(
					func(e *vecty.Event) {
						filter := e.Target.Get("checked").Bool()
						dispatcher.Dispatch(&actions.SetProcessResultsFilter{ResultsFilter: filter})
					},
				),
			),
			vecty.Markup(vecty.Class("contents")),
			elem.Input(
				vecty.Markup(
					prop.Value("checked_value"),
					prop.Type("checkbox"),
				),
				vecty.Text("results available"),
			),
		),
	)
}

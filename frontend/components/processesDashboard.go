package components

import (
	"fmt"
	"time"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/go-dvote/log"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/api"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/update"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// ProcessesDashboardView renders the processes dashboard page
type ProcessesDashboardView struct {
	vecty.Core
	envelopeIndex int
}

// Render renders the ProcessesDashboardView component
func (dash *ProcessesDashboardView) Render() vecty.ComponentOrHTML {
	if dash == nil || store.GatewayClient == nil {
		return &bootstrap.Alert{
			Contents: "Connecting to blockchain clients",
			Type:     "warning",
		}
	}

	if store.Processes.CurrentProcess.ProcessType == "" {
		dispatcher.Dispatch(&actions.SetProcessType{Type: "unknown"})
	}
	if store.Processes.CurrentProcess.State == "" {
		dispatcher.Dispatch(&actions.SetProcessState{State: "unknown"})
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
						Body: dash.ProcessDetails(),
					}),
				),
			),
		),
		elem.Section(
			vecty.Markup(vecty.Class("row")),
			elem.Div(
				vecty.Markup(vecty.Class("col-12")),
				bootstrap.Card(bootstrap.CardParams{
					Body: dash.ProcessTabs(),
				}),
			),
		),
	)
}

//ProcessDetails renders the details of a process
func (dash *ProcessesDashboardView) ProcessDetails() vecty.List {
	return vecty.List{
		elem.Heading1(
			vecty.Text("Process details"),
		),
		elem.Heading2(vecty.Text(store.Processes.CurrentProcessID)),
		elem.Div(
			vecty.Markup(vecty.Class("badges")),
			elem.Span(
				vecty.Markup(vecty.Class("badge", store.Processes.CurrentProcess.State)),
				vecty.Text(store.Processes.CurrentProcess.State),
			),
		),
		elem.HorizontalRule(),
		elem.DescriptionList(
			elem.DefinitionTerm(vecty.Text("Process type")),
			elem.Description(vecty.Text(store.Processes.CurrentProcess.ProcessType)),
			elem.DefinitionTerm(vecty.Text("State")),
			elem.Description(vecty.Text(store.Processes.CurrentProcess.State)),
			elem.DefinitionTerm(vecty.Text("Registered votes")),
			elem.Description(vecty.Text(util.IntToString(store.Processes.CurrentProcess.EnvelopeCount))),
		),
	}
}

//ProcessTab is a single tab of a process page
type ProcessTab struct {
	*Tab
}

func (p *ProcessTab) dispatch() interface{} {
	return &actions.ProcessesTabChange{
		Tab: p.alias(),
	}
}

func (p *ProcessTab) store() string {
	return store.Processes.Pagination.Tab
}

//ProcessTabs renders the tabs for a process
func (dash *ProcessesDashboardView) ProcessTabs() vecty.List {
	results := &ProcessTab{&Tab{
		Text:  "Results",
		Alias: "results",
	}}
	envelopes := &ProcessTab{&Tab{
		Text:  "Envelopes",
		Alias: "envelopes",
	}}

	return vecty.List{
		elem.Navigation(
			vecty.Markup(vecty.Class("tabs")),
			elem.UnorderedList(
				TabLink(dash, results),
				TabLink(dash, envelopes),
			),
		),
		elem.Div(
			vecty.Markup(vecty.Class("tabs-content")),
			TabContents(results, renderResults(store.Processes.CurrentProcess.Results)),
			TabContents(envelopes, &ProcessesEnvelopeListView{}),
		),
	}
}

func renderResults(results [][]uint32) vecty.ComponentOrHTML {
	if len(results) <= 0 {
		return elem.Preformatted(
			vecty.Markup(vecty.Class("empty")),
			vecty.Text("No results yet"),
		)
	}
	var resultList []vecty.MarkupOrChild
	var header []vecty.MarkupOrChild
	header = append(header, elem.TableHeader())
	numCols := 0
	for i, row := range results {
		var resultRow []vecty.MarkupOrChild
		resultRow = append(resultRow, elem.TableHeader(vecty.Text("Question "+util.IntToString(i)+": ")))
		for _, val := range row {
			resultRow = append(resultRow, elem.TableData(vecty.Text(util.IntToString(val)+" ")))
		}
		resultList = append(resultList, elem.TableRow(resultRow...))
		numCols = util.Max(numCols, len(row))
	}
	for i := 0; i < numCols; i++ {
		header = append(header, elem.TableHeader(vecty.Text("Option "+util.IntToString(i)+": ")))
	}
	resultList = append(resultList, elem.TableHead(
		elem.TableRow(header...),
	))
	return elem.Div(
		elem.Heading5(vecty.Text("Process Results: ")),
		elem.Table(resultList...),
	)
}

// UpdateAndRenderProcessesDashboard keeps the data for the processes dashboard up-to-date
func UpdateAndRenderProcessesDashboard(d *ProcessesDashboardView) {
	ticker := time.NewTicker(time.Duration(store.Config.RefreshTime) * time.Second)
	updateProcessesDashboard(d)
	for {
		select {
		case <-store.RedirectChan:
			fmt.Println("Redirecting...")
			ticker.Stop()
			return
		case <-ticker.C:
			updateProcessesDashboard(d)
		case i := <-store.Processes.Pagination.PagChannel:
		loop:
			for {
				// If many indices waiting in buffer, scan to last one.
				select {
				case i = <-store.Processes.Pagination.PagChannel:
				default:
					break loop
				}
			}
			d.envelopeIndex = i
			oldEnvelopes := store.Processes.CurrentProcess.EnvelopeCount
			newVal, ok := api.GetProcessEnvelopeHeight(store.Processes.CurrentProcessID)
			if ok {
				dispatcher.Dispatch(&actions.SetCurrentProcessEnvelopeHeight{Height: int(newVal)})
			}
			if i < 1 {
				oldEnvelopes = store.Processes.CurrentProcess.EnvelopeCount
			}
			if store.Processes.CurrentProcess.EnvelopeCount > 0 {
				updateProcessEnvelopes(d, util.Max(oldEnvelopes-d.envelopeIndex, config.ListSize))
			}
		}
	}
}

func updateProcessesDashboard(d *ProcessesDashboardView) {
	dispatcher.Dispatch(&actions.GatewayConnected{Connected: api.PingGateway(store.Config.GatewayHost)})

	dispatcher.Dispatch(&actions.ServerConnected{Connected: api.Ping()})

	update.CurrentProcessResults()
	newVal, ok := api.GetProcessEnvelopeHeight(store.Processes.CurrentProcessID)
	if ok {
		dispatcher.Dispatch(&actions.SetCurrentProcessEnvelopeHeight{Height: int(newVal)})
	}
	if !store.Envelopes.Pagination.DisableUpdate && store.Processes.CurrentProcess.EnvelopeCount > 0 {
		updateProcessEnvelopes(d, util.Max(store.Processes.CurrentProcess.EnvelopeCount-d.envelopeIndex, config.ListSize))
	}
}

func updateProcessEnvelopes(d *ProcessesDashboardView, index int) {
	log.Infof("Getting envelopes from index %d", util.IntToString(index))
	list, ok := api.GetEnvelopeListByProcess(index, store.Processes.CurrentProcessID)
	if ok {
		reverseEnvelopeList(&list)
		dispatcher.Dispatch(&actions.SetEnvelopeList{EnvelopeList: list})
	}
}

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

// ProcessContentsView renders the processes dashboard page
type ProcessContentsView struct {
	vecty.Core
	vecty.Mounter
	Rendered bool
}

// Mount is called after the component renders to signal that it can be rerendered safely
func (dash *ProcessContentsView) Mount() {
	if !dash.Rendered {
		dash.Rendered = true
		vecty.Rerender(dash)
	}
}

// Render renders the ProcessContentsView component
func (dash *ProcessContentsView) Render() vecty.ComponentOrHTML {
	if !dash.Rendered {
		return LoadingBar()
	}
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
func (dash *ProcessContentsView) ProcessDetails() vecty.List {
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
	return &actions.ProcessTabChange{
		Tab: p.alias(),
	}
}

func (p *ProcessTab) store() string {
	return store.Processes.Pagination.Tab
}

//ProcessTabs renders the tabs for a process
func (dash *ProcessContentsView) ProcessTabs() vecty.List {
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

// UpdateProcessContents keeps the data for the processes dashboard up-to-date
func UpdateProcessContents(d *ProcessContentsView) {
	dispatcher.Dispatch(&actions.EnableAllUpdates{})
	ticker := time.NewTicker(time.Duration(store.Config.RefreshTime) * time.Second)
	updateProcessContents(d)
	for {
		select {
		case <-store.RedirectChan:
			fmt.Println("Redirecting...")
			ticker.Stop()
			return
		case <-ticker.C:
			updateProcessContents(d)
		case i := <-store.Processes.EnvelopePagination.PagChannel:
		loop:
			for {
				// If many indices waiting in buffer, scan to last one.
				select {
				case i = <-store.Processes.EnvelopePagination.PagChannel:
				default:
					break loop
				}
			}
			dispatcher.Dispatch(&actions.ProcessEnvelopesIndexChange{Index: i})
			oldEnvelopes := store.Processes.CurrentProcess.EnvelopeCount
			newVal, ok := api.GetProcessEnvelopeCount(store.Processes.CurrentProcessID)
			if ok {
				dispatcher.Dispatch(&actions.SetCurrentProcessEnvelopeHeight{Height: int(newVal)})
			}
			if i < 1 {
				oldEnvelopes = store.Processes.CurrentProcess.EnvelopeCount
			}
			if store.Processes.CurrentProcess.EnvelopeCount > 0 {
				updateProcessEnvelopes(d, util.Max(oldEnvelopes-store.Processes.EnvelopePagination.Index, 1))
			}
		}
	}
}

func updateProcessContents(d *ProcessContentsView) {
	dispatcher.Dispatch(&actions.GatewayConnected{Connected: store.GatewayClient.Ping()})
	dispatcher.Dispatch(&actions.ServerConnected{Connected: api.PingServer()})
	update.CurrentProcessResults()
	newVal, ok := api.GetProcessEnvelopeCount(store.Processes.CurrentProcessID)
	if ok {
		dispatcher.Dispatch(&actions.SetCurrentProcessEnvelopeHeight{Height: int(newVal)})
	}
	if !store.Envelopes.Pagination.DisableUpdate && store.Processes.CurrentProcess.EnvelopeCount > 0 {
		updateProcessEnvelopes(d, util.Max(store.Processes.CurrentProcess.EnvelopeCount-store.Processes.EnvelopePagination.Index, 1))
	}
}

func updateProcessEnvelopes(d *ProcessContentsView, index int) {
	fmt.Printf("Getting envelopes from index %d\n", index)
	list, ok := api.GetEnvelopeListByProcess(index, store.Processes.CurrentProcessID)
	if ok {
		reverseEnvelopeList(&list)
		dispatcher.Dispatch(&actions.SetCurrentProcessEnvelopes{EnvelopeList: list})
	}
}

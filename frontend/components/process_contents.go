package components

import (
	"fmt"
	"time"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/api"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
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

	if store.Processes.CurrentProcessResults.ProcessType == "" {
		dispatcher.Dispatch(&actions.SetProcessType{Type: "unknown"})
	}
	if store.Processes.CurrentProcessResults.State == "" {
		dispatcher.Dispatch(&actions.SetProcessState{State: "unknown"})
	}

	return Container(
		renderGatewayConnectionBanner(),
		renderServerConnectionBanner(),
		DetailsView(
			dash.ProcessDetails(),
			dash.ProcessTabs(),
		),
	)
}

//ProcessDetails renders the details of a process
func (dash *ProcessContentsView) ProcessDetails() vecty.List {
	return vecty.List{
		elem.Heading1(
			vecty.Text("Process details"),
		),
		elem.Heading2(vecty.Text(store.Processes.CurrentProcess.ID)),
		elem.Div(
			vecty.Markup(vecty.Class("badges")),
			elem.Span(
				vecty.Markup(vecty.Class("badge", store.Processes.CurrentProcessResults.State)),
				vecty.Text(store.Processes.CurrentProcessResults.State),
			),
		),
		elem.HorizontalRule(),
		elem.DescriptionList(
			elem.DefinitionTerm(vecty.Text("Process type")),
			elem.Description(vecty.Text(store.Processes.CurrentProcessResults.ProcessType)),
			elem.DefinitionTerm(vecty.Text("State")),
			elem.Description(vecty.Text(store.Processes.CurrentProcessResults.State)),
			elem.DefinitionTerm(vecty.Text("Registered votes")),
			elem.Description(vecty.Text(util.IntToString(store.Processes.CurrentProcessResults.EnvelopeCount))),
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
			TabContents(results, renderResults(store.Processes.CurrentProcessResults.Results)),
			TabContents(envelopes, &ProcessesEnvelopeListView{}),
		),
	}
}

func renderPollAnswers(answers []uint32) vecty.ComponentOrHTML {
	items := vecty.List{}
	for _, a := range answers {
		items = append(items, elem.ListItem(
			vecty.Text(fmt.Sprintf("%d", a)),
		))
	}

	return items
}

func renderResults(results [][]uint32) vecty.ComponentOrHTML {
	if len(results) <= 0 {
		return elem.Preformatted(
			vecty.Markup(vecty.Class("empty")),
			vecty.Text("No results yet"),
		)
	}

	content := vecty.List{}

	for i, row := range results {
		res := elem.OrderedList(
			renderPollAnswers(row),
		)
		content = append(content, elem.Div(
			elem.Span(
				vecty.Markup(vecty.Class("question")),
				vecty.Text(fmt.Sprintf("Question %d", i+1)),
			),
			res,
		))
	}

	return elem.Div(
		vecty.Markup(vecty.Class("poll-results")),
		content,
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
			oldEnvelopes := store.Processes.CurrentProcessResults.EnvelopeCount
			newVal, ok := api.GetProcessEnvelopeCount(store.Processes.CurrentProcess.ID)
			if ok {
				dispatcher.Dispatch(&actions.SetCurrentProcessEnvelopeHeight{Height: int(newVal)})
			}
			if i < 1 {
				oldEnvelopes = store.Processes.CurrentProcessResults.EnvelopeCount
			}
			if store.Processes.CurrentProcessResults.EnvelopeCount > 0 {
				updateProcessEnvelopes(d, util.Max(oldEnvelopes-store.Processes.EnvelopePagination.Index, 1))
			}
		}
	}
}

func updateProcessContents(d *ProcessContentsView) {
	go dispatcher.Dispatch(&actions.GatewayConnected{Connected: store.GatewayClient.Ping()})
	go dispatcher.Dispatch(&actions.ServerConnected{Connected: api.PingServer()})
	update.CurrentProcessResults()
	newVal, ok := api.GetProcessEnvelopeCount(store.Processes.CurrentProcess.ID)
	if ok {
		dispatcher.Dispatch(&actions.SetCurrentProcessEnvelopeHeight{Height: int(newVal)})
	}
	if !store.Envelopes.Pagination.DisableUpdate && store.Processes.CurrentProcessResults.EnvelopeCount > 0 {
		updateProcessEnvelopes(d, util.Max(store.Processes.CurrentProcessResults.EnvelopeCount-store.Processes.EnvelopePagination.Index, 1))
	}
}

func updateProcessEnvelopes(d *ProcessContentsView, index int) {
	fmt.Printf("Getting envelopes from index %d\n", index)
	list, ok := api.GetEnvelopeListByProcess(index, store.Processes.CurrentProcess.ID)
	if ok {
		reverseEnvelopeList(&list)
		dispatcher.Dispatch(&actions.SetCurrentProcessEnvelopes{EnvelopeList: list})
	}
}

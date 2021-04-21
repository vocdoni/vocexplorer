package components

import (
	"fmt"
	"strings"
	"time"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/client"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/frontend/update"
	"gitlab.com/vocdoni/vocexplorer/logger"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// ProcessContentsView renders the processes dashboard page
type ProcessContentsView struct {
	vecty.Core
	vecty.Mounter
	Rendered    bool
	Unavailable bool
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
	if dash.Unavailable {
		return Unavailable("Process unavailable", "")
	}
	if store.Processes.CurrentProcess == nil || store.Processes.CurrentProcess.EntityID == "" {
		return Unavailable("Loading process...", "")
	}

	if store.Processes.CurrentProcessResults.ProcessInfo.Type == "" {
		dispatcher.Dispatch(&actions.SetProcessType{Type: "Unknown"})
	}
	if store.Processes.CurrentProcessResults.ProcessInfo.State == "" {
		dispatcher.Dispatch(&actions.SetProcessState{State: "Unknown"})
	}

	return Container(
		vecty.Markup(vecty.Attribute("id", "main")),
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
			elem.Span(
				vecty.Markup(vecty.Class("title")),
				elem.Anchor(
					vecty.Markup(
						vecty.Attribute("href", store.ProcessDomain+store.Processes.CurrentProcess.EntityID+"/0x"+store.Processes.CurrentProcess.ID),
						vecty.Property("target", store.Processes.CurrentProcess.ID),
					),
					vecty.Markup(vecty.Attribute("aria-label", "Link to process "+store.Processes.CurrentProcess.ID+"'s profile page")),
					vecty.Text("Process Profile"),
				),
			),
		),
		elem.Div(
			vecty.Markup(vecty.Class("badges")),
			elem.Span(
				vecty.Markup(vecty.Class("badge", store.Processes.CurrentProcessResults.ProcessInfo.State)),
				vecty.Text(strings.Title(store.Processes.CurrentProcessResults.ProcessInfo.State)),
			),
		),
		elem.HorizontalRule(),
		elem.DescriptionList(
			elem.DefinitionTerm(vecty.Text("Host entity")),
			elem.Description(
				Link(
					"/entity/"+store.Processes.CurrentProcess.EntityID,
					store.Processes.CurrentProcess.EntityID,
					"",
				),
			),
			elem.DefinitionTerm(vecty.Text("Process type")),
			elem.Description(vecty.Text(util.GetProcessName(store.Processes.CurrentProcessResults.ProcessInfo.Type))),
			elem.DefinitionTerm(vecty.Text("State")),
			elem.Description(vecty.Text(strings.Title(store.Processes.CurrentProcessResults.ProcessInfo.State))),
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
	processDetails := &ProcessTab{&Tab{
		Text:  "Details",
		Alias: "details",
	}}

	return vecty.List{
		elem.Navigation(
			vecty.Markup(vecty.Attribute("aria-label", "Tab navigation: results, envelopes and details")),
			vecty.Markup(vecty.Class("tabs")),
			elem.UnorderedList(
				TabLink(dash, results),
				TabLink(dash, envelopes),
				TabLink(dash, processDetails),
			),
		),
		elem.Div(
			vecty.Markup(vecty.Class("tabs-content")),
			TabContents(results, renderResults(store.Processes.CurrentProcessResults.ProcessInfo.Results)),
			TabContents(envelopes, &ProcessesEnvelopeListView{}),
			TabContents(processDetails, renderProcessDetails(store.Processes.CurrentProcessResults.ProcessInfo)),
		),
	}
}

func renderPollAnswers(answers []string) vecty.ComponentOrHTML {
	items := vecty.List{}
	for _, a := range answers {
		items = append(items, elem.ListItem(
			vecty.Text(fmt.Sprintf("%s", a)),
		))
	}

	return items
}

func renderResults(results [][]string) vecty.ComponentOrHTML {
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
				vecty.Text(fmt.Sprintf("Field %d", i+1)),
			),
			res,
		))
	}

	return elem.Div(
		vecty.Markup(vecty.Class("poll-results")),
		content,
	)
}

func renderProcessDetails(details api.ProcessResults) vecty.ComponentOrHTML {
	row1 := vecty.List{}
	row2 := vecty.List{}

	// Add EnvelopeType
	row1 = append(row1, renderEnvelopeType(details.EnvelopeType))
	row1 = append(row1, renderProcessMode(details.Mode))
	row1 = append(row1, renderProcessVoteOptions(details.VoteOptions))
	row2 = append(row2, renderProcessConfigs(details))
	row2 = append(row2, renderTimingDetails(details))
	row2 = append(row2, renderProcessCensusDetails(details))

	return elem.Div(
		elem.Div(
			vecty.Markup(vecty.Class("poll-details")),
			row1,
		),
		elem.Div(
			vecty.Markup(vecty.Class("poll-details")),
			row2,
		),
	)
}

func renderEnvelopeType(envelopeType api.EnvelopeType) vecty.ComponentOrHTML {
	return elem.Div(
		elem.Span(
			vecty.Markup(vecty.Class("detail")),
			vecty.Text("Envelope Type"),
		),
		elem.OrderedList(
			elem.ListItem(vecty.Text(fmt.Sprintf("Serial: %t", envelopeType.Serial))),
			elem.ListItem(vecty.Text(fmt.Sprintf("Anonymous: %t", envelopeType.Anonymous))),
			elem.ListItem(vecty.Text(fmt.Sprintf("Encrypted votes: %t", envelopeType.EncryptedVotes))),
			elem.ListItem(vecty.Text(fmt.Sprintf("Unique values: %t", envelopeType.UniqueValues))),
		))
}
func renderProcessMode(mode api.ProcessMode) vecty.ComponentOrHTML {
	return elem.Div(
		elem.Span(
			vecty.Markup(vecty.Class("detail")),
			vecty.Text("Process Mode"),
		),
		elem.OrderedList(
			elem.ListItem(vecty.Text(fmt.Sprintf("Auto start: %t", mode.AutoStart))),
			elem.ListItem(vecty.Text(fmt.Sprintf("Interruptible: %t", mode.Interruptible))),
			elem.ListItem(vecty.Text(fmt.Sprintf("Dynamic census: %t", mode.DynamicCensus))),
			elem.ListItem(vecty.Text(fmt.Sprintf("Encrypted metadata: %t", mode.EncryptedMetaData))),
		))
}
func renderProcessVoteOptions(options api.ProcessVoteOptions) vecty.ComponentOrHTML {
	return elem.Div(
		elem.Span(
			vecty.Markup(vecty.Class("detail")),
			vecty.Text("Vote Options"),
		),
		elem.OrderedList(
			elem.ListItem(vecty.Text(fmt.Sprintf("Max count: %d", options.MaxCount))),
			elem.ListItem(vecty.Text(fmt.Sprintf("Max value: %d", options.MaxValue))),
			elem.ListItem(vecty.Text(fmt.Sprintf("Max vote overwrites: %d", options.MaxVoteOverwrites))),
			elem.ListItem(vecty.Text(fmt.Sprintf("Max total cost: %d", options.MaxTotalCost))),
			elem.ListItem(vecty.Text(fmt.Sprintf("Cost exponent: %d", options.CostExponent))),
		))
}
func renderTimingDetails(details api.ProcessResults) vecty.ComponentOrHTML {
	return elem.Div(
		elem.Span(
			vecty.Markup(vecty.Class("detail")),
			vecty.Text("Timing & Results"),
		),
		elem.OrderedList(
			elem.ListItem(vecty.Text(fmt.Sprintf("Creation time: %s", details.CreationTime.Format("Mon Jan _2 15:04:05 +3:00 2006")))),
			elem.ListItem(vecty.Text(fmt.Sprintf("Start block: %d", details.StartBlock))),
			elem.ListItem(vecty.Text(fmt.Sprintf("End block: %d", details.EndBlock))),
		))
}

func renderProcessCensusDetails(details api.ProcessResults) vecty.ComponentOrHTML {
	return elem.Div(
		elem.Span(
			vecty.Markup(vecty.Class("detail")),
			vecty.Text("Census"),
		),
		elem.OrderedList(
			elem.ListItem(vecty.Text(fmt.Sprintf("Census origin: %s", details.CensusOrigin))),
			elem.ListItem(vecty.Text(fmt.Sprintf("Census root: %X", details.CensusRoot))),
			vecty.If(details.CensusURI != "",
				elem.ListItem(vecty.Text(fmt.Sprintf("Census URI: %s", details.CensusURI)))),
		))
}

func renderProcessConfigs(details api.ProcessResults) vecty.ComponentOrHTML {
	return elem.Div(
		elem.Span(
			vecty.Markup(vecty.Class("detail")),
			vecty.Text("Other Details"),
		),
		elem.OrderedList(
			elem.ListItem(vecty.Text(fmt.Sprintf("Namespace: %d", details.Namespace))),
			vecty.If(len(details.PrivateKeys) > 0,
				elem.ListItem(vecty.Text(fmt.Sprintf("Private Keys: %s", details.PrivateKeys)))),
			vecty.If(len(details.PublicKeys) > 0,
				elem.ListItem(vecty.Text(fmt.Sprintf("Public Keys: %s", details.PublicKeys)))),
			elem.ListItem(vecty.Text(fmt.Sprintf("Question Index: %d", details.QuestionIndex))),
		))
}

// UpdateProcessContents keeps the data for the processes dashboard up-to-date
func UpdateProcessContents(d *ProcessContentsView) {
	dispatcher.Dispatch(&actions.EnableAllUpdates{})
	process, ok := api.GetProcess(store.Processes.CurrentProcess.ID)
	if ok && process != nil {
		d.Unavailable = false
		dispatcher.Dispatch(&actions.SetCurrentProcessStruct{Process: process})
	} else {
		d.Unavailable = true
		dispatcher.Dispatch(&actions.SetCurrentProcessStruct{Process: nil})
		return
	}
	ticker := time.NewTicker(time.Duration(store.Config.RefreshTime) * time.Second)
	if !update.CheckCurrentPage("process", ticker) {
		return
	}
	updateProcessContents(d)
	for {
		select {
		case <-store.RedirectChan:
			if !update.CheckCurrentPage("process", ticker) {
				return
			}
		case <-ticker.C:
			if !update.CheckCurrentPage("process", ticker) {
				return
			}
			updateProcessContents(d)
		case i := <-store.Processes.EnvelopePagination.PagChannel:
			if !update.CheckCurrentPage("process", ticker) {
				return
			}
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
			if i < 1 {
				newVal, ok := api.GetProcessEnvelopeCount(store.Processes.CurrentProcess.ID)
				if ok {
					dispatcher.Dispatch(&actions.SetCurrentProcessEnvelopeHeight{Height: int(newVal)})
				}
			}
			if store.Processes.CurrentProcessResults.EnvelopeCount > 0 {
				updateProcessEnvelopes(d, util.Max(store.Processes.CurrentProcessResults.EnvelopeCount-store.Processes.EnvelopePagination.Index, 1))
			}
		}
	}
}

func updateProcessContents(d *ProcessContentsView) {
	dispatcher.Dispatch(&actions.ServerConnected{Connected: api.PingServer()})
	update.CurrentProcessResults()
	if !store.Envelopes.Pagination.DisableUpdate && store.Processes.CurrentProcessResults.EnvelopeCount > 0 {
		updateProcessEnvelopes(d, util.Max(store.Processes.CurrentProcessResults.EnvelopeCount-store.Processes.EnvelopePagination.Index, 1))
	}
}

func updateProcessEnvelopes(d *ProcessContentsView, index int) {
	logger.Info(fmt.Sprintf("Getting envelopes from index %d\n", index))
	list, ok := api.GetEnvelopeListByProcess(index, store.Processes.CurrentProcess.ID)
	if ok {
		reverseEnvelopeList(&list)
		dispatcher.Dispatch(&actions.SetCurrentProcessEnvelopes{EnvelopeList: list})
	}
}

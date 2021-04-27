package components

import (
	"fmt"
	"strings"
	"time"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/frontend/store/storeutil"
	"gitlab.com/vocdoni/vocexplorer/frontend/update"
	"gitlab.com/vocdoni/vocexplorer/logger"
	"gitlab.com/vocdoni/vocexplorer/util"
	"go.vocdoni.io/proto/build/go/models"
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
	if store.Processes.CurrentProcess == nil {
		return Unavailable("Loading process...", "")
	}

	if store.Processes.ProcessResults[util.HexToString(store.Processes.CurrentProcess.Process.ProcessId)].Type == "" {
		dispatcher.Dispatch(&actions.SetProcessType{Type: "Unknown"})
	}
	if store.Processes.ProcessResults[util.HexToString(store.Processes.CurrentProcess.Process.ProcessId)].State == "" {
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
		elem.Heading2(vecty.Text(util.HexToString(store.Processes.CurrentProcess.Process.ProcessId))),
		elem.Div(
			elem.Span(
				vecty.Markup(vecty.Class("title")),
				elem.Anchor(
					vecty.Markup(
						vecty.Attribute("href", store.ProcessDomain+util.HexToString(store.Processes.CurrentProcess.Process.EntityId)+"/0x"+util.HexToString(store.Processes.CurrentProcess.Process.ProcessId)),
						vecty.Property("target", util.HexToString(store.Processes.CurrentProcess.Process.ProcessId)),
					),
					vecty.Markup(vecty.Attribute("aria-label", "Link to process "+util.HexToString(store.Processes.CurrentProcess.Process.ProcessId)+"'s profile page")),
					vecty.Text("Process Profile"),
				),
			),
		),
		elem.Div(
			vecty.Markup(vecty.Class("badges")),
			elem.Span(
				vecty.Markup(vecty.Class("badge", store.Processes.ProcessResults[util.HexToString(store.Processes.CurrentProcess.Process.ProcessId)].State)),
				vecty.Text(strings.Title(store.Processes.ProcessResults[util.HexToString(store.Processes.CurrentProcess.Process.ProcessId)].State)),
			),
		),
		elem.HorizontalRule(),
		elem.DescriptionList(
			elem.DefinitionTerm(vecty.Text("Host entity")),
			elem.Description(
				Link(
					"/entity/"+util.HexToString(store.Processes.CurrentProcess.Process.EntityId),
					util.HexToString(store.Processes.CurrentProcess.Process.EntityId),
					"",
				),
			),
			elem.DefinitionTerm(vecty.Text("Process type")),
			elem.Description(vecty.Text(util.GetProcessName(store.Processes.ProcessResults[util.HexToString(store.Processes.CurrentProcess.Process.ProcessId)].Type))),
			elem.DefinitionTerm(vecty.Text("State")),
			elem.Description(vecty.Text(strings.Title(store.Processes.ProcessResults[util.HexToString(store.Processes.CurrentProcess.Process.ProcessId)].State))),
			elem.DefinitionTerm(vecty.Text("Registered votes")),
			elem.Description(vecty.Text(util.IntToString(store.Processes.Processes[util.HexToString(store.Processes.CurrentProcess.Process.ProcessId)].EnvelopeCount))),
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
			TabContents(results, renderResults(store.Processes.ProcessResults[util.HexToString(store.Processes.CurrentProcess.Process.ProcessId)].Results)),
			TabContents(envelopes, &ProcessesEnvelopeListView{}),
			TabContents(processDetails, renderProcessDetails(store.Processes.Processes[util.HexToString(store.Processes.CurrentProcess.Process.ProcessId)].Process)),
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

func renderProcessDetails(process *models.Process) vecty.ComponentOrHTML {
	row1 := vecty.List{}
	row2 := vecty.List{}

	// Add EnvelopeType
	row1 = append(row1, renderEnvelopeType(process.EnvelopeType))
	row1 = append(row1, renderProcessMode(process.Mode))
	row1 = append(row1, renderProcessVoteOptions(process.VoteOptions))
	row2 = append(row2, renderProcessConfigs(process))
	row2 = append(row2, renderTimingDetails(process))
	row2 = append(row2, renderProcessCensusDetails(process))

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

func renderEnvelopeType(envelopeType *models.EnvelopeType) vecty.ComponentOrHTML {
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
func renderProcessMode(mode *models.ProcessMode) vecty.ComponentOrHTML {
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
func renderProcessVoteOptions(options *models.ProcessVoteOptions) vecty.ComponentOrHTML {
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
func renderTimingDetails(process *models.Process) vecty.ComponentOrHTML {
	return elem.Div(
		elem.Span(
			vecty.Markup(vecty.Class("detail")),
			vecty.Text("Timing & Results"),
		),
		elem.OrderedList(
			// elem.ListItem(vecty.Text(fmt.Sprintf("Creation time: %s", process.Format("Mon Jan _2 15:04:05 +3:00 2006")))),
			elem.ListItem(vecty.Text(fmt.Sprintf("Start block: %d", process.StartBlock))),
			elem.ListItem(vecty.Text(fmt.Sprintf("End block: %d", process.StartBlock+process.BlockCount))),
		))
}

func renderProcessCensusDetails(process *models.Process) vecty.ComponentOrHTML {
	return elem.Div(
		elem.Span(
			vecty.Markup(vecty.Class("detail")),
			vecty.Text("Census"),
		),
		elem.OrderedList(
			elem.ListItem(vecty.Text(fmt.Sprintf("Census origin: %s", process.CensusOrigin))),
			elem.ListItem(vecty.Text(fmt.Sprintf("Census root: %X", process.CensusRoot))),
			vecty.If(*process.CensusURI != "",
				elem.ListItem(vecty.Text(fmt.Sprintf("Census URI: %s", *process.CensusURI)))),
		))
}

func renderProcessConfigs(process *models.Process) vecty.ComponentOrHTML {
	return elem.Div(
		elem.Span(
			vecty.Markup(vecty.Class("detail")),
			vecty.Text("Other Details"),
		),
		elem.OrderedList(
			elem.ListItem(vecty.Text(fmt.Sprintf("Namespace: %d", process.Namespace))),
			vecty.If(len(process.EncryptionPrivateKeys) > 0,
				elem.ListItem(vecty.Text(fmt.Sprintf("Private Keys: %s", process.EncryptionPrivateKeys)))),
			vecty.If(len(process.EncryptionPublicKeys) > 0,
				elem.ListItem(vecty.Text(fmt.Sprintf("Public Keys: %s", process.EncryptionPublicKeys)))),
			elem.ListItem(vecty.Text(fmt.Sprintf("Question Index: %d", process.QuestionIndex))),
		))
}

// UpdateProcessContents keeps the data for the processes dashboard up-to-date
func UpdateProcessContents(d *ProcessContentsView) {
	dispatcher.Dispatch(&actions.EnableAllUpdates{})
	process, rheight, creationTime, final, err := store.Client.GetProcess(store.Processes.CurrentProcess.Process.ProcessId)
	if err != nil {
		logger.Error(err)
		d.Unavailable = true
		dispatcher.Dispatch(&actions.SetCurrentProcessStruct{Process: nil})
		return
	}
	envelopeHeight, err := store.Client.GetEnvelopeHeight(store.Processes.CurrentProcess.Process.ProcessId)
	if err != nil {
		logger.Error(err)
	}
	d.Unavailable = false
	dispatcher.Dispatch(&actions.SetCurrentProcessStruct{
		Process: &storeutil.Process{
			EnvelopeCount: int(envelopeHeight),
			Process:       process,
			RHeight:       rheight,
			CreationTime:  creationTime,
			FinalResults:  final,
		},
	})
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
				newVal, err := store.Client.GetEnvelopeHeight(store.Processes.CurrentProcess.Process.ProcessId)
				if err == nil {
					dispatcher.Dispatch(&actions.SetCurrentProcessEnvelopeHeight{Height: int(newVal)})
				} else {
					logger.Error(err)
				}
			}
			if store.Processes.CurrentProcess.EnvelopeCount > 0 {
				updateProcessEnvelopes(d, util.Max(store.Processes.CurrentProcess.EnvelopeCount-store.Processes.EnvelopePagination.Index, 1))
			}
		}
	}
}

func updateProcessContents(d *ProcessContentsView) {
	dispatcher.Dispatch(&actions.GatewayConnected{GatewayErr: store.Client.GetGatewayInfo()})
	update.CurrentProcessResults()
	if !store.Envelopes.Pagination.DisableUpdate && store.Processes.CurrentProcess.EnvelopeCount > 0 {
		updateProcessEnvelopes(d, util.Max(store.Processes.CurrentProcess.EnvelopeCount-store.Processes.EnvelopePagination.Index, 1))
	}
}

func updateProcessEnvelopes(d *ProcessContentsView, index int) {
	logger.Info(fmt.Sprintf("Getting envelopes from index %d\n", index))
	list, err := store.Client.GetEnvelopeList(store.Processes.CurrentProcess.Process.ProcessId, config.ListSize)
	if err == nil {
		// TODO reverseEnvelopeList(list)
		dispatcher.Dispatch(&actions.SetCurrentProcessEnvelopes{EnvelopeList: list})
	} else {
		logger.Error(err)
	}
}

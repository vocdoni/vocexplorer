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
	sctypes "go.vocdoni.io/dvote/vochain/scrutinizer/types"
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

	if store.Processes.ProcessResults[util.HexToString(store.Processes.CurrentProcess.Process.ID)].Type == "" {
		dispatcher.Dispatch(&actions.SetProcessType{Type: "Unknown"})
	}
	if store.Processes.ProcessResults[util.HexToString(store.Processes.CurrentProcess.Process.ID)].State == "" {
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
	results := store.Processes.ProcessResults[util.HexToString(store.Processes.CurrentProcess.Process.ID)]
	if results.State == "" {
		results.State = "unknown"
	}
	if results.Type == "" {
		results.Type = "unknown"
	}
	return vecty.List{
		elem.Heading1(
			vecty.Text("Process details"),
		),
		elem.Heading2(vecty.Text(util.HexToString(store.Processes.CurrentProcess.Process.ID))),
		elem.Div(
			elem.Span(
				vecty.Markup(vecty.Class("title")),
				elem.Anchor(
					vecty.Markup(
						vecty.Attribute("href", store.ProcessDomain+util.HexToString(store.Processes.CurrentProcess.Process.EntityID)+"/0x"+util.HexToString(store.Processes.CurrentProcess.Process.ID)),
						vecty.Property("target", util.HexToString(store.Processes.CurrentProcess.Process.ID)),
					),
					vecty.Markup(vecty.Attribute("aria-label", "Link to process "+util.HexToString(store.Processes.CurrentProcess.Process.ID)+"'s profile page")),
					vecty.Text("Process Profile"),
				),
			),
		),
		elem.Div(
			vecty.Markup(vecty.Class("badges")),
			elem.Span(
				vecty.Markup(vecty.Class("badge", results.State)),
				vecty.Text(strings.Title(results.State)),
			),
		),
		elem.HorizontalRule(),
		elem.DescriptionList(
			elem.DefinitionTerm(vecty.Text("Host entity")),
			elem.Description(
				Link(
					"/entity/"+util.HexToString(store.Processes.CurrentProcess.Process.EntityID),
					util.HexToString(store.Processes.CurrentProcess.Process.EntityID),
					"",
				),
			),
			elem.DefinitionTerm(vecty.Text("Process type")),
			elem.Description(vecty.Text(util.GetProcessName(results.Type))),
			elem.DefinitionTerm(vecty.Text("State")),
			elem.Description(vecty.Text(strings.Title(results.State))),
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
			TabContents(results, renderResults(store.Processes.ProcessResults[util.HexToString(store.Processes.CurrentProcess.Process.ID)].Results)),
			TabContents(envelopes, &ProcessesEnvelopeListView{}),
			TabContents(processDetails, renderProcessDetails(store.Processes.CurrentProcess.Process)),
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

func renderProcessDetails(process *sctypes.Process) vecty.ComponentOrHTML {
	row1 := vecty.List{}
	row2 := vecty.List{}

	// Add EnvelopeType
	row1 = append(row1, renderEnvelopeType(process.Envelope))
	row1 = append(row1, renderProcessMode(process.Mode))
	row1 = append(row1, renderProcessVoteOptions(process.VoteOpts))
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
	if envelopeType == nil {
		return vecty.Text("Envelope Type unavailable")
	}
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
	if mode == nil {
		return vecty.Text("Process Mode unavailable")
	}
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
	if options == nil {
		return vecty.Text("Process Options unavailable")
	}
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
func renderTimingDetails(process *sctypes.Process) vecty.ComponentOrHTML {
	return elem.Div(
		elem.Span(
			vecty.Markup(vecty.Class("detail")),
			vecty.Text("Timing & Results"),
		),
		elem.OrderedList(
			// elem.ListItem(vecty.Text(fmt.Sprintf("Creation time: %s", process.Format("Mon Jan _2 15:04:05 +3:00 2006")))),
			elem.ListItem(vecty.Text(fmt.Sprintf("Start block: %d", process.StartBlock))),
			elem.ListItem(vecty.Text(fmt.Sprintf("End block: %d", process.EndBlock))),
		))
}

func renderProcessCensusDetails(process *sctypes.Process) vecty.ComponentOrHTML {
	return elem.Div(
		elem.Span(
			vecty.Markup(vecty.Class("detail")),
			vecty.Text("Census"),
		),
		elem.OrderedList(
			elem.ListItem(vecty.Text(fmt.Sprintf("Census origin: %d", process.CensusOrigin))),
			elem.ListItem(vecty.Text(fmt.Sprintf("Census root: %X", process.CensusRoot))),
			vecty.If(process.CensusURI != "",
				elem.ListItem(vecty.Text(fmt.Sprintf("Census URI: %s", process.CensusURI)))),
		))
}

func renderProcessConfigs(process *sctypes.Process) vecty.ComponentOrHTML {
	return elem.Div(
		elem.Span(
			vecty.Markup(vecty.Class("detail")),
			vecty.Text("Other Details"),
		),
		elem.OrderedList(
			elem.ListItem(vecty.Text(fmt.Sprintf("Namespace: %d", process.Namespace))),
			vecty.If(len(process.PrivateKeys) > 0,
				elem.ListItem(vecty.Text(fmt.Sprintf("Private Keys: %v", process.PrivateKeys)))),
			vecty.If(len(process.PublicKeys) > 0,
				elem.ListItem(vecty.Text(fmt.Sprintf("Public Keys: %v", process.PublicKeys)))),
			elem.ListItem(vecty.Text(fmt.Sprintf("Question Index: %d", process.QuestionIndex))),
		))
}

// UpdateProcessContents keeps the data for the processes dashboard up-to-date
func UpdateProcessContents(d *ProcessContentsView, pid []byte) {
	dispatcher.Dispatch(&actions.EnableAllUpdates{})
	process, err := store.Client.GetProcess(pid)
	if err != nil {
		logger.Error(err)
		d.Unavailable = true
		dispatcher.Dispatch(&actions.SetCurrentProcessStruct{Process: nil})
		return
	}
	pubKeys, privKeys, _, _, err := store.Client.GetProcessKeys(pid)
	if err != nil {
		logger.Error(err)
	}
	for _, key := range pubKeys {
		process.PublicKeys = append(process.PublicKeys, key.Key)
	}
	for _, key := range privKeys {
		process.PrivateKeys = append(process.PrivateKeys, key.Key)
	}
	envelopeHeight, err := store.Client.GetEnvelopeHeight(pid)
	if err != nil {
		logger.Error(err)
	}
	d.Unavailable = false
	dispatcher.Dispatch(&actions.SetCurrentProcessStruct{
		Process: &storeutil.Process{
			EnvelopeCount: int(envelopeHeight),
			Process:       process,
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
				newVal, err := store.Client.GetEnvelopeHeight(pid)
				if err == nil {
					dispatcher.Dispatch(&actions.SetCurrentProcessEnvelopeCount{Count: int(newVal)})
				} else {
					logger.Error(err)
				}
			}
			if store.Processes.CurrentProcess.EnvelopeCount > 0 {
				updateProcessEnvelopes(d, store.Processes.CurrentProcess.EnvelopeCount-store.Processes.EnvelopePagination.Index-config.ListSize)
			}
		}
	}
}

func updateProcessContents(d *ProcessContentsView) {
	dispatcher.Dispatch(&actions.GatewayConnected{GatewayErr: store.Client.GetGatewayInfo()})
	update.CurrentProcessResults()
	if !store.Envelopes.Pagination.DisableUpdate && store.Processes.CurrentProcess.EnvelopeCount > 0 {
		updateProcessEnvelopes(d, store.Processes.CurrentProcess.EnvelopeCount-store.Processes.EnvelopePagination.Index-config.ListSize)
	}
}

func updateProcessEnvelopes(d *ProcessContentsView, index int) {
	listSize := config.ListSize
	if index < 0 {
		listSize += index
		index = 0
	}
	logger.Info(fmt.Sprintf("Getting %d envelopes from index %d\n", listSize, index))
	list, err := store.Client.GetEnvelopeList(store.Processes.CurrentProcess.Process.ID, index, config.ListSize)
	if err == nil {
		dispatcher.Dispatch(&actions.SetCurrentProcessEnvelopes{EnvelopeList: list})
	} else {
		logger.Error(err)
	}
}

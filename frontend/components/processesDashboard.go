package components

import (
	"context"
	"fmt"
	"time"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/go-dvote/log"
	"gitlab.com/vocdoni/vocexplorer/client"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/dbapi"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// ProcessesDashboardView renders the processes dashboard page
type ProcessesDashboardView struct {
	vecty.Core
	gatewayConnected       bool
	serverConnected        bool
	gwClient               *client.Client
	process                *client.FullProcessInfo
	processID              string
	envelopeIndex          int
	disableEnvelopesUpdate bool
	quitCh                 chan struct{}
	refreshCh              chan int
}

// Render renders the ProcessesDashboardView component
func (dash *ProcessesDashboardView) Render() vecty.ComponentOrHTML {
	if dash == nil || dash.gwClient == nil || dash.process == nil {
		return &bootstrap.Alert{
			Contents: "Connecting to blockchain clients",
			Type:     "warning",
		}
	}

	t := dash.process.ProcessType
	if t == "" {
		t = "unknown"
	}
	st := dash.process.State
	if st == "" {
		st = "unknown"
	}

	return Container(
		renderGatewayConnectionBanner(dash.gatewayConnected),
		renderServerConnectionBanner(dash.serverConnected),
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

func (p *ProcessesDashboardView) ProcessDetails() vecty.List {
	t := p.process.ProcessType
	if t == "" {
		t = "unknown"
	}
	st := p.process.State
	if st == "" {
		st = "unknown"
	}

	return vecty.List{
		elem.Heading1(
			vecty.Text("Process details"),
		),
		elem.Heading2(vecty.Text(p.processID)),
		elem.Div(
			vecty.Markup(vecty.Class("badges")),
			elem.Span(
				vecty.Markup(vecty.Class("badge", st)),
				vecty.Text(st),
			),
		),
		elem.HorizontalRule(),
		elem.DescriptionList(
			elem.DefinitionTerm(vecty.Text("Process type")),
			elem.Description(vecty.Text(t)),
			elem.DefinitionTerm(vecty.Text("State")),
			elem.Description(vecty.Text(st)),
			elem.DefinitionTerm(vecty.Text("Registered votes")),
			elem.Description(vecty.Text(util.IntToString(p.process.EnvelopeHeight))),
		),
	}
}

type ProcessTab struct {
	*Tab
}

func (p *ProcessTab) dispatch() interface{} {
	return &actions.ProcessesTabChange{
		Tab: p.alias(),
	}
}

func (p *ProcessTab) store() string {
	return store.Processes.Tab
}

func (p *ProcessesDashboardView) ProcessTabs() vecty.List {
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
				TabLink(p, results),
				TabLink(p, envelopes),
			),
		),
		elem.Div(
			vecty.Markup(vecty.Class("tabs-content")),
			TabContents(results, renderResults(p.process.Results)),
			TabContents(envelopes, &ProcessesEnvelopeListView{
				process:       p.process,
				refreshCh:     p.refreshCh,
				disableUpdate: &p.disableEnvelopesUpdate,
			}),
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

// InitProcessesDashboardView initializes the processes dashboard view
func InitProcessesDashboardView(process *client.FullProcessInfo, ProcessesDashboardView *ProcessesDashboardView, processID string, cfg *config.Cfg) *ProcessesDashboardView {
	gwClient, cancel := client.InitGateway(cfg.GatewayHost)
	if gwClient == nil {
		return ProcessesDashboardView
	}
	ProcessesDashboardView.gwClient = gwClient
	ProcessesDashboardView.process = process
	ProcessesDashboardView.processID = processID
	ProcessesDashboardView.quitCh = make(chan struct{})
	ProcessesDashboardView.refreshCh = make(chan int, 50)
	BeforeUnload(func() {
		close(ProcessesDashboardView.quitCh)
	})
	go updateAndRenderProcessesDashboard(ProcessesDashboardView, cancel, processID, cfg)
	return ProcessesDashboardView
}

func updateAndRenderProcessesDashboard(d *ProcessesDashboardView, cancel context.CancelFunc, processID string, cfg *config.Cfg) {
	ticker := time.NewTicker(time.Duration(cfg.RefreshTime) * time.Second)
	updateProcessesDashboard(d, processID)
	vecty.Rerender(d)
	for {
		select {
		case <-d.quitCh:
			ticker.Stop()
			d.gwClient.Close()
			fmt.Println("Gateway connection closed")
			return
		case <-ticker.C:
			updateProcessesDashboard(d, processID)
			vecty.Rerender(d)
		case i := <-d.refreshCh:
		loop:
			for {
				// If many indices waiting in buffer, scan to last one.
				select {
				case i = <-d.refreshCh:
				default:
					break loop
				}
			}
			d.envelopeIndex = i
			oldEnvelopes := d.process.EnvelopeHeight
			newVal, ok := dbapi.GetProcessEnvelopeHeight(processID)
			if ok {
				d.process.EnvelopeHeight = int(newVal)
			}
			if i < 1 {
				oldEnvelopes = d.process.EnvelopeHeight
			}
			if d.process.EnvelopeHeight > 0 {
				updateProcessEnvelopes(d, util.Max(oldEnvelopes-d.envelopeIndex, config.ListSize))
			}
			vecty.Rerender(d)
		}
	}
}

func updateProcessesDashboard(d *ProcessesDashboardView, processID string) {
	if d.gwClient.Conn.Ping(d.gwClient.Ctx) != nil {
		d.gatewayConnected = false
	} else {
		d.gatewayConnected = true
	}
	if !dbapi.Ping() {
		d.serverConnected = false
	} else {
		d.serverConnected = true
	}
	client.UpdateProcessesDashboardInfo(d.gwClient, d.process, processID)
	newVal, ok := dbapi.GetProcessEnvelopeHeight(processID)
	if ok {
		d.process.EnvelopeHeight = int(newVal)
	}
	if !d.disableEnvelopesUpdate && d.process.EnvelopeHeight > 0 {
		updateProcessEnvelopes(d, util.Max(d.process.EnvelopeHeight-d.envelopeIndex, config.ListSize))
	}
}

func updateProcessEnvelopes(d *ProcessesDashboardView, index int) {
	log.Infof("Getting envelopes from index %d", util.IntToString(index))
	list, ok := dbapi.GetEnvelopeListByProcess(index, d.processID)
	if ok {
		reverseEnvelopeList(&list)
		d.process.EnvelopeList = list
	}
}

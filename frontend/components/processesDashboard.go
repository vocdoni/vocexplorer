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
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// ProcessesDashboardView renders the processes dashboard page
type ProcessesDashboardView struct {
	vecty.Core
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
	if dash != nil && dash.gwClient != nil && dash.process != nil {
		t := dash.process.ProcessType
		if t == "" {
			t = "unknown"
		}
		st := dash.process.State
		if st == "" {
			st = "unknown"
		}
		return Container(
			elem.Section(
				elem.Heading4(vecty.Text(
					fmt.Sprintf("Process %s", dash.processID),
				)),
				elem.Heading5(vecty.Text("Process type: "+t+", state: "+st)),
				elem.Heading5(vecty.Text("Number of votes : "+util.IntToString(dash.process.EnvelopeHeight))),

				renderResults(dash.process.Results),
				vecty.Markup(vecty.Class("info-pane")),
				&ProcessesEnvelopeListView{
					process:       dash.process,
					refreshCh:     dash.refreshCh,
					disableUpdate: &dash.disableEnvelopesUpdate,
				},
			),
		)
	}
	return &bootstrap.Alert{
		Contents: "Connecting to blockchain clients",
		Type:     "warning",
	}
}

func renderResults(results [][]uint32) vecty.ComponentOrHTML {
	if len(results) <= 0 {
		return elem.Heading6(vecty.Text("No results yet"))
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
	client.UpdateProcessesDashboardInfo(d.gwClient, d.process, processID)
	d.process.EnvelopeHeight = int(dbapi.GetProcessEnvelopeHeight(processID))
	if d.process.EnvelopeHeight > 0 {
		updateProcessEnvelopes(d, util.Max(d.process.EnvelopeHeight-d.envelopeIndex, config.ListSize))
	}
	vecty.Rerender(d)
	for {
		select {
		case <-d.quitCh:
			ticker.Stop()
			d.gwClient.Close()
			fmt.Println("Gateway connection closed")
			return
		case <-ticker.C:
			client.UpdateProcessesDashboardInfo(d.gwClient, d.process, processID)
			d.process.EnvelopeHeight = int(dbapi.GetProcessEnvelopeHeight(processID))
			if !d.disableEnvelopesUpdate && d.process.EnvelopeHeight > 0 {
				updateProcessEnvelopes(d, util.Max(d.process.EnvelopeHeight-d.envelopeIndex, config.ListSize))
			}
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
			d.process.EnvelopeHeight = int(dbapi.GetProcessEnvelopeHeight(processID))
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

func updateProcessEnvelopes(d *ProcessesDashboardView, index int) {
	log.Infof("Getting envelopes from index %d", util.IntToString(index))
	list := dbapi.GetEnvelopeListByProcess(index, d.processID)
	reverseEnvelopeList(&list)
	d.process.EnvelopeList = list
}

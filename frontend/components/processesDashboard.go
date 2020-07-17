package components

import (
	"context"
	"fmt"
	"time"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/client"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// ProcessesDashboardView renders the processes dashboard page
type ProcessesDashboardView struct {
	vecty.Core
	processID string
	process   *client.FullProcessInfo
	gwClient  *client.Client
	quitCh    chan struct{}
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
			t = "unknown"
		}
		return elem.Div(
			elem.Main(
				elem.Heading4(vecty.Text("Process "+dash.processID)),
				elem.Heading5(vecty.Text("Process type: "+t+", state: "+st)),
				renderResults(dash.process.Results),
				vecty.Markup(vecty.Class("info-pane")),
				&EnvelopeListView{
					process: dash.process,
				},
			),
		)
	}
	return vecty.Text("Connecting to blockchain clients")
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

func initProcessesDashboardView(process *client.FullProcessInfo, ProcessesDashboardView *ProcessesDashboardView, processID string) *ProcessesDashboardView {
	gwClient, cancel := InitGateway()
	if gwClient == nil {
		return ProcessesDashboardView
	}
	ProcessesDashboardView.gwClient = gwClient
	ProcessesDashboardView.process = process
	ProcessesDashboardView.processID = processID
	ProcessesDashboardView.quitCh = make(chan struct{})
	BeforeUnload(func() {
		close(ProcessesDashboardView.quitCh)
	})
	go updateAndRenderProcessesDashboard(ProcessesDashboardView, cancel, processID)
	return ProcessesDashboardView
}

func updateAndRenderProcessesDashboard(d *ProcessesDashboardView, cancel context.CancelFunc, processID string) {
	ticker := time.NewTicker(config.RefreshTime * time.Second)
	// Wait for data structs to load
	for d == nil || d.process == nil {
	}
	client.UpdateProcessesDashboardInfo(d.gwClient, d.process, processID)
	vecty.Rerender(d)
	// time.Sleep(250 * time.Millisecond)
	// client.UpdateAuxProcessInfo(d.gwClient, d.vc)
	// vecty.Rerender(d)
	for {
		select {
		case <-d.quitCh:
			ticker.Stop()
			d.gwClient.Close()
			fmt.Println("Gateway connection closed")
			return
		case <-ticker.C:
			client.UpdateProcessesDashboardInfo(d.gwClient, d.process, processID)
			// client.UpdateAuxProcessInfo(d.gwClient, d.process)
			vecty.Rerender(d)
		}
	}
}

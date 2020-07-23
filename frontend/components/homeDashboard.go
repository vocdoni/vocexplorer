package components

import (
	"context"
	"fmt"
	"syscall/js"
	"time"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/tendermint/tendermint/rpc/client/http"
	"gitlab.com/vocdoni/vocexplorer/client"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/rpc"
)

// DashboardView renders the dashboard landing page
type DashboardView struct {
	vecty.Core
	t          *rpc.TendermintInfo
	vc         *client.VochainInfo
	gwClient   *client.Client
	tClient    *http.HTTP
	quitCh     chan struct{}
	refreshCh  chan int
	blockIndex int
}

// Render renders the DashboardView component
func (dash *DashboardView) Render() vecty.ComponentOrHTML {
	if dash != nil && dash.gwClient != nil && dash.tClient != nil && dash.t != nil && dash.vc != nil {
		return elem.Div(
			elem.Main(
				vecty.Markup(vecty.Class("info-pane")),
				&StatsView{
					t:  dash.t,
					vc: dash.vc,
				},
			),
			vecty.Markup(
				event.BeforeUnload(func(i *vecty.Event) {
					js.Global().Get("alert").Invoke("Closing page")
					dash.gwClient.Close()
				},
				),
			),
		)
	}
	return vecty.Text("Connecting to blockchain clients")
}

func initDashboardView(t *rpc.TendermintInfo, vc *client.VochainInfo, DashboardView *DashboardView) *DashboardView {
	// Init tendermint client
	tClient := rpc.StartClient()
	// Init Gateway client
	gwClient, cancel := InitGateway()
	if gwClient == nil || tClient == nil {
		return DashboardView
	}
	DashboardView.tClient = tClient
	DashboardView.gwClient = gwClient
	DashboardView.t = t
	DashboardView.vc = vc
	DashboardView.quitCh = make(chan struct{})
	BeforeUnload(func() {
		close(DashboardView.quitCh)
	})
	go updateAndRenderDashboard(DashboardView, cancel)
	return DashboardView
}

func updateAndRenderDashboard(d *DashboardView, cancel context.CancelFunc) {
	ticker := time.NewTicker(config.RefreshTime * time.Second)
	// Wait for data structs to load
	for d == nil || d.vc == nil {
	}
	rpc.UpdateTendermintInfo(d.tClient, d.t)
	client.UpdateDashboardInfo(d.gwClient, d.vc)
	vecty.Rerender(d)
	for {
		select {
		case <-d.quitCh:
			ticker.Stop()
			d.gwClient.Close()
			//cancel()
			fmt.Println("Gateway connection closed")
			return
		case <-ticker.C:
			rpc.UpdateTendermintInfo(d.tClient, d.t, d.d.blockIndex)
			client.UpdateDashboardInfo(d.gwClient, d.vc)
			vecty.Rerender(d)
		case i <- d.refreshCh:
			d.blockIndex = i
			rpc.UpdateBlockList(c.tClient, d.t, d.blockIndex)
		}
	}
}
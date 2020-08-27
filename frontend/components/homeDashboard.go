package components

import (
	"context"
	"fmt"
	"time"

	"github.com/gopherjs/vecty"
	"github.com/tendermint/tendermint/rpc/client/http"
	"gitlab.com/vocdoni/vocexplorer/client"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/dbapi"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/rpc"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// DashboardView renders the dashboard landing page
type DashboardView struct {
	vecty.Core
	gatewayConnected bool
	serverConnected  bool
	blockIndex       int
	gwClient         *client.Client
	t                *rpc.TendermintInfo
	tClient          *http.HTTP
	vc               *client.VochainInfo
}

// Render renders the DashboardView component
func (dash *DashboardView) Render() vecty.ComponentOrHTML {
	if dash != nil && dash.gwClient != nil && dash.tClient != nil && dash.t != nil && dash.vc != nil {
		return Container(
			renderGatewayConnectionBanner(dash.gatewayConnected),
			renderServerConnectionBanner(dash.serverConnected),
			&StatsView{
				t:        dash.t,
				vc:       dash.vc,
				gwClient: dash.gwClient,
			},
		)
	}
	return &bootstrap.Alert{
		Contents: "Connecting to blockchain clients",
		Type:     "warning",
	}
}

// InitDashboardView returns the home dashboard components
func InitDashboardView(t *rpc.TendermintInfo, vc *client.VochainInfo, DashboardView *DashboardView, cfg *config.Cfg) *DashboardView {
	// Init tendermint client
	tClient := rpc.StartClient(cfg.TendermintHost)
	// Init Gateway client
	gwClient, cancel := client.InitGateway(cfg.GatewayHost)
	if gwClient == nil || tClient == nil {
		return DashboardView
	}
	DashboardView.tClient = tClient
	DashboardView.gwClient = gwClient
	DashboardView.t = t
	DashboardView.vc = vc
	DashboardView.blockIndex = 0
	DashboardView.serverConnected = true
	DashboardView.gatewayConnected = true
	BeforeUnload(func() {
		dispatcher.Dispatch(&actions.SignalRedirect{})
	})
	go updateAndRenderDashboard(DashboardView, cancel, cfg)
	return DashboardView
}

func updateHeight(t *rpc.TendermintInfo) {
	newVal, ok := dbapi.GetBlockHeight()
	if ok {
		t.TotalBlocks = int(newVal - 1)
		dispatcher.Dispatch(&actions.BlocksHeightUpdate{
			Height: int64(newVal),
		})
	}
	newVal, ok = dbapi.GetTxHeight()
	if ok {
		t.TotalTxs = int(newVal - 1)
	}
	newVal, ok = dbapi.GetEntityHeight()
	if ok {
		t.TotalEntities = int(newVal)
	}
	newVal, ok = dbapi.GetProcessHeight()
	if ok {
		t.TotalProcesses = int(newVal)
	}
	newVal, ok = dbapi.GetEnvelopeHeight()
	if ok {
		t.TotalEnvelopes = int(newVal)
	}
}

func updateAndRenderDashboard(d *DashboardView, cancel context.CancelFunc, cfg *config.Cfg) {
	ticker := time.NewTicker(time.Duration(cfg.RefreshTime) * time.Second)
	updateHomeDashboardInfo(d)
	vecty.Rerender(d)
	for {
		select {
		case <-store.RedirectChan:
			fmt.Println("Redirecting...")
			ticker.Stop()
			d.gwClient.Close()
			//cancel()
			return
		case <-ticker.C:
			updateHomeDashboardInfo(d)
			vecty.Rerender(d)
		}
	}
}

func updateHomeDashboardInfo(d *DashboardView) {
	if !rpc.Ping(d.tClient) || d.gwClient.Conn.Ping(d.gwClient.Ctx) != nil {
		d.gatewayConnected = false
	} else {
		d.gatewayConnected = true
	}
	if !dbapi.Ping() {
		d.serverConnected = false
	} else {
		d.serverConnected = true
	}
	rpc.UpdateTendermintInfo(d.tClient, d.t)
	client.UpdateDashboardInfo(d.gwClient, d.vc)
	updateHeight(d.t)
	updateHomeBlocks(d, util.Max(d.t.TotalBlocks-d.blockIndex, config.HomeWidgetBlocksListSize))
}

func updateHomeBlocks(d *DashboardView, index int) {
	fmt.Println("Getting blocks from index " + util.IntToString(index))
	list, ok := dbapi.GetBlockList(index)
	if ok {
		reverseBlockList(&list)
		d.t.BlockList = list
	}
}

package components

import (
	"context"
	"fmt"
	"time"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/tendermint/tendermint/rpc/client/http"
	"gitlab.com/vocdoni/vocexplorer/client"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/dbapi"
	"gitlab.com/vocdoni/vocexplorer/rpc"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// DashboardView renders the dashboard landing page
type DashboardView struct {
	vecty.Core
	blockIndex int
	gwClient   *client.Client
	quitCh     chan struct{}
	refreshCh  chan int
	t          *rpc.TendermintInfo
	tClient    *http.HTTP
	vc         *client.VochainInfo
}

// Render renders the DashboardView component
func (dash *DashboardView) Render() vecty.ComponentOrHTML {
	if dash != nil && dash.gwClient != nil && dash.tClient != nil && dash.t != nil && dash.vc != nil {
		return elem.Main(
			vecty.Markup(vecty.Class("home")),
			&StatsView{
				t:         dash.t,
				vc:        dash.vc,
				refreshCh: dash.refreshCh,
			},
			vecty.Markup(
				event.BeforeUnload(func(i *vecty.Event) {
					dash.gwClient.Close()
				},
				),
			),
		)
	}
	return vecty.Text("Connecting to blockchain clients")
}

func initDashboardView(t *rpc.TendermintInfo, vc *client.VochainInfo, DashboardView *DashboardView, cfg *config.Cfg) *DashboardView {
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
	DashboardView.quitCh = make(chan struct{})
	DashboardView.refreshCh = make(chan int, 50)
	DashboardView.blockIndex = 0
	BeforeUnload(func() {
		close(DashboardView.quitCh)
	})
	go updateAndRenderDashboard(DashboardView, cancel, cfg)
	return DashboardView
}

func updateAndRenderDashboard(d *DashboardView, cancel context.CancelFunc, cfg *config.Cfg) {
	ticker := time.NewTicker(time.Duration(cfg.RefreshTime) * time.Second)
	rpc.UpdateTendermintInfo(d.tClient, d.t)
	client.UpdateDashboardInfo(d.gwClient, d.vc)
	d.t.TotalBlocks = int(dbapi.GetBlockHeight()) - 1
	updateHomeBlocks(d, util.Max(d.t.TotalBlocks-d.blockIndex, config.HomeWidgetBlocksListSize))
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
			rpc.UpdateTendermintInfo(d.tClient, d.t)
			d.t.TotalBlocks = int(dbapi.GetBlockHeight()) - 1
			updateHomeBlocks(d, util.Max(d.t.TotalBlocks-d.blockIndex, config.HomeWidgetBlocksListSize))

			client.UpdateDashboardInfo(d.gwClient, d.vc)
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
			d.blockIndex = i
			oldBlocks := d.t.TotalBlocks
			d.t.TotalBlocks = int(dbapi.GetBlockHeight()) - 1
			if i < 1 {
				oldBlocks = d.t.TotalBlocks
			}
			updateHomeBlocks(d, util.Max(oldBlocks-d.blockIndex, config.HomeWidgetBlocksListSize))
			vecty.Rerender(d)
		}
	}
}

func updateHomeBlocks(d *DashboardView, index int) {
	fmt.Println("Getting blocks from index " + util.IntToString(index))
	list := dbapi.GetBlockList(index)
	for i := len(list)/2 - 1; i >= 0; i-- {
		opp := len(list) - 1 - i
		list[i], list[opp] = list[opp], list[i]
	}
	d.t.BlockList = list
}

package components

import (
	"context"
	"fmt"
	"time"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/client"
	"gitlab.com/vocdoni/vocexplorer/config"
)

// VocDashDashboardView renders the processes dashboard page
type VocDashDashboardView struct {
	vecty.Core
	gwClient  *client.Client
	quitCh    chan struct{}
	refreshCh chan bool
	vc        *client.VochainInfo
}

// Render renders the VocDashDashboardView component
func (dash *VocDashDashboardView) Render() vecty.ComponentOrHTML {
	if dash != nil && dash.gwClient != nil && dash.vc != nil {
		return elem.Div(
			elem.Main(
				vecty.Markup(vecty.Class("info-pane")),
				&VochainInfoView{
					vc:        dash.vc,
					refreshCh: dash.refreshCh,
				},
			),
		)
	}
	return vecty.Text("Connecting to blockchain clients")
}

func initVocDashDashboardView(vc *client.VochainInfo, VocDashDashboardView *VocDashDashboardView, cfg *config.Cfg) *VocDashDashboardView {
	gwClient, cancel := client.InitGateway(cfg.GatewayHost)
	if gwClient == nil {
		return VocDashDashboardView
	}
	VocDashDashboardView.gwClient = gwClient
	VocDashDashboardView.vc = vc
	VocDashDashboardView.quitCh = make(chan struct{})
	VocDashDashboardView.refreshCh = make(chan bool, 20)
	BeforeUnload(func() {
		close(VocDashDashboardView.quitCh)
	})
	go updateAndRenderVocDashDashboard(VocDashDashboardView, cancel, cfg)
	return VocDashDashboardView
}

func updateAndRenderVocDashDashboard(d *VocDashDashboardView, cancel context.CancelFunc, cfg *config.Cfg) {
	ticker := time.NewTicker(time.Duration(cfg.RefreshTime) * time.Second)
	//TODO: update to  use real index
	client.UpdateVocDashDashboardInfo(d.gwClient, d.vc, 0)
	vecty.Rerender(d)
	time.Sleep(250 * time.Millisecond)
	client.UpdateAuxProcessInfo(d.gwClient, d.vc)
	vecty.Rerender(d)
	for {
		select {
		case <-d.quitCh:
			ticker.Stop()
			d.gwClient.Close()
			fmt.Println("Gateway connection closed")
			return
		case <-ticker.C:
			//TODO: update to  use real index
			client.UpdateVocDashDashboardInfo(d.gwClient, d.vc, 0)
			client.UpdateAuxProcessInfo(d.gwClient, d.vc)
			vecty.Rerender(d)
		case <-d.refreshCh:
			//TODO: update to  use real index
			client.UpdateVocDashDashboardInfo(d.gwClient, d.vc, 0)
			client.UpdateAuxProcessInfo(d.gwClient, d.vc)
			vecty.Rerender(d)
		}
	}
}

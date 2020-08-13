package components

import (
	"context"
	"fmt"
	"time"

	"github.com/gopherjs/vecty"
	"gitlab.com/vocdoni/vocexplorer/client"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
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
		return Container(
			vecty.Markup(vecty.Class("info-pane")),
			&VochainInfoView{
				vc:        dash.vc,
				refreshCh: dash.refreshCh,
			},
		)
	}
	return &bootstrap.Alert{
		Contents: "Connecting to blockchain clients",
		Type:     "warning",
	}
}

// InitVocDashDashboardView initializes the vocdash page
func InitVocDashDashboardView(vc *client.VochainInfo, VocDashDashboardView *VocDashDashboardView, cfg *config.Cfg) *VocDashDashboardView {
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
	// Wait for data structs to load
	for d == nil || d.vc == nil {
	}
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

package components

import (
	"context"
	"fmt"
	"syscall/js"
	"time"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/client"
	"gitlab.com/vocdoni/vocexplorer/config"
)

// VocDashDashboardView renders the processes dashboard page
type VocDashDashboardView struct {
	vecty.Core
	vc       *client.VochainInfo
	gwClient *client.Client
	quitCh   chan struct{}
}

// Render renders the VocDashDashboardView component
func (dash *VocDashDashboardView) Render() vecty.ComponentOrHTML {
	if dash != nil && dash.gwClient != nil && dash.vc != nil {
		return elem.Div(
			elem.Main(
				vecty.Markup(vecty.Class("info-pane")),
				&VochainInfoView{
					vc: dash.vc,
				},
			),
		)
	}
	return vecty.Text("Connecting to blockchain clients")
}

func initVocDashDashboardView(vc *client.VochainInfo, VocDashDashboardView *VocDashDashboardView) *VocDashDashboardView {
	js.Global().Set("apiEnabled", true)
	gwClient, cancel := InitGateway()
	if gwClient == nil {
		return VocDashDashboardView
	}
	VocDashDashboardView.gwClient = gwClient
	VocDashDashboardView.vc = vc
	VocDashDashboardView.quitCh = make(chan struct{})
	BeforeUnload(func() {
		close(VocDashDashboardView.quitCh)
	})
	go updateAndRenderVocDashDashboard(VocDashDashboardView, cancel)
	return VocDashDashboardView
}

func updateAndRenderVocDashDashboard(d *VocDashDashboardView, cancel context.CancelFunc) {
	ticker := time.NewTicker(config.RefreshTime * time.Second)
	// Wait for data structs to load
	for d == nil || d.vc == nil {
	}
	client.UpdateVocDashDashboardInfo(d.gwClient, d.vc)
	vecty.Rerender(d)
	time.Sleep(250 * time.Millisecond)
	client.UpdateAuxProcessInfo(d.gwClient, d.vc)
	vecty.Rerender(d)
	for {
		select {
		case <-d.quitCh:
			ticker.Stop()
			d.gwClient.Close()
			// cancel()
			fmt.Println("Gateway connection closed")
			return
		case <-ticker.C:
			client.UpdateVocDashDashboardInfo(d.gwClient, d.vc)

			// if ticks%10 == 0 {
			client.UpdateAuxProcessInfo(d.gwClient, d.vc)
			// }
			vecty.Rerender(d)
		}
	}
}

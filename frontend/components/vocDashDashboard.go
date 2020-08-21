package components

import (
	"context"
	"fmt"
	"time"

	"github.com/gopherjs/vecty"
	"gitlab.com/vocdoni/go-dvote/log"
	"gitlab.com/vocdoni/vocexplorer/client"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/dbapi"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// VocDashDashboardView renders the processes dashboard page
type VocDashDashboardView struct {
	vecty.Core
	gwClient               *client.Client
	envelopeIndex          int
	quitCh                 chan struct{}
	refreshCh              chan bool
	disableEnvelopesUpdate bool
	refreshEnvelopes       chan int
	vc                     *client.VochainInfo
}

// Render renders the VocDashDashboardView component
func (dash *VocDashDashboardView) Render() vecty.ComponentOrHTML {
	if dash != nil && dash.gwClient != nil && dash.vc != nil {
		return Container(
			&VochainInfoView{
				vc:        dash.vc,
				refreshCh: dash.refreshCh,
			},
			&EnvelopeListView{
				vochain:       dash.vc,
				refreshCh:     dash.refreshEnvelopes,
				disableUpdate: &dash.disableEnvelopesUpdate,
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
	VocDashDashboardView.refreshEnvelopes = make(chan int, 50)
	VocDashDashboardView.disableEnvelopesUpdate = false

	BeforeUnload(func() {
		close(VocDashDashboardView.quitCh)
	})
	go updateAndRenderVocDashDashboard(VocDashDashboardView, cancel, cfg)
	return VocDashDashboardView
}

func updateAndRenderVocDashDashboard(d *VocDashDashboardView, cancel context.CancelFunc, cfg *config.Cfg) {
	ticker := time.NewTicker(time.Duration(cfg.RefreshTime) * time.Second)
	d.vc.EnvelopeHeight = int(dbapi.GetEnvelopeHeight())
	updateEnvelopes(d, util.Max(d.vc.EnvelopeHeight-d.envelopeIndex, config.ListSize))
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
			d.vc.EnvelopeHeight = int(dbapi.GetEnvelopeHeight())
			if !d.disableEnvelopesUpdate {
				updateEnvelopes(d, util.Max(d.vc.EnvelopeHeight-d.envelopeIndex, config.ListSize))
			}
			client.UpdateVocDashDashboardInfo(d.gwClient, d.vc, 0)
			client.UpdateAuxProcessInfo(d.gwClient, d.vc)
			vecty.Rerender(d)
		case <-d.refreshCh:
			d.vc.EnvelopeHeight = int(dbapi.GetEnvelopeHeight())
			updateEnvelopes(d, util.Max(d.vc.EnvelopeHeight-d.envelopeIndex, config.ListSize))
			client.UpdateVocDashDashboardInfo(d.gwClient, d.vc, 0)
			client.UpdateAuxProcessInfo(d.gwClient, d.vc)
			vecty.Rerender(d)
		case i := <-d.refreshEnvelopes:
		loop:
			for {
				// If many indices waiting in buffer, scan to last one.
				select {
				case i = <-d.refreshEnvelopes:
				default:
					break loop
				}
			}
			d.envelopeIndex = i
			oldEnvelopes := d.vc.EnvelopeHeight
			d.vc.EnvelopeHeight = int(dbapi.GetEnvelopeHeight())
			if i < 1 {
				oldEnvelopes = d.vc.EnvelopeHeight
			}
			if d.vc.EnvelopeHeight > 0 {
				updateEnvelopes(d, util.Max(oldEnvelopes-d.envelopeIndex, config.ListSize))
			}
			vecty.Rerender(d)
		}
	}
}

func updateEnvelopes(d *VocDashDashboardView, index int) {
	log.Infof("Getting envelopes from index %d", util.IntToString(index))
	list := dbapi.GetEnvelopeList(index)
	reverseEnvelopeList(&list)
	d.vc.EnvelopeList = list
}

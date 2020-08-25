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
	"gitlab.com/vocdoni/vocexplorer/types"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// VocDashDashboardView renders the processes dashboard page
type VocDashDashboardView struct {
	vecty.Core
	gwClient               *client.Client
	envelopeIndex          int
	entityIndex            int
	processIndex           int
	quitCh                 chan struct{}
	disableEnvelopesUpdate bool
	disableEntitiesUpdate  bool
	disableProcessesUpdate bool
	refreshEnvelopes       chan int
	refreshEntities        chan int
	refreshProcesses       chan int
	vc                     *client.VochainInfo
}

// Render renders the VocDashDashboardView component
func (dash *VocDashDashboardView) Render() vecty.ComponentOrHTML {
	if dash != nil && dash.gwClient != nil && dash.vc != nil {
		return Container(
			&ProcessListView{
				vochain:       dash.vc,
				refreshCh:     dash.refreshProcesses,
				disableUpdate: &dash.disableProcessesUpdate,
			},
			&EntityListView{
				vochain:       dash.vc,
				refreshCh:     dash.refreshEntities,
				disableUpdate: &dash.disableEntitiesUpdate,
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
	VocDashDashboardView.refreshEnvelopes = make(chan int, 50)
	VocDashDashboardView.refreshProcesses = make(chan int, 50)
	VocDashDashboardView.refreshEntities = make(chan int, 50)
	VocDashDashboardView.disableEnvelopesUpdate = false

	BeforeUnload(func() {
		close(VocDashDashboardView.quitCh)
	})
	go updateAndRenderVocDashDashboard(VocDashDashboardView, cancel, cfg)
	return VocDashDashboardView
}

func updateAndRenderVocDashDashboard(d *VocDashDashboardView, cancel context.CancelFunc, cfg *config.Cfg) {
	ticker := time.NewTicker(time.Duration(cfg.RefreshTime) * time.Second)
	updateVocdash(d)
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
			updateVocdash(d)
			client.UpdateAuxProcessInfo(d.gwClient, d.vc)
			vecty.Rerender(d)
		case i := <-d.refreshEntities:
		entityLoop:
			for {
				// If many indices waiting in buffer, scan to last one.
				select {
				case i = <-d.refreshEntities:
				default:
					break entityLoop
				}
			}
			d.entityIndex = i
			oldEntities := d.vc.EntityCount
			d.vc.EntityCount = int(dbapi.GetEntityHeight())
			if i < 1 {
				oldEntities = d.vc.EntityCount
			}
			if d.vc.EntityCount > 0 {
				updateEntities(d, util.Max(oldEntities-d.entityIndex-1, config.ListSize-1))
			}
			vecty.Rerender(d)
		case i := <-d.refreshProcesses:
		processLoop:
			for {
				// If many indices waiting in buffer, scan to last one.
				select {
				case i = <-d.refreshProcesses:
				default:
					break processLoop

				}
			}
			d.processIndex = i
			oldProcesses := d.vc.ProcessCount
			d.vc.ProcessCount = int(dbapi.GetProcessHeight())
			if i < 1 {
				oldProcesses = d.vc.ProcessCount
			}
			if d.vc.ProcessCount > 0 {
				updateProcesses(d, util.Max(oldProcesses-d.processIndex, config.ListSize))
			}
			vecty.Rerender(d)
		case i := <-d.refreshEnvelopes:
		envelopeLoop:
			for {
				// If many indices waiting in buffer, scan to last one.
				select {
				case i = <-d.refreshEnvelopes:
				default:
					break envelopeLoop
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

func updateVocdash(d *VocDashDashboardView) {
	updateHeights(d)
	if !d.disableEnvelopesUpdate {
		updateEnvelopes(d, util.Max(d.vc.EnvelopeHeight-d.envelopeIndex, config.ListSize))
	}
	if !d.disableEntitiesUpdate {
		updateEntities(d, util.Max(d.vc.EntityCount-d.entityIndex-1, config.ListSize-1))
	}
	if !d.disableProcessesUpdate {
		updateProcesses(d, util.Max(d.vc.ProcessCount-d.processIndex, config.ListSize))
	}
}

func updateEnvelopes(d *VocDashDashboardView, index int) {
	log.Infof("Getting envelopes from index %d", util.IntToString(index))
	list := dbapi.GetEnvelopeList(index)
	reverseEnvelopeList(&list)
	d.vc.EnvelopeList = list
}

func updateEntities(d *VocDashDashboardView, index int) {
	log.Infof("Getting entities from index %d", util.IntToString(index))
	list := dbapi.GetEntityList(index)
	reverseIDList(&list)
	d.vc.EntityIDs = list
}

func updateProcesses(d *VocDashDashboardView, index int) {
	log.Infof("Getting processes from index %d", util.IntToString(index))
	list := dbapi.GetProcessList(index)
	reverseIDList(&list)
	d.vc.ProcessIDs = list
}

func updateHeights(d *VocDashDashboardView) {
	d.vc.EnvelopeHeight = int(dbapi.GetEnvelopeHeight())
	d.vc.EntityCount = int(dbapi.GetEntityHeight())
	d.vc.ProcessCount = int(dbapi.GetProcessHeight())
}

func reverseEnvelopeList(list *[config.ListSize]*types.Envelope) {
	for i := len(list)/2 - 1; i >= 0; i-- {
		opp := len(list) - 1 - i
		list[i], list[opp] = list[opp], list[i]
	}
}

func reverseIDList(list *[config.ListSize]string) {
	for i := len(list)/2 - 1; i >= 0; i-- {
		opp := len(list) - 1 - i
		list[i], list[opp] = list[opp], list[i]
	}
}

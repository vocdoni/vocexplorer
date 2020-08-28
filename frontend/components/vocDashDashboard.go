package components

import (
	"fmt"
	"time"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/go-dvote/log"
	"gitlab.com/vocdoni/vocexplorer/client"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/dbapi"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/types"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// VocDashDashboardView renders the processes dashboard page
type VocDashDashboardView struct {
	vecty.Core
	gatewayConnected       bool
	serverConnected        bool
	gwClient               *client.Client
	envelopeIndex          int
	entityIndex            int
	processIndex           int
	disableEnvelopesUpdate bool
	disableEntitiesUpdate  bool
	disableProcessesUpdate bool
	refreshEnvelopes       chan int
	refreshEntities        chan int
	refreshProcesses       chan int
	rendered               *bool
	vc                     *client.VochainInfo
}

// Render renders the VocDashDashboardView component
func (dash *VocDashDashboardView) Render() vecty.ComponentOrHTML {
	*dash.rendered = true
	if dash != nil && dash.gwClient != nil && dash.vc != nil {
		return Container(
			renderGatewayConnectionBanner(dash.gatewayConnected),
			renderServerConnectionBanner(dash.serverConnected),
			elem.Section(
				bootstrap.Card(bootstrap.CardParams{
					Body: vecty.List{
						elem.Heading2(vecty.Text("Processes")),
						&ProcessListView{
							vochain:       dash.vc,
							refreshCh:     dash.refreshProcesses,
							disableUpdate: &dash.disableProcessesUpdate,
						},
					},
				}),
			),
			elem.Section(
				bootstrap.Card(bootstrap.CardParams{
					Body: vecty.List{
						elem.Heading2(vecty.Text("Entities")),
						&EntityListView{
							vochain:       dash.vc,
							refreshCh:     dash.refreshEntities,
							disableUpdate: &dash.disableEntitiesUpdate,
						},
					},
				}),
			),
			elem.Section(
				bootstrap.Card(bootstrap.CardParams{
					Body: vecty.List{
						elem.Heading2(vecty.Text("Envelopes")),
						&EnvelopeList{
							vochain:       dash.vc,
							refreshCh:     dash.refreshEnvelopes,
							disableUpdate: &dash.disableEnvelopesUpdate,
						},
					},
				}),
			),
		)
	}
	return &bootstrap.Alert{
		Contents: "Connecting to blockchain clients",
		Type:     "warning",
	}
}

// InitVocDashDashboardView initializes the vocdash page
func InitVocDashDashboardView(vc *client.VochainInfo, VocDashDashboardView *VocDashDashboardView, cfg *config.Cfg) *VocDashDashboardView {
	// gwClient, cancel := client.InitGateway(cfg.GatewayHost)
	// if gwClient == nil {
	// 	return VocDashDashboardView
	// }
	gwClient := store.GatewayClient
	VocDashDashboardView.gwClient = gwClient
	VocDashDashboardView.vc = vc
	VocDashDashboardView.refreshEnvelopes = make(chan int, 50)
	VocDashDashboardView.refreshProcesses = make(chan int, 50)
	VocDashDashboardView.refreshEntities = make(chan int, 50)
	VocDashDashboardView.disableEnvelopesUpdate = false
	store.Entities.PagChannel = make(chan int, 50)
	store.Processes.PagChannel = make(chan int, 50)
	VocDashDashboardView.refreshEntities = store.Entities.PagChannel
	VocDashDashboardView.refreshProcesses = store.Processes.PagChannel
	VocDashDashboardView.serverConnected = true
	VocDashDashboardView.gatewayConnected = true
	rendered := false
	VocDashDashboardView.rendered = &rendered
	BeforeUnload(func() {
		dispatcher.Dispatch(&actions.SignalRedirect{})
	})
	// go updateAndRenderVocDashDashboard(VocDashDashboardView, cancel, cfg)
	go updateAndRenderVocDashDashboard(VocDashDashboardView, cfg)
	return VocDashDashboardView
}

// func updateAndRenderVocDashDashboard(d *VocDashDashboardView, cancel context.CancelFunc, cfg *config.Cfg) {
func updateAndRenderVocDashDashboard(d *VocDashDashboardView, cfg *config.Cfg) {
	ticker := time.NewTicker(time.Duration(cfg.RefreshTime) * time.Second)
	for d != nil && !*d.rendered {
		fmt.Println("Not rendered yet")
	}
	updateVocdash(d)
	vecty.Rerender(d)
	time.Sleep(250 * time.Millisecond)
	vecty.Rerender(d)
	for {
		select {
		case <-store.RedirectChan:
			fmt.Println("Redirecting...")
			ticker.Stop()
			// d.gwClient.Close()
			return
		case <-ticker.C:
			updateVocdash(d)
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
			newVal, ok := dbapi.GetEntityHeight()
			if ok {
				d.vc.EntityCount = int(newVal)
			}
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
			newVal, ok := dbapi.GetProcessHeight()
			if ok {
				d.vc.ProcessCount = int(newVal)
			}
			if i < 1 {
				oldProcesses = d.vc.ProcessCount
			}
			if d.vc.ProcessCount > 0 {
				updateProcesses(d, util.Max(oldProcesses-d.processIndex, config.ListSize))
				client.UpdateProcessResults(d.gwClient, d.vc)
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
			newVal, ok := dbapi.GetEnvelopeHeight()
			if ok {
				d.vc.EnvelopeHeight = int(newVal)
			}
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
	if d.gwClient.Conn.Ping(d.gwClient.Ctx) != nil {
		d.gatewayConnected = false
	} else {
		d.gatewayConnected = true
	}
	if !dbapi.Ping() {
		d.serverConnected = false
	} else {
		d.serverConnected = true
	}
	updateHeights(d)
	if !d.disableEnvelopesUpdate {
		updateEnvelopes(d, util.Max(d.vc.EnvelopeHeight-d.envelopeIndex, config.ListSize))
	}
	if !d.disableEntitiesUpdate {
		updateEntities(d, util.Max(d.vc.EntityCount-d.entityIndex-1, config.ListSize-1))
	}
	if !d.disableProcessesUpdate {
		updateProcesses(d, util.Max(d.vc.ProcessCount-d.processIndex, config.ListSize))
		client.UpdateProcessResults(d.gwClient, d.vc)
	}
}

func updateEnvelopes(d *VocDashDashboardView, index int) {
	log.Infof("Getting envelopes from index %d", util.IntToString(index))
	list, ok := dbapi.GetEnvelopeList(index)
	if ok {
		reverseEnvelopeList(&list)
		d.vc.EnvelopeList = list
	}
}

func updateEntities(d *VocDashDashboardView, index int) {
	log.Infof("Getting entities from index %d", util.IntToString(index))
	list, ok := dbapi.GetEntityList(index)
	if ok {
		reverseIDList(&list)
		d.vc.EntityIDs = list
	}
}

func updateProcesses(d *VocDashDashboardView, index int) {
	log.Infof("Getting processes from index %d", util.IntToString(index))
	list, ok := dbapi.GetProcessList(index)
	if ok {
		reverseIDList(&list)
		d.vc.ProcessIDs = list
	}
	newVal, ok := dbapi.GetProcessEnvelopeHeightMap()
	if ok {
		d.vc.EnvelopeHeights = newVal
	}
	newVal, ok = dbapi.GetEntityProcessHeightMap()
	if ok {
		d.vc.ProcessHeights = newVal
	}
}

func updateHeights(d *VocDashDashboardView) {
	newVal, ok := dbapi.GetEnvelopeHeight()
	if ok {
		d.vc.EnvelopeHeight = int(newVal)
	}
	newVal, ok = dbapi.GetEntityHeight()
	if ok {
		d.vc.EntityCount = int(newVal)
	}
	newVal, ok = dbapi.GetProcessHeight()
	if ok {
		d.vc.ProcessCount = int(newVal)
	}
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

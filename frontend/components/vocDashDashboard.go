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
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/types"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// VocDashDashboardView renders the processes dashboard page
type VocDashDashboardView struct {
	vecty.Core
	GatewayConnected       bool
	ServerConnected        bool
	EnvelopeIndex          int
	EntityIndex            int
	ProcessIndex           int
	QuitCh                 chan struct{}
	DisableEnvelopesUpdate bool
	DisableEntitiesUpdate  bool
	DisableProcessesUpdate bool
	RefreshEnvelopes       chan int
	RefreshEntities        chan int
	RefreshProcesses       chan int
	Vc                     *client.VochainInfo
	Rendered               *bool
}

// Render renders the VocDashDashboardView component
func (dash *VocDashDashboardView) Render() vecty.ComponentOrHTML {
	*dash.Rendered = true
	if dash != nil && store.Vochain != nil && dash.Vc != nil {
		return Container(
			renderGatewayConnectionBanner(dash.GatewayConnected),
			renderServerConnectionBanner(dash.ServerConnected),
			elem.Section(
				bootstrap.Card(bootstrap.CardParams{
					Body: vecty.List{
						elem.Heading2(vecty.Text("Processes")),
						&ProcessListView{
							vochain:       dash.Vc,
							refreshCh:     dash.RefreshProcesses,
							disableUpdate: &dash.DisableProcessesUpdate,
						},
					},
				}),
			),
			elem.Section(
				bootstrap.Card(bootstrap.CardParams{
					Body: vecty.List{
						elem.Heading2(vecty.Text("Entities")),
						&EntityListView{
							vochain:       dash.Vc,
							refreshCh:     dash.RefreshEntities,
							disableUpdate: &dash.DisableEntitiesUpdate,
						},
					},
				}),
			),
			elem.Section(
				bootstrap.Card(bootstrap.CardParams{
					Body: vecty.List{
						elem.Heading2(vecty.Text("Envelopes")),
						&EnvelopeList{
							vochain:       dash.Vc,
							refreshCh:     dash.RefreshEnvelopes,
							disableUpdate: &dash.DisableEnvelopesUpdate,
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

// // InitVocDashDashboardView initializes the vocdash page
// func InitVocDashDashboardView(vc *client.VochainInfo, VocDashDashboardView *VocDashDashboardView, cfg *config.Cfg) *VocDashDashboardView {
// 	VocDashDashboardView.vc = vc
// 	VocDashDashboardView.quitCh = make(chan struct{})
// 	VocDashDashboardView.refreshEnvelopes = make(chan int, 50)
// 	VocDashDashboardView.refreshProcesses = make(chan int, 50)
// 	VocDashDashboardView.refreshEntities = make(chan int, 50)
// 	VocDashDashboardView.disableEnvelopesUpdate = false
// 	store.Entities.PagChannel = make(chan int, 50)
// 	store.Processes.PagChannel = make(chan int, 50)
// 	VocDashDashboardView.refreshEntities = store.Entities.PagChannel
// 	VocDashDashboardView.refreshProcesses = store.Processes.PagChannel
// 	VocDashDashboardView.serverConnected = true
// 	VocDashDashboardView.gatewayConnected = true

// 	go updateAndRenderVocDashDashboard(VocDashDashboardView, cfg)
// 	return VocDashDashboardView
// }

// UpdateAndRenderVocDashDashboard continuously updates the information needed by the vocdash dashboard
func UpdateAndRenderVocDashDashboard(d *VocDashDashboardView, cfg *config.Cfg) {
	for !*d.Rendered {
		fmt.Println("Not rendered")
		time.Sleep(20 * time.Millisecond)
	}
	ticker := time.NewTicker(time.Duration(cfg.RefreshTime) * time.Second)
	updateVocdash(d)
	vecty.Rerender(d)
	time.Sleep(250 * time.Millisecond)
	vecty.Rerender(d)
	for {
		select {
		case <-d.QuitCh:
			ticker.Stop()
			// store.Vochain.Close()
			fmt.Println("Gateway connection closed")
			return
		case <-ticker.C:
			updateVocdash(d)
			vecty.Rerender(d)
		case i := <-d.RefreshEntities:
		entityLoop:
			for {
				// If many indices waiting in buffer, scan to last one.
				select {
				case i = <-d.RefreshEntities:
				default:
					break entityLoop
				}
			}
			d.EntityIndex = i
			oldEntities := d.Vc.EntityCount
			newVal, ok := dbapi.GetEntityHeight()
			if ok {
				d.Vc.EntityCount = int(newVal)
			}
			if i < 1 {
				oldEntities = d.Vc.EntityCount
			}
			if d.Vc.EntityCount > 0 {
				updateEntities(d, util.Max(oldEntities-d.EntityIndex-1, config.ListSize-1))
			}
			vecty.Rerender(d)
		case i := <-d.RefreshProcesses:
		processLoop:
			for {
				// If many indices waiting in buffer, scan to last one.
				select {
				case i = <-d.RefreshProcesses:
				default:
					break processLoop

				}
			}
			d.ProcessIndex = i
			oldProcesses := d.Vc.ProcessCount
			newVal, ok := dbapi.GetProcessHeight()
			if ok {
				d.Vc.ProcessCount = int(newVal)
			}
			if i < 1 {
				oldProcesses = d.Vc.ProcessCount
			}
			if d.Vc.ProcessCount > 0 {
				updateProcesses(d, util.Max(oldProcesses-d.ProcessIndex, config.ListSize))
				client.UpdateProcessResults(store.Vochain, d.Vc)
			}
			vecty.Rerender(d)
		case i := <-d.RefreshEnvelopes:
		envelopeLoop:
			for {
				// If many indices waiting in buffer, scan to last one.
				select {
				case i = <-d.RefreshEnvelopes:
				default:
					break envelopeLoop
				}
			}
			d.EnvelopeIndex = i
			oldEnvelopes := d.Vc.EnvelopeHeight
			newVal, ok := dbapi.GetEnvelopeHeight()
			if ok {
				d.Vc.EnvelopeHeight = int(newVal)
			}
			if i < 1 {
				oldEnvelopes = d.Vc.EnvelopeHeight
			}
			if d.Vc.EnvelopeHeight > 0 {
				updateEnvelopes(d, util.Max(oldEnvelopes-d.EnvelopeIndex, config.ListSize))
			}
			vecty.Rerender(d)
		}
	}
}

func updateVocdash(d *VocDashDashboardView) {
	if store.Vochain.Conn.Ping(store.Vochain.Ctx) != nil {
		d.GatewayConnected = false
	} else {
		d.GatewayConnected = true
	}
	if !dbapi.Ping() {
		d.ServerConnected = false
	} else {
		d.ServerConnected = true
	}
	updateHeights(d)
	if !d.DisableEnvelopesUpdate {
		updateEnvelopes(d, util.Max(d.Vc.EnvelopeHeight-d.EnvelopeIndex, config.ListSize))
	}
	if !d.DisableEntitiesUpdate {
		updateEntities(d, util.Max(d.Vc.EntityCount-d.EntityIndex-1, config.ListSize-1))
	}
	if !d.DisableProcessesUpdate {
		updateProcesses(d, util.Max(d.Vc.ProcessCount-d.ProcessIndex, config.ListSize))
		client.UpdateProcessResults(store.Vochain, d.Vc)
	}
}

func updateEnvelopes(d *VocDashDashboardView, index int) {
	log.Infof("Getting envelopes from index %d", util.IntToString(index))
	list, ok := dbapi.GetEnvelopeList(index)
	if ok {
		reverseEnvelopeList(&list)
		d.Vc.EnvelopeList = list
	}
}

func updateEntities(d *VocDashDashboardView, index int) {
	log.Infof("Getting entities from index %d", util.IntToString(index))
	list, ok := dbapi.GetEntityList(index)
	if ok {
		reverseIDList(&list)
		d.Vc.EntityIDs = list
	}
}

func updateProcesses(d *VocDashDashboardView, index int) {
	log.Infof("Getting processes from index %d", util.IntToString(index))
	list, ok := dbapi.GetProcessList(index)
	if ok {
		reverseIDList(&list)
		d.Vc.ProcessIDs = list
	}
	newVal, ok := dbapi.GetProcessEnvelopeHeightMap()
	if ok {
		d.Vc.EnvelopeHeights = newVal
	}
	newVal, ok = dbapi.GetEntityProcessHeightMap()
	if ok {
		d.Vc.ProcessHeights = newVal
	}
}

func updateHeights(d *VocDashDashboardView) {
	newVal, ok := dbapi.GetEnvelopeHeight()
	if ok {
		d.Vc.EnvelopeHeight = int(newVal)
	}
	newVal, ok = dbapi.GetEntityHeight()
	if ok {
		d.Vc.EntityCount = int(newVal)
	}
	newVal, ok = dbapi.GetProcessHeight()
	if ok {
		d.Vc.ProcessCount = int(newVal)
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

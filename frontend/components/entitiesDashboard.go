package components

import (
	"context"
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
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// EntitiesDashboardView renders the entities dashboard page
type EntitiesDashboardView struct {
	vecty.Core
	gatewayConnected       bool
	serverConnected        bool
	gwClient               *client.Client
	entity                 *client.EntityInfo
	entityID               string
	processIndex           int
	disableProcessesUpdate bool
	quitCh                 chan struct{}
	refreshCh              chan int
}

//EntitiesTab is the tab component for entities
type EntitiesTab struct {
	*Tab
}

func (e *EntitiesTab) dispatch() interface{} {
	return &actions.EntitiesTabChange{
		Tab: e.alias(),
	}
}

func (e *EntitiesTab) store() string {
	return store.Entities.Tab
}

// Render renders the EntitiesDashboardView component
func (dash *EntitiesDashboardView) Render() vecty.ComponentOrHTML {
	if dash == nil || dash.gwClient == nil || dash.entity == nil {
		return Container(&bootstrap.Alert{
			Type:     "warning",
			Contents: "Connecting to blockchain client",
		})
	}

	return Container(
		renderGatewayConnectionBanner(dash.gatewayConnected),
		renderServerConnectionBanner(dash.serverConnected),
		elem.Section(
			vecty.Markup(vecty.Class("details-view", "no-column")),
			elem.Div(
				vecty.Markup(vecty.Class("row")),
				elem.Div(
					vecty.Markup(vecty.Class("main-column")),
					bootstrap.Card(bootstrap.CardParams{
						Body: dash.EntityDetails(),
					}),
				),
			),
		),
		elem.Section(
			vecty.Markup(vecty.Class("row")),
			elem.Div(
				vecty.Markup(vecty.Class("col-12")),
				bootstrap.Card(bootstrap.CardParams{
					Body: &EntityProcessListView{
						entity:        dash.entity,
						refreshCh:     dash.refreshCh,
						disableUpdate: &dash.disableProcessesUpdate,
					},
				}),
			),
		),
	)
}

//EntityDetails renders the details of a single entity
func (dash *EntitiesDashboardView) EntityDetails() vecty.List {
	return vecty.List{
		elem.Heading1(
			vecty.Text("Entity details"),
		),
		elem.Heading2(vecty.Text(dash.entityID)),
		elem.Anchor(
			vecty.Markup(vecty.Class("hash")),
			vecty.Markup(vecty.Attribute("href", "https://manage.vocdoni.net/entities/#/0x"+dash.entityID)),
			vecty.Text("Entity Manager Page"),
		),
	}
}

// InitEntitiesDashboardView initializes the entities dashboard view
func InitEntitiesDashboardView(entity *client.EntityInfo, EntitiesDashboardView *EntitiesDashboardView, entityID string, cfg *config.Cfg) *EntitiesDashboardView {
	gwClient, cancel := client.InitGateway(cfg.GatewayHost)
	if gwClient == nil {
		return EntitiesDashboardView
	}
	EntitiesDashboardView.gwClient = gwClient
	EntitiesDashboardView.entity = entity
	EntitiesDashboardView.entityID = entityID
	EntitiesDashboardView.quitCh = make(chan struct{})
	EntitiesDashboardView.refreshCh = make(chan int, 50)
	EntitiesDashboardView.serverConnected = true
	EntitiesDashboardView.gatewayConnected = true
	BeforeUnload(func() {
		close(EntitiesDashboardView.quitCh)
	})
	go updateAndRenderEntitiesDashboard(EntitiesDashboardView, cancel, entityID, cfg)
	return EntitiesDashboardView
}

func updateAndRenderEntitiesDashboard(d *EntitiesDashboardView, cancel context.CancelFunc, entityID string, cfg *config.Cfg) {
	ticker := time.NewTicker(time.Duration(cfg.RefreshTime) * time.Second)
	updateEntityProcesses(d, util.Max(d.entity.ProcessCount-d.processIndex, config.ListSize))
	vecty.Rerender(d)
	for {
		select {
		case <-d.quitCh:
			ticker.Stop()
			d.gwClient.Close()
			fmt.Println("Gateway connection closed")
			return
		case <-ticker.C:
			updateEntityProcesses(d, util.Max(d.entity.ProcessCount-d.processIndex, config.ListSize))
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
			d.processIndex = i
			oldProcesses := d.entity.ProcessCount
			newHeight, _ := dbapi.GetEntityProcessHeight(entityID)
			d.entity.ProcessCount = int(newHeight)
			if i < 1 {
				oldProcesses = d.entity.ProcessCount
			}
			updateEntityProcesses(d, util.Max(oldProcesses-d.processIndex, config.ListSize))
			vecty.Rerender(d)
		}
	}
}

func updateEntityProcesses(d *EntitiesDashboardView, index int) {
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
	newCount, ok := dbapi.GetEntityProcessHeight(d.entityID)
	if ok {
		d.entity.ProcessCount = int(newCount)
	}
	if d.entity.ProcessCount > 0 && !d.disableProcessesUpdate {
		log.Infof("Getting processes from entity %s, index %d", d.entityID, util.IntToString(index))
		list, ok := dbapi.GetProcessListByEntity(index, d.entityID)
		if ok {
			reverseIDList(&list)
			d.entity.ProcessIDs = list
		}
		newMap, ok := dbapi.GetProcessEnvelopeHeightMap()
		if ok {
			d.entity.EnvelopeHeights = newMap
		}
		client.UpdateAuxEntityInfo(d.gwClient, d.entity)
	}
}

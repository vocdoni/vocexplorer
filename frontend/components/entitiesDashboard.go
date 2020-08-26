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
	gwClient               *client.Client
	entity                 *client.EntityInfo
	entityID               string
	processIndex           int
	disableProcessesUpdate bool
	quitCh                 chan struct{}
	refreshCh              chan int
}

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

func (e *EntitiesDashboardView) EntityDetails() vecty.List {
	return vecty.List{
		elem.Heading1(
			vecty.Text("Entity details"),
		),
		elem.Heading2(vecty.Text(e.entityID)),
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
			d.entity.ProcessCount = int(dbapi.GetEntityProcessHeight(entityID))
			if i < 1 {
				oldProcesses = d.entity.ProcessCount
			}
			updateEntityProcesses(d, util.Max(oldProcesses-d.processIndex, config.ListSize))
			vecty.Rerender(d)
		}
	}
}

func updateEntityProcesses(d *EntitiesDashboardView, index int) {
	d.entity.ProcessCount = int(dbapi.GetEntityProcessHeight(d.entityID))
	if d.entity.ProcessCount > 0 && !d.disableProcessesUpdate {
		log.Infof("Getting processes from entity %s, index %d", d.entityID, util.IntToString(index))
		list := dbapi.GetProcessListByEntity(index, d.entityID)
		reverseIDList(&list)
		d.entity.ProcessIDs = list
		d.entity.EnvelopeHeights = dbapi.GetProcessEnvelopeHeightMap()
		client.UpdateAuxEntityInfo(d.gwClient, d.entity)
	}
}

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

// EntitiesDashboardView renders the entities dashboard page
type EntitiesDashboardView struct {
	vecty.Core
	entity    *client.EntityInfo
	entityID  string
	gwClient  *client.Client
	quitCh    chan struct{}
	refreshCh chan bool
}

// Render renders the EntitiesDashboardView component
func (dash *EntitiesDashboardView) Render() vecty.ComponentOrHTML {
	if dash != nil && dash.gwClient != nil && dash.entity != nil {
		return Container(
			elem.Heading4(vecty.Text(
				fmt.Sprintf("Entity %s", dash.entityID),
			)),
			vecty.Markup(vecty.Class("info-pane")),
			&ProcessListView{
				entity:    dash.entity,
				refreshCh: dash.refreshCh,
			},
		)
	}
	return Container(vecty.Text("Connecting to blockchain client"))
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
	EntitiesDashboardView.refreshCh = make(chan bool, 20)
	BeforeUnload(func() {
		close(EntitiesDashboardView.quitCh)
	})
	go updateAndRenderEntitiesDashboard(EntitiesDashboardView, cancel, entityID, cfg)
	return EntitiesDashboardView
}

func updateAndRenderEntitiesDashboard(d *EntitiesDashboardView, cancel context.CancelFunc, entityID string, cfg *config.Cfg) {
	ticker := time.NewTicker(time.Duration(cfg.RefreshTime) * time.Second)
	// Wait for data structs to load
	for d == nil || d.entity == nil {
	}
	// TODO change to accept real index
	client.UpdateEntitiesDashboardInfo(d.gwClient, d.entity, entityID, 0)
	vecty.Rerender(d)
	time.Sleep(250 * time.Millisecond)
	client.UpdateAuxEntityInfo(d.gwClient, d.entity)
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
			client.UpdateEntitiesDashboardInfo(d.gwClient, d.entity, entityID, 0)
			client.UpdateAuxEntityInfo(d.gwClient, d.entity)
			vecty.Rerender(d)
		case <-d.refreshCh:
			//TODO: update to  use real index
			client.UpdateEntitiesDashboardInfo(d.gwClient, d.entity, entityID, 0)
			client.UpdateAuxEntityInfo(d.gwClient, d.entity)
			vecty.Rerender(d)
		}
	}
}

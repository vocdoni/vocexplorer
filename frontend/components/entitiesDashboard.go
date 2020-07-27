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
	entityID string
	entity   *client.EntityInfo
	gwClient *client.Client
	quitCh   chan struct{}
}

// Render renders the EntitiesDashboardView component
func (dash *EntitiesDashboardView) Render() vecty.ComponentOrHTML {
	if dash != nil && dash.gwClient != nil && dash.entity != nil {
		return elem.Div(
			elem.Main(
				elem.Heading4(vecty.Text("Entity "+dash.entityID)),
				vecty.Markup(vecty.Class("info-pane")),
				&ProcessListView{
					entity: dash.entity,
				},
			),
		)
	}
	return vecty.Text("Connecting to blockchain client")
}

func initEntitiesDashboardView(entity *client.EntityInfo, EntitiesDashboardView *EntitiesDashboardView, entityID string) *EntitiesDashboardView {
	gwClient, cancel := InitGateway()
	if gwClient == nil {
		return EntitiesDashboardView
	}
	EntitiesDashboardView.gwClient = gwClient
	EntitiesDashboardView.entity = entity
	EntitiesDashboardView.entityID = entityID
	EntitiesDashboardView.quitCh = make(chan struct{})
	BeforeUnload(func() {
		close(EntitiesDashboardView.quitCh)
	})
	go updateAndRenderEntitiesDashboard(EntitiesDashboardView, cancel, entityID)
	return EntitiesDashboardView
}

func updateAndRenderEntitiesDashboard(d *EntitiesDashboardView, cancel context.CancelFunc, entityID string) {
	ticker := time.NewTicker(config.RefreshTime * time.Second)
	// Wait for data structs to load
	for d == nil || d.entity == nil {
	}
	client.UpdateEntitiesDashboardInfo(d.gwClient, d.entity, entityID)
	vecty.Rerender(d)
	for {
		select {
		case <-d.quitCh:
			ticker.Stop()
			d.gwClient.Close()
			fmt.Println("Gateway connection closed")
			return
		case <-ticker.C:
			client.UpdateEntitiesDashboardInfo(d.gwClient, d.entity, entityID)
			client.UpdateAuxEntityInfo(d.gwClient, d.entity)
			vecty.Rerender(d)
		}
	}
}

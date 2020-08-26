package components

import (
	"context"
	"fmt"
	"time"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/client"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
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
					Body: &ProcessListView{
						entity:    dash.entity,
						refreshCh: dash.refreshCh,
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
	EntitiesDashboardView.refreshCh = make(chan bool, 20)
	BeforeUnload(func() {
		close(EntitiesDashboardView.quitCh)
	})
	go updateAndRenderEntitiesDashboard(EntitiesDashboardView, cancel, entityID, cfg)
	return EntitiesDashboardView
}

func updateAndRenderEntitiesDashboard(d *EntitiesDashboardView, cancel context.CancelFunc, entityID string, cfg *config.Cfg) {
	ticker := time.NewTicker(time.Duration(cfg.RefreshTime) * time.Second)
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

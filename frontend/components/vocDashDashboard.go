package components

import (
	"fmt"
	"time"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/api"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/types"
	"gitlab.com/vocdoni/vocexplorer/update"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// VocDashDashboardView renders the processes dashboard page
type VocDashDashboardView struct {
	vecty.Core
	vecty.Mounter
	Rendered      bool
	EntityIndex   int
	EnvelopeIndex int
	ProcessIndex  int
}

// Mount is called after the component renders to signal that it can be rerendered safely
func (dash *VocDashDashboardView) Mount() {
	if !dash.Rendered {
		dash.Rendered = true
		vecty.Rerender(dash)
	}
}

// Render renders the VocDashDashboardView component
func (dash *VocDashDashboardView) Render() vecty.ComponentOrHTML {
	if !dash.Rendered {
		return elem.Div(vecty.Text("Loading..."))
	}
	if dash != nil && store.GatewayClient != nil && store.TendermintClient != nil {
		return Container(
			renderGatewayConnectionBanner(),
			renderServerConnectionBanner(),
			elem.Section(
				bootstrap.Card(bootstrap.CardParams{
					Body: vecty.List{
						elem.Heading2(vecty.Text("Processes")),
						&ProcessListView{},
					},
				}),
			),
			elem.Section(
				bootstrap.Card(bootstrap.CardParams{
					Body: vecty.List{
						elem.Heading2(vecty.Text("Entities")),
						&EntityListView{},
					},
				}),
			),
			elem.Section(
				bootstrap.Card(bootstrap.CardParams{
					Body: vecty.List{
						elem.Heading2(vecty.Text("Envelopes")),
						&EnvelopeList{},
					},
				}),
			),
		)
	}
	return elem.Div(
		&bootstrap.Alert{
			Contents: "Connecting to blockchain clients",
			Type:     "warning",
		},
	)
}

// UpdateAndRenderVocDashDashboard continuously updates the information needed by the vocdash dashboard
func UpdateAndRenderVocDashDashboard(d *VocDashDashboardView) {
	actions.EnableUpdates()
	ticker := time.NewTicker(time.Duration(store.Config.RefreshTime) * time.Second)
	updateVocdash(d)
	for {
		select {
		case <-store.RedirectChan:
			fmt.Println("Redirecting...")
			ticker.Stop()
			return
		case <-ticker.C:
			updateVocdash(d)
		case i := <-store.Entities.Pagination.PagChannel:
		entityLoop:
			for {
				// If many indices waiting in buffer, scan to last one.
				select {
				case i = <-store.Entities.Pagination.PagChannel:
				default:
					break entityLoop
				}
			}
			d.EntityIndex = i
			oldEntities := store.Entities.Count
			newVal, ok := api.GetEntityHeight()
			if ok {
				dispatcher.Dispatch(&actions.SetEntityCount{Count: int(newVal)})
			}
			if i < 1 {
				oldEntities = store.Entities.Count
			}
			if store.Entities.Count > 0 {
				updateEntities(d, util.Max(oldEntities-d.EntityIndex-1, config.ListSize-1))
			}
		case i := <-store.Processes.Pagination.PagChannel:
		processLoop:
			for {
				// If many indices waiting in buffer, scan to last one.
				select {
				case i = <-store.Processes.Pagination.PagChannel:
				default:
					break processLoop

				}
			}
			d.ProcessIndex = i
			oldProcesses := store.Processes.Count
			newVal, ok := api.GetProcessHeight()
			if ok {
				dispatcher.Dispatch(&actions.SetProcessCount{Count: int(newVal)})
			}
			if i < 1 {
				oldProcesses = store.Processes.Count
			}
			if store.Processes.Count > 0 {
				updateProcesses(d, util.Max(oldProcesses-d.ProcessIndex, config.ListSize))
				update.ProcessResults()
			}
		case i := <-store.Envelopes.Pagination.PagChannel:
		envelopeLoop:
			for {
				// If many indices waiting in buffer, scan to last one.
				select {
				case i = <-store.Envelopes.Pagination.PagChannel:
				default:
					break envelopeLoop
				}
			}
			d.EnvelopeIndex = i
			oldEnvelopes := store.Envelopes.Count
			newVal, ok := api.GetEnvelopeHeight()
			if ok {
				dispatcher.Dispatch(&actions.SetEnvelopeCount{Count: int(newVal)})
			}
			if i < 1 {
				oldEnvelopes = store.Envelopes.Count
			}
			if store.Envelopes.Count > 0 {
				updateEnvelopes(d, util.Max(oldEnvelopes-d.EnvelopeIndex, config.ListSize))
			}
		}
	}
}

func updateVocdash(d *VocDashDashboardView) {
	dispatcher.Dispatch(&actions.GatewayConnected{Connected: api.PingGateway(store.Config.GatewayHost)})
	dispatcher.Dispatch(&actions.ServerConnected{Connected: api.Ping()})
	actions.UpdateCounts()
	if !store.Envelopes.Pagination.DisableUpdate {
		updateEnvelopes(d, util.Max(store.Envelopes.Count-d.EnvelopeIndex, config.ListSize))
	}
	if !store.Entities.Pagination.DisableUpdate {
		updateEntities(d, util.Max(store.Entities.Count-d.EntityIndex-1, config.ListSize-1))
	}
	if !store.Processes.Pagination.DisableUpdate {
		updateProcesses(d, util.Max(store.Processes.Count-d.ProcessIndex, config.ListSize))
		update.ProcessResults()
	}
}

func updateEnvelopes(d *VocDashDashboardView, index int) {
	fmt.Printf("Getting envelopes from index %d\n", index)
	list, ok := api.GetEnvelopeList(index)
	if ok {
		reverseEnvelopeList(&list)
		dispatcher.Dispatch(&actions.SetEnvelopeList{EnvelopeList: list})
	}
}

func updateEntities(d *VocDashDashboardView, index int) {
	fmt.Printf("Getting entities from index %d\n", index)
	list, ok := api.GetEntityList(index)
	if ok {
		reverseIDList(&list)
		dispatcher.Dispatch(&actions.SetEntityIDs{EntityIDs: list})
	}
}

func updateProcesses(d *VocDashDashboardView, index int) {
	fmt.Printf("Getting processes from index %d\n", index)
	list, ok := api.GetProcessList(index)
	if ok {
		reverseIDList(&list)
		dispatcher.Dispatch(&actions.SetProcessIDs{ProcessIDs: list})
	}
	newVal, ok := api.GetProcessEnvelopeHeightMap()
	if ok {
		dispatcher.Dispatch(&actions.SetEnvelopeHeights{EnvelopeHeights: newVal})
	}
	newVal, ok = api.GetEntityProcessHeightMap()
	if ok {
		dispatcher.Dispatch(&actions.SetProcessHeights{ProcessHeights: newVal})
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

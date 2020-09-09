package components

import (
	"fmt"
	"log"
	"time"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/api"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/frontend/update"
	"gitlab.com/vocdoni/vocexplorer/proto"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// ParticipationDashboardView renders the processes dashboard page
type ParticipationDashboardView struct {
	vecty.Core
	vecty.Mounter
	Rendered bool
}

// Mount is called after the component renders to signal that it can be rerendered safely
func (dash *ParticipationDashboardView) Mount() {
	if !dash.Rendered {
		dash.Rendered = true
		vecty.Rerender(dash)
	}
}

// Render renders the ParticipationDashboardView component
func (dash *ParticipationDashboardView) Render() vecty.ComponentOrHTML {
	if !dash.Rendered {
		return LoadingBar()
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

// UpdateAndRenderParticipationDashboard continuously updates the information needed by the participation dashboard
func UpdateAndRenderParticipationDashboard(d *ParticipationDashboardView) {
	dispatcher.Dispatch(&actions.EnableAllUpdates{})

	ticker := time.NewTicker(time.Duration(store.Config.RefreshTime) * time.Second)
	updateParticipation(d)
	for {
		select {
		case <-store.RedirectChan:
			fmt.Println("Redirecting...")
			ticker.Stop()
			return
		case <-ticker.C:
			updateParticipation(d)
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
			dispatcher.Dispatch(&actions.EntitiesIndexChange{Index: i})
			oldEntities := store.Entities.Count
			newVal, ok := api.GetEntityCount()
			if ok {
				dispatcher.Dispatch(&actions.SetEntityCount{Count: int(newVal)})
			}
			if i < 1 {
				oldEntities = store.Entities.Count
			}
			if store.Entities.Count > 0 {
				updateEntities(d, util.Max(oldEntities-store.Entities.Pagination.Index, 1))
			}
		case search := <-store.Entities.Pagination.SearchChannel:
		entitySearch:
			for {
				// If many indices waiting in buffer, scan to last one.
				select {
				case search = <-store.Entities.Pagination.SearchChannel:
				default:
					break entitySearch
				}
			}
			log.Println("search: " + search)
			dispatcher.Dispatch(&actions.EntitiesIndexChange{Index: 0})
			list, ok := api.GetEntitySearch(search)
			if ok {
				reverseIDList(&list)
				dispatcher.Dispatch(&actions.SetEntityIDs{EntityIDs: list})
			} else {
				dispatcher.Dispatch(&actions.SetEntityIDs{EntityIDs: [config.ListSize]string{}})
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
			dispatcher.Dispatch(&actions.ProcessesIndexChange{Index: i})
			oldProcesses := store.Processes.Count
			newVal, ok := api.GetProcessCount()
			if ok {
				dispatcher.Dispatch(&actions.SetProcessCount{Count: int(newVal)})
			}
			if i < 1 {
				oldProcesses = store.Processes.Count
			}
			if store.Processes.Count > 0 {
				updateProcesses(d, util.Max(oldProcesses-store.Processes.Pagination.Index, 1))
				update.ProcessResults()
			}
		case search := <-store.Processes.Pagination.SearchChannel:
		processSearch:
			for {
				// If many indices waiting in buffer, scan to last one.
				select {
				case search = <-store.Processes.Pagination.SearchChannel:
				default:
					break processSearch
				}
			}
			log.Println("search: " + search)
			dispatcher.Dispatch(&actions.ProcessesIndexChange{Index: 0})
			list, ok := api.GetProcessSearch(search)
			if ok {
				dispatcher.Dispatch(&actions.SetProcessIDs{ProcessIDs: list})
			} else {
				dispatcher.Dispatch(&actions.SetProcessIDs{ProcessIDs: [config.ListSize]string{}})
			}
			update.ProcessResults()
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
			dispatcher.Dispatch(&actions.EnvelopesIndexChange{Index: i})
			oldEnvelopes := store.Envelopes.Count
			newVal, ok := api.GetEnvelopeCount()
			if ok {
				dispatcher.Dispatch(&actions.SetEnvelopeCount{Count: int(newVal)})
			}
			if i < 1 {
				oldEnvelopes = store.Envelopes.Count
			}
			if store.Envelopes.Count > 0 {
				updateEnvelopes(d, util.Max(oldEnvelopes-store.Envelopes.Pagination.Index, 1))
			}
		case search := <-store.Envelopes.Pagination.SearchChannel:
		envelopeSearch:
			for {
				// If many indices waiting in buffer, scan to last one.
				select {
				case search = <-store.Envelopes.Pagination.SearchChannel:
				default:
					break envelopeSearch
				}
			}
			log.Println("search: " + search)
			dispatcher.Dispatch(&actions.EnvelopesIndexChange{Index: 0})
			list, ok := api.GetEnvelopeSearch(search)
			if ok {
				reverseEnvelopeList(&list)
				dispatcher.Dispatch(&actions.SetEnvelopeList{EnvelopeList: list})
			} else {
				dispatcher.Dispatch(&actions.SetEnvelopeList{EnvelopeList: [config.ListSize]*proto.Envelope{}})
			}
		}
	}
}

func updateParticipation(d *ParticipationDashboardView) {
	dispatcher.Dispatch(&actions.GatewayConnected{Connected: store.GatewayClient.Ping()})
	dispatcher.Dispatch(&actions.ServerConnected{Connected: api.PingServer()})
	actions.UpdateCounts()
	if !store.Envelopes.Pagination.DisableUpdate {
		updateEnvelopes(d, util.Max(store.Envelopes.Count-store.Envelopes.Pagination.Index, 1))
	}
	if !store.Entities.Pagination.DisableUpdate {
		updateEntities(d, util.Max(store.Entities.Count-store.Entities.Pagination.Index, 1))
	}
	if !store.Processes.Pagination.DisableUpdate {
		updateProcesses(d, util.Max(store.Processes.Count-store.Processes.Pagination.Index, 1))
		update.ProcessResults()
	}
}

func updateEnvelopes(d *ParticipationDashboardView, index int) {
	fmt.Printf("Getting envelopes from index %d\n", index)
	list, ok := api.GetEnvelopeList(index)
	if ok {
		reverseEnvelopeList(&list)
		dispatcher.Dispatch(&actions.SetEnvelopeList{EnvelopeList: list})
	}
}

func updateEntities(d *ParticipationDashboardView, index int) {
	index--
	fmt.Printf("Getting entities from index %d\n", index)
	list, ok := api.GetEntityList(index)
	if ok {
		reverseIDList(&list)
		dispatcher.Dispatch(&actions.SetEntityIDs{EntityIDs: list})
	}
}

func updateProcesses(d *ParticipationDashboardView, index int) {
	// index--
	fmt.Printf("Getting processes from index %d\n", index)
	list, ok := api.GetProcessList(index)
	if ok {
		reverseIDList(&list)
		dispatcher.Dispatch(&actions.SetProcessIDs{ProcessIDs: list})
	}
	newVal, ok := api.GetProcessEnvelopeCountMap()
	if ok {
		dispatcher.Dispatch(&actions.SetEnvelopeHeights{EnvelopeHeights: newVal})
	}
	newVal, ok = api.GetEntityProcessCountMap()
	if ok {
		dispatcher.Dispatch(&actions.SetProcessHeights{ProcessHeights: newVal})
	}
}

func reverseEnvelopeList(list *[config.ListSize]*proto.Envelope) {
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

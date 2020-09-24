package components

import (
	"fmt"
	"log"
	"time"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
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

// EnvelopesDashboardView renders the envelopes dashboard page
type EnvelopesDashboardView struct {
	vecty.Core
	vecty.Mounter
	Rendered bool
}

// Mount is called after the component renders to signal that it can be rerendered safely
func (dash *EnvelopesDashboardView) Mount() {
	if !dash.Rendered {
		dash.Rendered = true
		vecty.Rerender(dash)
	}
}

// Render renders the EnvelopesDashboardView component
func (dash *EnvelopesDashboardView) Render() vecty.ComponentOrHTML {
	if !dash.Rendered {
		return LoadingBar()
	}
	return Container(
		renderGatewayConnectionBanner(),
		renderServerConnectionBanner(),
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

// UpdateEnvelopesDashboard continuously updates the information needed by the Envelopes dashboard
func UpdateEnvelopesDashboard(d *EnvelopesDashboardView) {
	dispatcher.Dispatch(&actions.EnableAllUpdates{})

	ticker := time.NewTicker(time.Duration(store.Config.RefreshTime) * time.Second)
	updateEnvelopes(d)
	for {
		select {
		case <-store.RedirectChan:
			fmt.Println("Redirecting...")
			ticker.Stop()
			return
		case <-ticker.C:
			updateEnvelopes(d)
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
				getEnvelopes(d, util.Max(oldEnvelopes-store.Envelopes.Pagination.Index, 1))
				update.EnvelopeProcessResults()
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
			update.EnvelopeProcessResults()
		}
	}
}

func updateEnvelopes(d *EnvelopesDashboardView) {
	go dispatcher.Dispatch(&actions.GatewayConnected{Connected: store.GatewayClient.Ping()})
	go dispatcher.Dispatch(&actions.ServerConnected{Connected: api.PingServer()})
	actions.UpdateCounts()
	if !store.Envelopes.Pagination.DisableUpdate {
		getEnvelopes(d, util.Max(store.Envelopes.Count-store.Envelopes.Pagination.Index, 1))
		update.EnvelopeProcessResults()
	}
}

func getEnvelopes(d *EnvelopesDashboardView, index int) {
	fmt.Printf("Getting envelopes from index %d\n", index)
	list, ok := api.GetEnvelopeList(index)
	if ok {
		reverseEnvelopeList(&list)
		dispatcher.Dispatch(&actions.SetEnvelopeList{EnvelopeList: list})
	}
}

func reverseEnvelopeList(list *[config.ListSize]*proto.Envelope) {
	for i := len(list)/2 - 1; i >= 0; i-- {
		opp := len(list) - 1 - i
		list[i], list[opp] = list[opp], list[i]
	}
}

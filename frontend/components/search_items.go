package components

import (
	"fmt"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/logger"
	indexertypes "go.vocdoni.io/dvote/vochain/scrutinizer/indexertypes"
)

// SearchItemsView renders the processes dashboard page
type SearchItemsView struct {
	vecty.Core
	vecty.Mounter
	Rendered bool
}

// Mount is called after the component renders to signal that it can be rerendered safely
func (dash *SearchItemsView) Mount() {
	if !dash.Rendered {
		dash.Rendered = true
		vecty.Rerender(dash)
	}
}

// Render renders the SearchItemsView component
func (dash *SearchItemsView) Render() vecty.ComponentOrHTML {
	if !dash.Rendered {
		return LoadingBar()
	}
	if store.Loading {
		return Unavailable("Loading search...", "")
	}
	envelopes := len(store.Envelopes.Envelopes) > 0
	processes := len(store.Processes.ProcessIds) > 0
	entities := len(store.Entities.EntityIDs) > 0
	return Container(
		vecty.Markup(vecty.Attribute("id", "main")),
		renderServerConnectionBanner(),
		elem.Section(
			vecty.Markup(vecty.Class("details-view", "no-column")),
			elem.Div(
				vecty.Markup(vecty.Class("row")),
				elem.Div(
					vecty.Markup(vecty.Class("main-column")),
					vecty.If(envelopes || processes || entities, elem.Heading1(vecty.Text(fmt.Sprintf("Search results for \"%s\"", store.SearchTerm)))),
					vecty.If(envelopes, dash.EnvelopeList()),
					vecty.If(processes, dash.ProcessList()),
					vecty.If(entities, dash.EntityList()),
					vecty.If(!envelopes && !processes && !entities,
						bootstrap.Card(bootstrap.CardParams{
							Body: elem.Heading1(vecty.Text(fmt.Sprintf("No search results found for \"%s\"", store.SearchTerm))),
						})),
				),
			),
		),
	)
}

func (d *SearchItemsView) EnvelopeList() vecty.ComponentOrHTML {
	var elemList []vecty.MarkupOrChild
	if len(store.Envelopes.Envelopes) == 0 {
		elemList = append(elemList,
			vecty.Text("No envelopes found"))

	}
	for _, envelope := range store.Envelopes.Envelopes {
		elemList = append(elemList, renderProcessEnvelope(envelope))
	}
	return bootstrap.Card(bootstrap.CardParams{
		Body: vecty.List{
			elem.Heading1(vecty.Text("Vote envelopes")),
			elem.Div(elemList...),
		},
	})
}

func (d *SearchItemsView) ProcessList() vecty.ComponentOrHTML {
	var elemList []vecty.MarkupOrChild
	if len(store.Processes.ProcessIds) == 0 {
		elemList = append(elemList,
			vecty.Text("No processes found"))
	}
	for _, pid := range store.Processes.ProcessIds {
		process := store.Processes.Processes[pid]
		if process != nil {
			elemList = append(
				elemList,
				ProcessBlock(process),
			)
		}
	}
	return bootstrap.Card(bootstrap.CardParams{
		Body: vecty.List{
			elem.Heading1(vecty.Text("Processes")),
			elem.Div(elemList...),
		},
	})
}

func (d *SearchItemsView) EntityList() vecty.ComponentOrHTML {
	var elemList []vecty.MarkupOrChild
	if len(store.Entities.EntityIDs) == 0 {
		elemList = append(elemList,
			vecty.Text("No entities found"))
	}
	for _, ID := range store.Entities.EntityIDs {
		if ID != "" {
			height, hok := store.Entities.ProcessHeights[ID]
			if !hok {
				height = 0
			}
			elemList = append(
				elemList,
				EntityBlock(ID, height),
			)
		}
	}
	return bootstrap.Card(bootstrap.CardParams{
		Body: vecty.List{
			elem.Heading1(vecty.Text("Entities")),
			elem.Div(elemList...),
		},
	})
}

// UpdateProcessContents keeps the data for the processes dashboard up-to-date
func (dash *SearchItemsView) UpdateSearchItems(searchTerm string) {
	dispatcher.Dispatch(&actions.EnableAllUpdates{})
	dispatcher.Dispatch(&actions.SetLoading{Loading: true})
	if len(searchTerm) > 1 && (searchTerm[:2] == "0x" || searchTerm[:2] == "0X") {
		searchTerm = searchTerm[2:]
	}
	dispatcher.Dispatch(&actions.SetSearchTerm{SearchTerm: searchTerm})
	updateEnvelopeSearch(searchTerm)
	updateProcessSearch(searchTerm)
	updateEntitySearch(searchTerm)
	dispatcher.Dispatch(&actions.SetLoading{Loading: false})
}

func updateEnvelopeSearch(searchTerm string) {
	dispatcher.Dispatch(&actions.SetEnvelopeList{EnvelopeList: []*indexertypes.EnvelopeMetadata{}})
	list, err := store.Client.GetEnvelopeList([]byte{}, 0, config.ListSize, searchTerm)
	if err != nil {
		logger.Error(err)
		return
	}
	dispatcher.Dispatch(&actions.SetEnvelopeList{EnvelopeList: list})
}
func updateProcessSearch(searchTerm string) {
	dispatcher.Dispatch(&actions.SetProcessIds{Processes: []string{}})
	list, err := store.Client.GetProcessList([]byte{}, searchTerm, 0, "", false, 0, config.ListSize)
	if err != nil {
		logger.Error(err)
		return
	}
	fetchProcesses(list)
	dispatcher.Dispatch(&actions.SetProcessIds{Processes: list})
}
func updateEntitySearch(searchTerm string) {
	dispatcher.Dispatch(&actions.SetEntityIDs{EntityIDs: []string{}})
	list, err := store.Client.GetEntityList(searchTerm, config.ListSize, 0)
	if err != nil {
		logger.Error(err)
		return
	}
	dispatcher.Dispatch(&actions.SetEntityIDs{EntityIDs: list})
}

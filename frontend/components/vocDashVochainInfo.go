package components

import (
	"fmt"
	"strings"
	"syscall/js"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/client"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// VochainInfoView renders the vochainInfo pane
type VochainInfoView struct {
	vecty.Core
	entitiesIndex  int
	numEntities    int
	numProcesses   int
	processesIndex int
	refreshCh      chan bool
	vc             *client.VochainInfo
}

// Render renders the VochainInfoView component
func (b *VochainInfoView) Render() vecty.ComponentOrHTML {

	if b.vc != nil {
		if js.Global().Get("searchTerm").IsUndefined() || js.Global().Get("searchTerm").String() == "" {
			b.vc.EntitySearchIDs = util.TrimSlice(b.vc.EntityIDs, config.ListSize, &b.entitiesIndex)
			b.numEntities = len(b.vc.EntityIDs)
			b.vc.ProcessSearchIDs = util.TrimSlice(b.vc.ProcessIDs, config.ListSize, &b.processesIndex)
			b.numProcesses = len(b.vc.ProcessIDs)
		} else {
			search := js.Global().Get("searchTerm").String()
			temp := util.SearchSlice(b.vc.EntityIDs, search)
			b.vc.EntitySearchIDs = util.TrimSlice(temp, config.ListSize, &b.entitiesIndex)
			b.numEntities = len(temp)
			temp = util.SearchSlice(b.vc.ProcessIDs, search)
			if len(temp) <= 0 {
				temp = searchStateType(b.vc, search)
			}
			b.vc.ProcessSearchIDs = util.TrimSlice(temp, config.ListSize, &b.processesIndex)
			b.numProcesses = len(temp)
		}

		return elem.Section(
			bootstrap.Card(bootstrap.CardParams{
				Body: vecty.List{
					elem.Heading1(
						vecty.Text("Processes"),
					),
					b.ProcessesPaginatedList(),
				},
			}),
			bootstrap.Card(bootstrap.CardParams{
				Body: b.EntitiesPaginatedList(),
			}),
		)
	}
	return elem.Div(vecty.Text("Waiting for blockchain statistics..."))
}

func (b *VochainInfoView) EntitiesPaginatedList() vecty.ComponentOrHTML {
	slice := b.vc.EntityIDs

	return &Pagination{
		CurrentPage:     &store.Entities.CurrentPage,
		ListSize:        config.ListSize,
		TotalItems:      &b.numEntities,
		RefreshCh:       store.Entities.PagChannel,
		DisableUpdate:   &store.Entities.DisableUpdate,
		RenderSearchBar: false,
		RenderFunc: func(index int) vecty.ComponentOrHTML {
			if len(slice) == 0 {
				return elem.Div(
					vecty.Markup(vecty.Class("error")),
					vecty.Text("No valid entities"),
				)
			}

			list := make(vecty.List, len(slice))
			for _, ID := range slice {
				list = append(
					list,
					elem.ListItem(
						elem.Anchor(vecty.Markup(vecty.Attribute("href", "/entities/"+ID)), vecty.Text(ID)),
					),
				)
			}

			return elem.Div(list)
		},
	}
}

func (b *VochainInfoView) ProcessesPaginatedList() vecty.ComponentOrHTML {
	return &Pagination{
		CurrentPage:     &store.Processes.CurrentPage,
		ListSize:        config.ListSize,
		TotalItems:      &b.numProcesses,
		RefreshCh:       store.Processes.PagChannel,
		DisableUpdate:   &store.Processes.DisableUpdate,
		RenderSearchBar: false,
		RenderFunc: func(index int) vecty.ComponentOrHTML {
			list := make(vecty.List, len(b.vc.ProcessSearchIDs))
			IDs := b.vc.ProcessSearchIDs
			heights := b.vc.EnvelopeHeights
			procs := b.vc.ProcessSearchList

			for _, ID := range IDs {
				height, hok := heights[ID]
				info, iok := procs[ID]

				if !iok {
					list = append(
						list,
						elem.Div(
							vecty.Markup(vecty.Class("loading")),
							vecty.Text("Loading process info..."),
						),
					)
					continue
				}

				list = append(
					list,
					ProcessBlock(ID, hok, height, info),
				)
			}
			return elem.Div(list)
		},
	}
}

func ProcessBlock(ID string, hok bool, height int64, info client.ProcessInfo) vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(vecty.Class("tile", info.State)),
		elem.Div(
			vecty.Markup(vecty.Class("tile-body")),
			elem.Div(
				vecty.Markup(vecty.Class("type")),
				elem.Div(
					elem.Span(
						vecty.Markup(vecty.Class("title")),
						vecty.Text(info.ProcessType),
					),
					elem.Span(
						vecty.Markup(vecty.Class("status")),
						vecty.Text(info.State),
					),
				),
			),
			elem.Div(
				vecty.Markup(vecty.Class("contents")),
				elem.Div(
					elem.Div(
						elem.Anchor(
							vecty.Markup(vecty.Class("hash")),
							vecty.Markup(vecty.Attribute("href", "/processes/"+ID)),
							vecty.Text(ID),
						),
					),
					elem.Div(
						vecty.Markup(vecty.Class("envelopes")),
						vecty.Text(
							fmt.Sprintf("%d envelopes", height),
						),
					),
				),
			),
			elem.Div(
				vecty.Markup(vecty.Class("details")),
				elem.Div(
					vecty.Text("(date?)"),
				),
			),
		),
	)
}

func searchStateType(vc *client.VochainInfo, search string) []string {
	var IDList []string
	for _, key := range vc.ProcessIDs {
		// for key, info := range vc.ProcessSearchList {
		info, ok := vc.ProcessSearchList[key]
		if ok {
			if strings.Contains(info.State, search) || strings.Contains(info.ProcessType, search) {
				IDList = append(IDList, key)
			}
		}
	}
	return IDList
}

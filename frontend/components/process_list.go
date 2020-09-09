package components

import (
	"fmt"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/frontend/store/storeutil"
)

// ProcessListView renders the process list pane
type ProcessListView struct {
	vecty.Core
}

// Render renders the ProcessListView component
func (b *ProcessListView) Render() vecty.ComponentOrHTML {
	if store.Processes.Count > 0 {
		p := &Pagination{
			TotalPages:      int(store.Processes.Count) / config.ListSize,
			TotalItems:      &store.Processes.Count,
			CurrentPage:     &store.Processes.Pagination.CurrentPage,
			RefreshCh:       store.Processes.Pagination.PagChannel,
			ListSize:        config.ListSize,
			DisableUpdate:   &store.Processes.Pagination.DisableUpdate,
			SearchCh:        store.Processes.Pagination.SearchChannel,
			Searching:       &store.Processes.Pagination.Search,
			RenderSearchBar: true,
		}
		p.RenderFunc = func(index int) vecty.ComponentOrHTML {
			return elem.Div(renderProcessItems()...)
		}
		return p
	}
	return elem.Div(vecty.Text("No processes available"))
}

//ProcessBlock renders a single process card
func ProcessBlock(ID string, ok bool, height int64, info storeutil.Process) vecty.ComponentOrHTML {
	if !ok {
		return elem.Div(
			vecty.Markup(vecty.Class("tile", "empty")),
			elem.Div(
				vecty.Markup(vecty.Class("tile-body")),
				elem.Div(
					vecty.Markup(vecty.Class("type")),
					elem.Div(
						elem.Span(
							vecty.Markup(vecty.Class("title")),
							vecty.Text("Loading process..."),
						),
					),
				),
			),
		)
	}
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
						Link("/process/"+ID,
							ID,
							"hash",
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
			// elem.Div(
			// 	vecty.Markup(vecty.Class("details")),
			// 	elem.Div(
			// 		vecty.Text("(date?)"),
			// 	),
			// ),
		),
	)
}

func renderProcessItems() []vecty.MarkupOrChild {
	if len(store.Processes.ProcessIDs) == 0 {
		return []vecty.MarkupOrChild{vecty.Text("No valid processes")}
	}
	var elemList []vecty.MarkupOrChild
	for _, ID := range store.Processes.ProcessIDs {
		if ID != "" {
			height, _ := store.Processes.EnvelopeHeights[ID]
			info, iok := store.Processes.ProcessResults[ID]

			elemList = append(
				elemList,
				ProcessBlock(ID, iok, height, info),
			)
		}
	}
	return elemList
}

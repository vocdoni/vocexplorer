package components

import (
	"fmt"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/api/dbtypes"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/frontend/store/storeutil"
	"gitlab.com/vocdoni/vocexplorer/util"
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

func renderProcessItems() []vecty.MarkupOrChild {
	var elemList []vecty.MarkupOrChild
	for _, process := range store.Processes.Processes {
		if process != nil {
			height := store.Processes.EnvelopeHeights[process.ID]
			info, iok := store.Processes.ProcessResults[process.ID]

			elemList = append(
				elemList,
				ProcessBlock(process, iok, height, info),
			)
		}
	}
	return elemList
}

//ProcessBlock renders a single process card
func ProcessBlock(process *dbtypes.Process, ok bool, height int64, info storeutil.Process) vecty.ComponentOrHTML {
	if !ok || process == nil {
		return elem.Div(
			vecty.Markup(vecty.Class("tile")),
			elem.Div(
				vecty.Markup(vecty.Class("tile-body")),
				elem.Div(
					vecty.Markup(vecty.Class("type")),
					elem.Div(
						elem.Span(
							vecty.Markup(vecty.Class("title")),
							vecty.Text("Loading process..."),
						),
						elem.Span(
							vecty.Markup(vecty.Class("status")),
						),
					),
				),
				elem.Div(
					vecty.Markup(vecty.Class("contents")),
					elem.Div(),
					elem.Div(),
				),
			),
		)
	}
	entityHeight := store.Entities.ProcessHeights[process.EntityID]
	return elem.Div(
		vecty.Markup(vecty.Class("tile", strings.ToLower(info.ProcessInfo.State))),
		elem.Div(
			vecty.Markup(vecty.Class("tile-body")),
			elem.Div(
				vecty.Markup(vecty.Class("type")),
				elem.Div(
					elem.Span(
						vecty.Markup(vecty.Class("title")),
						vecty.Text(util.GetProcessName(info.ProcessInfo.Type)),
					),
					elem.Span(
						vecty.Markup(vecty.Class("status")),
						vecty.Text(strings.Title(info.ProcessInfo.State)),
					),
				),
			),
			elem.Div(
				vecty.Markup(vecty.Class("contents")),
				elem.Div(
					elem.Div(
						Link("/process/"+process.ID,
							process.ID,
							"hash",
						),
					),
					elem.Div(
						vecty.If(
							entityHeight < 1,
							vecty.Text(humanize.Ordinal(int(process.LocalHeight.Height+1))+" process hosted by entity "),
						),
						vecty.If(
							entityHeight > 1,
							vecty.Text(humanize.Ordinal(int(process.LocalHeight.Height+1))+" of "+util.IntToString(entityHeight)+" processes hosted by entity "),
						),
						vecty.If(
							entityHeight == 1,
							vecty.Text("only process hosted by entity "),
						),
						Link(
							"/entity/"+util.TrimHex(process.EntityID),
							util.TrimHex(process.EntityID),
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
		),
	)
}

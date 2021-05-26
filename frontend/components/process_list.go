package components

import (
	"fmt"
	"strings"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"

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
			SearchPrompt:    "search by process id, height",
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
	if len(store.Processes.Processes) == 0 {
		return []vecty.MarkupOrChild{elem.Div(vecty.Text("Loading processes..."))}
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
	return elemList
}

//ProcessBlock renders a single process card
func ProcessBlock(process *storeutil.Process) vecty.ComponentOrHTML {
	// Assume we only have access to process metadata: state, type, envelope height, enitity id
	if process == nil {
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
	if process.State == "" {
		process.State = "Unknown"
	}
	return elem.Div(
		vecty.Markup(vecty.Class("tile", strings.ToLower(process.State))),
		elem.Div(
			vecty.Markup(vecty.Class("tile-body")),
			elem.Div(
				vecty.Markup(vecty.Class("type")),
				elem.Div(
					elem.Span(
						vecty.Markup(vecty.Class("title")),
						vecty.Text(util.GetProcessName(process.Type)),
					),
					elem.Span(
						vecty.Markup(vecty.Class("status")),
						vecty.Text(strings.Title(process.State)),
					),
				),
			),
			elem.Div(
				vecty.Markup(vecty.Class("contents")),
				elem.Div(
					elem.Div(
						Link("/process/"+process.ProcessID,
							process.ProcessID,
							"hash",
						),
					),
					elem.Div(
						vecty.Text("Belongs to entity "),
						Link(
							"/entity/"+process.EntityID,
							process.EntityID,
							"hash",
						),
					),
					elem.Div(
						vecty.Markup(vecty.Class("envelopes")),
						vecty.If(process.EnvelopeCount == 1, vecty.Text(
							fmt.Sprintf("%d envelope", process.EnvelopeCount),
						)),
						vecty.If(process.EnvelopeCount != 1, vecty.Text(
							fmt.Sprintf("%d envelopes", process.EnvelopeCount),
						)),
					),
				),
			),
		),
	)
}

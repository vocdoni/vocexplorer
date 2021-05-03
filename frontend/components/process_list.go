package components

import (
	"fmt"
	"strings"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"go.vocdoni.io/proto/build/go/models"

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
	var tp string
	if process.Process.Envelope.Anonymous {
		tp = "anonymous"
	} else {
		tp = "poll"
	}
	if process.Process.Envelope.EncryptedVotes {
		tp += " encrypted"
	} else {
		tp += " open"
	}
	if process.Process.Envelope.Serial {
		tp += " serial"
	} else {
		tp += " single"
	}
	state := models.ProcessStatus(process.Process.Status).String()
	// entityHeight := store.Entities.ProcessHeights[util.HexToString(process.Process.EntityId)]
	return elem.Div(
		vecty.Markup(vecty.Class("tile", strings.ToLower(state))),
		elem.Div(
			vecty.Markup(vecty.Class("tile-body")),
			elem.Div(
				vecty.Markup(vecty.Class("type")),
				elem.Div(
					elem.Span(
						vecty.Markup(vecty.Class("title")),
						vecty.Text(util.GetProcessName(tp)),
					),
					elem.Span(
						vecty.Markup(vecty.Class("status")),
						vecty.Text(strings.Title(state)),
					),
				),
			),
			elem.Div(
				vecty.Markup(vecty.Class("contents")),
				elem.Div(
					elem.Div(
						Link("/process/"+util.HexToString(process.Process.ID),
							util.HexToString(process.Process.ID),
							"hash",
						),
					),
					elem.Div(
						// vecty.If(
						// 	entityHeight < 1,
						// 	vecty.Text(humanize.Ordinal(int(process.LocalHeight.Height+1))+" process hosted by entity "),
						// ),
						// vecty.If(
						// 	entityHeight > 1,
						// 	vecty.Text(humanize.Ordinal(int(process.LocalHeight.Height+1))+" of "+util.IntToString(entityHeight)+" processes hosted by entity "),
						// ),
						// vecty.If(
						// 	entityHeight == 1,
						// 	vecty.Text("only process hosted by entity "),
						// ),
						vecty.Text("Belongs to entity "),
						Link(
							"/entity/"+util.HexToString(process.Process.EntityID),
							util.HexToString(process.Process.EntityID),
							"hash",
						),
					),
					elem.Div(
						vecty.Markup(vecty.Class("envelopes")),
						vecty.Text(
							fmt.Sprintf("%d envelopes", process.EnvelopeCount),
						),
					),
				),
			),
		),
	)
}

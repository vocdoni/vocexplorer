package components

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
)

// EntityProcessListView renders the process list pane
type EntityProcessListView struct {
	vecty.Core
}

//Render renders the EntityProcessListView component
func (b *EntityProcessListView) Render() vecty.ComponentOrHTML {
	if store.Entities.CurrentEntity.ProcessCount > 0 {
		p := &Pagination{
			TotalPages:      int(store.Entities.CurrentEntity.ProcessCount) / config.ListSize,
			TotalItems:      &store.Entities.CurrentEntity.ProcessCount,
			CurrentPage:     &store.Entities.ProcessesPage,
			RefreshCh:       store.Entities.Pagination.PagChannel,
			ListSize:        config.ListSize,
			DisableUpdate:   &store.Entities.Pagination.DisableUpdate,
			SearchCh:        store.Entities.Pagination.SearchChannel,
			Searching:       &store.Entities.Pagination.Search,
			RenderSearchBar: true,
		}
		p.RenderFunc = func(index int) vecty.ComponentOrHTML {
			return elem.Div(renderEntityProcessItems()...)
		}

		return elem.Div(
			vecty.Markup(vecty.Class("recent-processes")),
			elem.Heading3(
				vecty.Text("Processes"),
			),
			p,
		)
	}
	if store.Entities.CurrentEntity.ProcessCount < 1 {
		return elem.Div(vecty.Text("No processes available"))
	}
	return elem.Div(vecty.Text("Waiting for processes..."))
}

func renderEntityProcessItems() []vecty.MarkupOrChild {
	if len(store.Entities.CurrentEntity.ProcessIDs) == 0 {
		return []vecty.MarkupOrChild{vecty.Text("No valid processes")}
	}
	var elemList []vecty.MarkupOrChild
	for _, ID := range store.Entities.CurrentEntity.ProcessIDs {
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

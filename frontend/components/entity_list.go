package components

import (
	"fmt"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
)

// EntityListView renders the entity list pane
type EntityListView struct {
	vecty.Core
}

// Render renders the EntityListView component
func (b *EntityListView) Render() vecty.ComponentOrHTML {
	if store.Entities.Count > 0 {
		p := &Pagination{
			TotalPages:      int(store.Entities.Count) / config.ListSize,
			TotalItems:      &store.Entities.Count,
			CurrentPage:     &store.Entities.Pagination.CurrentPage,
			RefreshCh:       store.Entities.Pagination.PagChannel,
			ListSize:        config.ListSize,
			DisableUpdate:   &store.Entities.Pagination.DisableUpdate,
			SearchCh:        store.Entities.Pagination.SearchChannel,
			Searching:       &store.Entities.Pagination.Search,
			RenderSearchBar: true,
		}
		p.RenderFunc = func(index int) vecty.ComponentOrHTML {
			return elem.Div(renderEntityItems()...)
		}
		return p
	}
	return elem.Div(vecty.Text("No entities available"))
}

//EntityBlock renders a single entity card
func EntityBlock(ID string, height int64) vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(vecty.Class("tile")),
		elem.Div(
			vecty.Markup(vecty.Class("tile-body")),
			elem.Div(
				vecty.Markup(vecty.Class("type")),
				elem.Div(
					elem.Span(
						vecty.Markup(vecty.Class("title")),
						elem.Anchor(
							vecty.Markup(
								vecty.Attribute("href", store.EntityDomain+ID),
								vecty.Property("target", ID),
							),
							vecty.Markup(vecty.Attribute("aria-label", "Link to entity "+store.Entities.CurrentEntityID+"'s profile page")),
							vecty.Text("Entity Profile"),
						),
					),
				),
			),
			elem.Div(
				vecty.Markup(vecty.Class("contents")),
				elem.Div(
					elem.Div(
						Link(
							"/entity/"+ID,
							ID,
							"hash",
						),
					),
					elem.Div(
						vecty.Markup(vecty.Class("envelopes")),
						vecty.If(height == 1,
							vecty.Text(
								fmt.Sprintf("%d process", height),
							),
						),
						vecty.If(height != 1,
							vecty.Text(
								fmt.Sprintf("%d processes", height),
							),
						),
					),
				),
			),
		),
	)
}

func renderEntityItems() []vecty.MarkupOrChild {
	if len(store.Entities.EntityIDs) == 0 {
		return []vecty.MarkupOrChild{vecty.Text("No valid entities")}
	}
	var elemList []vecty.MarkupOrChild
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
	return elemList
}

package components

import (
	"fmt"
	"strconv"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/gopherjs/vecty/prop"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/util"
	router "marwan.io/vecty-router"
)

// EntityListView renders the entity list pane
type EntityListView struct {
	vecty.Core
	currentPage int
}

// Render renders the EntityListView component
func (b *EntityListView) Render() vecty.ComponentOrHTML {
	if store.Entities.Count > 0 {
		p := &Pagination{
			TotalPages:      int(store.Entities.Count) / config.ListSize,
			TotalItems:      &store.Entities.Count,
			CurrentPage:     &b.currentPage,
			RefreshCh:       store.Entities.Pagination.PagChannel,
			ListSize:        config.ListSize,
			DisableUpdate:   &store.Entities.Pagination.DisableUpdate,
			RenderSearchBar: true,
		}
		p.RenderFunc = func(index int) vecty.ComponentOrHTML {
			return elem.Div(renderEntityItems()...)
		}
		p.SearchBar = func(self *Pagination) vecty.ComponentOrHTML {
			return elem.Input(vecty.Markup(
				event.Input(func(e *vecty.Event) {
					search := e.Target.Get("value").String()
					index, err := strconv.Atoi(e.Target.Get("value").String())
					if err != nil || index < 0 || index > int(*self.TotalItems) || search == "" {
						*self.CurrentPage = 0
						dispatcher.Dispatch(&actions.DisableEntityUpdate{Disabled: false})
						self.RefreshCh <- *self.CurrentPage * config.ListSize
					} else {
						*self.CurrentPage = util.Max(int(*self.TotalItems)-index-1, 0) / config.ListSize
						dispatcher.Dispatch(&actions.DisableEntityUpdate{Disabled: true})
						self.RefreshCh <- int(*self.TotalItems) - index
					}
					vecty.Rerender(self)
				}),
				prop.Placeholder("search entities"),
			))
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
							vecty.Markup(vecty.Attribute("href", "https://manage.vocdoni.net/entities/#/0x"+ID)),
							vecty.Text("Entity Manager Page"),
						),
					),
				),
			),
			elem.Div(
				vecty.Markup(vecty.Class("contents")),
				elem.Div(
					elem.Div(
						router.Link(
							"/entities/"+ID,
							ID,
							router.LinkOptions{
								Class: "hash",
							},
						),
					),
					elem.Div(
						vecty.Markup(vecty.Class("envelopes")),
						vecty.Text(
							fmt.Sprintf("%d processes", height),
						),
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

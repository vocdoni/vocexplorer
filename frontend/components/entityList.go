package components

import (
	"fmt"
	"strconv"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/gopherjs/vecty/prop"
	"gitlab.com/vocdoni/vocexplorer/client"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// EntityListView renders the entity list pane
type EntityListView struct {
	vecty.Core
	currentPage   int
	disableUpdate *bool
	refreshCh     chan int
	vochain       *client.VochainInfo
}

// Render renders the EntityListView component
func (b *EntityListView) Render() vecty.ComponentOrHTML {
	if b.vochain != nil && b.vochain.EntityCount > 0 {
		p := &Pagination{
			TotalPages:      int(b.vochain.EntityCount) / config.ListSize,
			TotalItems:      &b.vochain.EntityCount,
			CurrentPage:     &b.currentPage,
			RefreshCh:       b.refreshCh,
			ListSize:        config.ListSize,
			DisableUpdate:   b.disableUpdate,
			RenderSearchBar: true,
		}
		p.RenderFunc = func(index int) vecty.ComponentOrHTML {
			return elem.Div(renderEntityItems(b.vochain.EntityIDs, b.vochain.ProcessHeights)...)
		}
		p.SearchBar = func(self *Pagination) vecty.ComponentOrHTML {
			return elem.Input(vecty.Markup(
				event.Input(func(e *vecty.Event) {
					search := e.Target.Get("value").String()
					index, err := strconv.Atoi(e.Target.Get("value").String())
					if err != nil || index < 0 || index > int(*self.TotalItems) || search == "" {
						*self.CurrentPage = 0
						*b.disableUpdate = false
						self.RefreshCh <- *self.CurrentPage * config.ListSize
					} else {
						*self.CurrentPage = util.Max(int(*self.TotalItems)-index-1, 0) / config.ListSize
						*b.disableUpdate = true
						self.RefreshCh <- int(*self.TotalItems) - index
					}
					vecty.Rerender(self)
				}),
				prop.Placeholder("search entities"),
			))
		}
		return p
	}
	if b.vochain.EntityCount < 1 {
		return elem.Div(vecty.Text("No entities available"))
	}
	return elem.Div(vecty.Text("Waiting for entities..."))
}
func EntityBlock(ID string, height int64) vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(vecty.Class("tile")),
		elem.Div(
			vecty.Markup(vecty.Class("tile-body")),
			elem.Div(
				vecty.Markup(vecty.Class("type")),
			// elem.Div(
			// 	elem.Span(
			// 		vecty.Markup(vecty.Class("title")),
			// 		vecty.Text(ID),
			// 	),
			),
		),
		elem.Div(
			vecty.Markup(vecty.Class("contents")),
			elem.Div(
				elem.Div(
					elem.Anchor(
						vecty.Markup(vecty.Class("hash")),
						vecty.Markup(vecty.Attribute("href", "/entities/"+ID)),
						vecty.Text(ID),
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
	// elem.Div(
	// 	vecty.Markup(vecty.Class("details")),
	// 	elem.Div(
	// 		vecty.Text("(date?)"),
	// 	),
	// ),
	)
}

func renderEntityItems(slice [config.ListSize]string, heights map[string]int64) []vecty.MarkupOrChild {
	if len(slice) == 0 {
		return []vecty.MarkupOrChild{vecty.Text("No valid entities")}
	}
	var elemList []vecty.MarkupOrChild
	for _, ID := range slice {
		if ID != "" {
			height, hok := heights[ID]
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

package components

import (
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
			RenderSearchBar: true,
		}
		p.RenderFunc = func(index int) vecty.ComponentOrHTML {
			return elem.Div(renderEntityProcessItems()...)
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
				prop.Placeholder("search processes"),
			))
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

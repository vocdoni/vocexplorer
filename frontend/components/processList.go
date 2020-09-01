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
	"gitlab.com/vocdoni/vocexplorer/frontend/store/storeutil"
	"gitlab.com/vocdoni/vocexplorer/util"
	router "marwan.io/vecty-router"
)

// ProcessListView renders the process list pane
type ProcessListView struct {
	vecty.Core
	currentPage int
}

// Render renders the ProcessListView component
func (b *ProcessListView) Render() vecty.ComponentOrHTML {
	fmt.Println(store.Processes.Count)
	if store.Processes.Count > 0 {
		p := &Pagination{
			TotalPages:      int(store.Processes.Count) / config.ListSize,
			TotalItems:      &store.Processes.Count,
			CurrentPage:     &b.currentPage,
			RefreshCh:       store.Processes.Pagination.PagChannel,
			ListSize:        config.ListSize,
			DisableUpdate:   &store.Processes.Pagination.DisableUpdate,
			RenderSearchBar: true,
		}
		p.RenderFunc = func(index int) vecty.ComponentOrHTML {
			return elem.Div(renderProcessItems()...)
		}
		p.SearchBar = func(self *Pagination) vecty.ComponentOrHTML {
			return elem.Input(vecty.Markup(
				event.Input(func(e *vecty.Event) {
					search := e.Target.Get("value").String()
					index, err := strconv.Atoi(e.Target.Get("value").String())
					if err != nil || index < 0 || index > int(*self.TotalItems) || search == "" {
						*self.CurrentPage = 0
						dispatcher.Dispatch(&actions.DisableProcessUpdate{Disabled: false})
						self.RefreshCh <- *self.CurrentPage * config.ListSize
					} else {
						*self.CurrentPage = util.Max(int(*self.TotalItems)-index-1, 0) / config.ListSize
						dispatcher.Dispatch(&actions.DisableProcessUpdate{Disabled: true})
						self.RefreshCh <- int(*self.TotalItems) - index
					}
					vecty.Rerender(self)
				}),
				prop.Placeholder("search processes"),
			))
		}
		return p
	}
	return elem.Div(vecty.Text("No processes available"))
}

//ProcessBlock renders a single process card
func ProcessBlock(ID string, hok bool, height int64, info storeutil.Process) vecty.ComponentOrHTML {
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
						router.Link("/processes/"+ID, ID, router.LinkOptions{
							Class: "hash",
						}),
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

func renderProcessItems() []vecty.MarkupOrChild {
	if len(store.Processes.ProcessIDs) == 0 {
		return []vecty.MarkupOrChild{vecty.Text("No valid processes")}
	}
	var elemList []vecty.MarkupOrChild
	for _, ID := range store.Processes.ProcessIDs {
		if ID != "" {
			height, hok := store.Processes.EnvelopeHeights[ID]
			info, iok := store.Processes.ProcessResults[ID]

			if !iok {
				elemList = append(
					elemList,
					elem.Div(
						vecty.Markup(vecty.Class("loading")),
						vecty.Text("Loading process info..."),
					),
				)
				continue
			}

			elemList = append(
				elemList,
				ProcessBlock(ID, hok, height, info),
			)
		}
	}
	return elemList
}

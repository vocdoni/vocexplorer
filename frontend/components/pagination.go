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
	"gitlab.com/vocdoni/vocexplorer/util"
)

// Pagination holds pages of information (blocks, processes, etc)
type Pagination struct {
	vecty.Core
	CurrentPage     *int
	ListSize        int
	RefreshCh       chan int
	SearchCh        chan string
	RenderFunc      func(int) vecty.ComponentOrHTML
	RenderSearchBar bool
	Searching       *bool
	DisableUpdate   *bool
	PageLeft        func(e *vecty.Event)
	PageRight       func(e *vecty.Event)
	PageStart       func(e *vecty.Event)
	TotalItems      *int
	TotalPages      int
}

// Render renders the pagination component
func (p *Pagination) Render() vecty.ComponentOrHTML {
	p.TotalPages = (*p.TotalItems - 1) / p.ListSize
	searchBar := func(*Pagination) vecty.ComponentOrHTML {
		return nil
	}
	if p.RenderSearchBar {
		searchBar = func(self *Pagination) vecty.ComponentOrHTML {
			return elem.Input(vecty.Markup(
				// Trigger when 'enter' is pressed
				event.Change(func(e *vecty.Event) {
					search := e.Target.Get("value").String()
					index, err := strconv.Atoi(search)
					if index < 0 || index > int(*self.TotalItems) || search == "" { // Search term is empty or <0
						*self.CurrentPage = 0
						*self.Searching = false
						dispatcher.Dispatch(&actions.DisableUpdate{Updater: p.DisableUpdate, Disabled: false})
						self.RefreshCh <- *self.CurrentPage * config.ListSize
					} else if index > 0 && err == nil { // Search term is an integer > 0
						*self.CurrentPage = util.Max(int(*self.TotalItems)-index-1, 0) / config.ListSize
						*self.Searching = false
						dispatcher.Dispatch(&actions.DisableUpdate{Updater: p.DisableUpdate, Disabled: true})
						self.RefreshCh <- int(*self.TotalItems) - index
					} else { // Search term is not int, is ID
						*self.CurrentPage = 0
						dispatcher.Dispatch(&actions.DisableUpdate{Updater: p.DisableUpdate, Disabled: true})
						*self.Searching = true
						self.SearchCh <- search
					}
					vecty.Rerender(self)
				}),
				prop.Placeholder("search by height, id"),
			))
		}
	}
	if p.PageLeft == nil {
		p.PageLeft = func(e *vecty.Event) {
			*p.CurrentPage = util.Max(*p.CurrentPage-1, 0)
			if *p.CurrentPage != 0 || *p.Searching {
				*p.DisableUpdate = true
			} else {
				*p.DisableUpdate = false
			}
			p.RefreshCh <- *p.CurrentPage * p.ListSize
			vecty.Rerender(p)
		}
	}
	if p.PageRight == nil {
		p.PageRight = func(e *vecty.Event) {
			*p.CurrentPage = util.Min(*p.CurrentPage+1, p.TotalPages)
			if *p.CurrentPage != 0 || *p.Searching {
				*p.DisableUpdate = true
			} else {
				*p.DisableUpdate = false
			}
			p.RefreshCh <- *p.CurrentPage * p.ListSize
			vecty.Rerender(p)
		}
	}
	if p.PageStart == nil {
		p.PageStart = func(e *vecty.Event) {
			*p.CurrentPage = 0
			*p.DisableUpdate = false
			*p.Searching = false
			p.RefreshCh <- *p.CurrentPage * p.ListSize
			vecty.Rerender(p)
		}
	}
	return elem.Div(
		elem.Navigation(
			vecty.Markup(vecty.Class("pagination-wrapper")),
			elem.Div(
				vecty.Markup(vecty.Class("page-count")),
				vecty.If(
					!*p.Searching,
					elem.Span(
						vecty.Text(
							fmt.Sprintf("Page %d of %d", *p.CurrentPage+1, p.TotalPages+1),
						),
					),
				),
			),
			vecty.If(p.RenderSearchBar, elem.Div(
				vecty.Markup(vecty.Class("pagination-searchbar")),
				searchBar(p),
			)),
			elem.UnorderedList(
				vecty.Markup(vecty.Class("pagination")),
				elem.ListItem(
					vecty.Markup(
						vecty.MarkupIf(
							*p.CurrentPage != 0 || *p.Searching,
							vecty.Class("page-item"),
						),
						vecty.MarkupIf(
							*p.CurrentPage == 0 && !*p.Searching,
							vecty.Class("page-item", "disabled"),
						),
					),
					elem.Button(
						vecty.Markup(
							vecty.Class("page-link"),
							event.Click(p.PageStart),
							vecty.MarkupIf(
								*p.CurrentPage != 0 || *p.Searching,
								prop.Disabled(false),
							),
							vecty.MarkupIf(
								*p.CurrentPage == 0 && !*p.Searching,
								prop.Disabled(true),
							),
						),
						elem.Span(
							vecty.Text("«"),
						),
						elem.Span(
							vecty.Markup(vecty.Class("sr-only")),
							vecty.Text("Back to top"),
						),
					),
				),
				elem.ListItem(
					vecty.Markup(
						vecty.MarkupIf(
							*p.CurrentPage > 0,
							vecty.Class("page-item"),
						),
						vecty.MarkupIf(
							*p.CurrentPage <= 0 || *p.Searching,
							vecty.Class("page-item", "disabled"),
						),
					),
					elem.Button(
						vecty.Text("prev"),
						vecty.Markup(
							vecty.Class("page-link"),
							event.Click(p.PageLeft),
							vecty.MarkupIf(
								*p.CurrentPage > 0,
								prop.Disabled(false),
							),
							vecty.MarkupIf(
								*p.CurrentPage < 1 || *p.Searching,
								prop.Disabled(true),
							),
						),
					),
				),
				elem.ListItem(
					vecty.Markup(
						vecty.MarkupIf(
							*p.CurrentPage < p.TotalPages,

							vecty.Class("page-item"),
						),
						vecty.MarkupIf(
							*p.CurrentPage >= p.TotalPages || *p.Searching,
							vecty.Class("page-item", "disabled"),
						),
					),
					elem.Button(vecty.Text("next"),
						vecty.Markup(
							vecty.Class("page-link"),
							event.Click(p.PageRight),
							vecty.MarkupIf(
								*p.CurrentPage < p.TotalPages,
								prop.Disabled(false),
							),
							vecty.MarkupIf(
								*p.CurrentPage >= p.TotalPages || *p.Searching,
								prop.Disabled(true),
							),
						),
					),
				),
			),
		),
		p.RenderFunc(*p.CurrentPage),
	)
}

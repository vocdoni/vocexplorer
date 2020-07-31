package components

import (
	"fmt"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/gopherjs/vecty/prop"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// Pagination holds pages of information (blocks, processes, etc)
type Pagination struct {
	vecty.Core
	TotalPages      int
	TotalItems      *int
	CurrentPage     *int
	ListSize        int
	RefreshCh       chan int
	RenderSearchBar bool
	SearchBar       func(*Pagination) vecty.ComponentOrHTML
	RenderFunc      func(int) vecty.ComponentOrHTML
}

// Render renders the pagination component
func (p *Pagination) Render() vecty.ComponentOrHTML {
	p.TotalPages = (*p.TotalItems - 1) / p.ListSize
	return elem.Div(
		elem.Navigation(
			elem.Div(
				vecty.Markup(vecty.Class("pagination-wrapper")),
				elem.Span(
					vecty.Markup(vecty.Class("page-count")),
					vecty.Text(
						fmt.Sprintf("Page %d", *p.CurrentPage+1),
					),
				),
			),
			// vecty.If(p.RenderSearchBar, p.SearchBar(p)),
			p.SearchBar(p),
			elem.UnorderedList(
				vecty.Markup(vecty.Class("pagination")),
				elem.ListItem(
					vecty.Markup(
						vecty.MarkupIf(
							*p.CurrentPage != 0,
							vecty.Class("page-item"),
						),
						vecty.MarkupIf(
							*p.CurrentPage == 0,
							vecty.Class("page-item", "disabled"),
						),
					),
					elem.Button(
						vecty.Markup(
							vecty.Class("page-link"),
							event.Click(func(e *vecty.Event) {
								*p.CurrentPage = 0
								p.RefreshCh <- *p.CurrentPage * p.ListSize
								vecty.Rerender(p)
							}),
							vecty.MarkupIf(
								*p.CurrentPage != 0,
								prop.Disabled(false),
							),
							vecty.MarkupIf(
								*p.CurrentPage == 0,
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
							*p.CurrentPage <= 0,
							vecty.Class("page-item", "disabled"),
						),
					),
					elem.Button(
						vecty.Text("prev"),
						vecty.Markup(
							vecty.Class("page-link"),
							event.Click(func(e *vecty.Event) {
								*p.CurrentPage = util.Max(*p.CurrentPage-1, 0)
								p.RefreshCh <- *p.CurrentPage * p.ListSize
								vecty.Rerender(p)
							}),
							vecty.MarkupIf(
								*p.CurrentPage > 0,
								prop.Disabled(false),
							),
							vecty.MarkupIf(
								*p.CurrentPage < 1,
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
							*p.CurrentPage >= p.TotalPages,
							vecty.Class("page-item", "disabled"),
						),
					),
					elem.Button(vecty.Text("next"),
						vecty.Markup(
							vecty.Class("page-link"),
							event.Click(func(e *vecty.Event) {
								*p.CurrentPage = util.Min(*p.CurrentPage+1, p.TotalPages)
								p.RefreshCh <- *p.CurrentPage * p.ListSize
								vecty.Rerender(p)
							}),
							vecty.MarkupIf(
								*p.CurrentPage < p.TotalPages,
								prop.Disabled(false),
							),
							vecty.MarkupIf(
								*p.CurrentPage >= p.TotalPages,
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

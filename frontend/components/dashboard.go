package components

import (
	"strconv"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
)

// DashboardView renders the dashboard landing page
type DashboardView struct {
	vecty.Core
	val int
}

// Render renders the DashboardView component
func (dash *DashboardView) Render() vecty.ComponentOrHTML {
	return elem.Div(
		&Header{currentPage: "dashboard"},
		getDash(dash),
	)
}

func getDash(dash *DashboardView) *vecty.HTML {
	return elem.Span(elem.Button(
		vecty.Markup(
			vecty.Class("button"),
			event.Click(func(e *vecty.Event) {
				dash.val++
				vecty.Rerender(dash)
			}),
		),
		vecty.Text("Reset"),
	),
		vecty.Text(strconv.Itoa(dash.val)),
	)
}

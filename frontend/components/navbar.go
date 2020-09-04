package components

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/gopherjs/vecty/prop"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	router "marwan.io/vecty-router"
)

// NavBar renders the navigation bar
type NavBar struct {
	vecty.Core
}

// Render renders the NavBar component
func (n *NavBar) Render() vecty.ComponentOrHTML {
	return elem.Navigation(
		vecty.Markup(
			vecty.Class("navbar", "navbar-expand-lg", "navbar-light", "bg-light"),
		),
		elem.Div(
			vecty.Markup(vecty.Class("container-fluid")),
			Link("/", "Vochain Explorer", "navbar-brand"),
			elem.Button(
				vecty.Markup(
					vecty.Class("navbar-toggler"),
					vecty.Attribute("type", "button"),
					vecty.Attribute("data-toggle", "collapse"),
					vecty.Attribute("data-target", "#navbar-main"),
				),
				elem.Span(vecty.Markup(vecty.Class("navbar-toggler-icon"))),
			),
			elem.Div(
				vecty.Markup(
					vecty.Class("collapse", "navbar-collapse"),
					vecty.Attribute("id", "navbar-main"),
				),
				elem.UnorderedList(
					vecty.Markup(
						vecty.Class("navbar-nav", "mr-auto"),
					),
					elem.ListItem(
						vecty.Markup(
							vecty.Class("nav-item", "active"),
						),
						NavLink("/", "Home"),
					),
					elem.ListItem(
						vecty.Markup(
							vecty.Class("nav-item", "dropdown"),
						),
						NavLink("/participation", "Processes & Entities"),
					),
					elem.ListItem(
						vecty.Markup(
							vecty.Class("nav-item", "dropdown"),
						),
						NavLink("/blocktxs", "Blocks & Transactions"),
					),
					elem.ListItem(
						vecty.Markup(
							vecty.Class("nav-item", "dropdown"),
						),
						NavLink("/validators", "Validators"),
					),
					// elem.ListItem(
					// 	vecty.Markup(
					// 		vecty.Class("nav-item", "dropdown"),
					// 	),
					// 	NavLink("/blocks", "Blocks"),
					// ),
					elem.ListItem(
						vecty.Markup(
							vecty.Class("nav-item", "dropdown"),
						),
						NavLink("/stats", "Stats"),
					),
				),
				// &SearchBar{},
			),
		),
	)
}

func NavLink(route, text string) *vecty.HTML {
	return Link(route, text, "nav-link")
}

// Link renders a link which, when clicks, signals a redirect
func Link(route, text, class string) *vecty.HTML {
	attrs := []vecty.Applyer{
		prop.Href(route),
		event.Click(
			func(e *vecty.Event) {
				dispatcher.Dispatch(&actions.SignalRedirect{})
				router.Redirect(route)
			},
		).PreventDefault(),
	}

	if class != "" {
		attrs = append(attrs, vecty.Class(class))
	}

	return elem.Anchor(
		vecty.Markup(
			attrs...,
		),
		vecty.Text(text),
	)
}

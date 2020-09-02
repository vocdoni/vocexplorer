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
						Link("/", "Home", "nav-link"),
					),
					elem.ListItem(
						vecty.Markup(
							vecty.Class("nav-item", "dropdown"),
						),
						Link("/vocdash", "Processes & Entities", "nav-link"),
					),
					elem.ListItem(
						vecty.Markup(
							vecty.Class("nav-item", "dropdown"),
						),
						Link("/blocktxs", "Blocks & Transactions", "nav-link"),
					),
					elem.ListItem(
						vecty.Markup(
							vecty.Class("nav-item", "dropdown"),
						),
						Link("/validators", "Validators", "nav-link"),
					),
					// elem.ListItem(
					// 	vecty.Markup(
					// 		vecty.Class("nav-item", "dropdown"),
					// 	),
					// 	Link("/blocks", "Blocks", "nav-link"),
					// ),
					elem.ListItem(
						vecty.Markup(
							vecty.Class("nav-item", "dropdown"),
						),
						Link("/stats", "Stats", "nav-link"),
					),
				),
				&SearchBar{},
			),
		),
	)
}

// Link renders a link which, when clicks, signals a redirect
func Link(route, text, class string) *vecty.HTML {
	if class == "" {
		class = "nav-link"
	}
	return elem.Anchor(
		vecty.Markup(
			prop.Href(route),
			vecty.Class(class),
			event.Click(
				func(e *vecty.Event) {
					dispatcher.Dispatch(&actions.SignalRedirect{})
					router.Redirect(route)
				},
			).PreventDefault(),
		),
		vecty.Text(text),
	)
}

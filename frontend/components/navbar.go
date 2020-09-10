package components

import (
	"strings"
	"syscall/js"

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
	active := getActivePage()
	return elem.Navigation(
		vecty.Markup(
			vecty.Class("navbar", "navbar-expand-lg", "navbar-dark", "bg-dark"),
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
							vecty.Class("nav-item"),
							vecty.MarkupIf(
								active == "home",
								vecty.Class("nav-item", "active"),
							),
						),
						NavLink("/", "Home"),
					),
					elem.ListItem(
						vecty.Markup(
							vecty.Class("nav-item"),
							vecty.MarkupIf(
								active == "blocks",
								vecty.Class("nav-item", "active"),
							),
						),
						NavLink("/blocks", "Blocks"),
					),
					elem.ListItem(
						vecty.Markup(
							vecty.Class("nav-item"),
							vecty.MarkupIf(
								active == "transactions",
								vecty.Class("nav-item", "active"),
							),
						),
						NavLink("/transactions", "Transactions"),
					),
					elem.ListItem(
						vecty.Markup(
							vecty.Class("nav-item"),
							vecty.MarkupIf(
								active == "entities",
								vecty.Class("nav-item", "active"),
							),
						),
						NavLink("/entities", "Entities"),
					),
					elem.ListItem(
						vecty.Markup(
							vecty.Class("nav-item"),
							vecty.MarkupIf(
								active == "processes",
								vecty.Class("nav-item", "active"),
							),
						),
						NavLink("/processes", "Processes"),
					),
					elem.ListItem(
						vecty.Markup(
							vecty.Class("nav-item"),
							vecty.MarkupIf(
								active == "envelopes",
								vecty.Class("nav-item", "active"),
							),
						),
						NavLink("/envelopes", "Vote Envelopes"),
					),
					elem.ListItem(
						vecty.Markup(
							vecty.Class("nav-item"),
							vecty.MarkupIf(
								active == "validators",
								vecty.Class("nav-item", "active"),
							),
						),
						NavLink("/validators", "Validators"),
					),
					elem.ListItem(
						vecty.Markup(
							vecty.Class("nav-item"),
							vecty.MarkupIf(
								active == "stats",
								vecty.Class("nav-item", "active"),
							),
						),
						NavLink("/stats", "Stats"),
					),
				),
				// &SearchBar{},
			),
		),
	)
}

// NavLink generates a Link with nav-link styling
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

func getActivePage() string {
	path := js.Global().Get("location").Get("pathname").String()
	active := ""
	switch {
	case strings.Contains(path, "block"):
		active = "blocks"
	case strings.Contains(path, "transaction"):
		active = "transactions"
	case strings.Contains(path, "entit"):
		active = "entities"
	case strings.Contains(path, "process"):
		active = "processes"
	case strings.Contains(path, "envelope"):
		active = "envelopes"
	case strings.Contains(path, "validator"):
		active = "validators"
	case strings.Contains(path, "stats"):
		active = "stats"
	default:
		active = "home"
	}
	return active
}

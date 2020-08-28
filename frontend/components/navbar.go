package components

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
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
			router.Link("/", "Vochain Explorer", router.LinkOptions{
				Class: "navbar-brand",
			}),
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
						router.Link("/", "Home", router.LinkOptions{
							Class: "nav-link",
						}),
					),
					elem.ListItem(
						vecty.Markup(
							vecty.Class("nav-item", "dropdown"),
						),
						router.Link("/vocdash", "Processes & Entities", router.LinkOptions{
							Class: "nav-link",
						}),
					),
					elem.ListItem(
						vecty.Markup(
							vecty.Class("nav-item", "dropdown"),
						),
						router.Link("/blocktxs", "Blocks & Transactions", router.LinkOptions{
							Class: "nav-link",
						}),
					),
					elem.ListItem(
						router.Link(
							"/validators",
							"Validators",
							router.LinkOptions{
								Class: "nav-link",
							},
						),
					),
					elem.ListItem(
						router.Link("/blocks", "Blocks", router.LinkOptions{
							Class: "nav-link",
						}),
					),
					elem.ListItem(
						router.Link("/stats", "Stats", router.LinkOptions{
							Class: "nav-link",
						}),
					),
				),
				&SearchBar{},
			),
		),
	)
}

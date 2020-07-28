package components

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
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
		elem.Anchor(
			vecty.Markup(
				vecty.Attribute("href", "/"),
				vecty.Class("navbar-brand"),
			),
			vecty.Text("Vochain Explorer"),
		),
		elem.Button(
			vecty.Markup(
				vecty.Class("navbar-toggler"),
				vecty.Attribute("type", "button"),
				vecty.Attribute("data-toggle", "collapse"),
				vecty.Attribute("data-target", "navbar-main"),
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
					elem.Anchor(
						vecty.Markup(
							vecty.Class("nav-link"),
							vecty.Attribute("href", "/"),
						),
						vecty.Text("Home"),
					),
				),
				elem.ListItem(
					vecty.Markup(
						vecty.Class("nav-item", "dropdown"),
					),
					elem.Anchor(
						vecty.Markup(
							vecty.Class("nav-link"),
							vecty.Attribute("href", "/vocdash"),
						),
						vecty.Text("Processes & Entities"),
					),
				),
			),
		),
		// &nav.HorzBar{
		// 	Divs: []vecty.MarkupOrChild{
		// 		elem.Div(
		// 			vecty.Markup(
		// 				marginLeft.Auto,
		// 				float.Left,
		// 				style.Padding(14, 16),
		// 				fontFamily.Arial,
		// 			),
		// 			elem.Anchor(vecty.Markup(vecty.Attribute("href", "/")), vecty.Text("Dashboard")),
		// 			// router.Link("/", "Dashboard", router.LinkOptions{}),
		// 		),
		// 		elem.Div(
		// 			vecty.Markup(
		// 				marginLeft.Auto,
		// 				float.Left,
		// 				style.Padding(14, 16),
		// 				fontFamily.Arial,
		// 			),
		// 			elem.Anchor(vecty.Markup(vecty.Attribute("href", "/vocdash")), vecty.Text("Voting Processes & Entities")),
		// 			// router.Link("/processes", "Voting Processes & Entities", router.LinkOptions{}),
		// 		),
		// 	},
		// },
	)
}

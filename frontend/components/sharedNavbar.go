package components

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/nathanhack/vectyUI/nav"
	"github.com/nathanhack/vectyUI/style"
	"github.com/nathanhack/vectyUI/style/backgroundImage"
	"github.com/nathanhack/vectyUI/style/float"
	"github.com/nathanhack/vectyUI/style/fontFamily"
	"github.com/nathanhack/vectyUI/style/marginLeft"
	router "marwan.io/vecty-router"
)

// NavBar renders the navigation bar
type NavBar struct {
	vecty.Core
}

// Render renders the NavBar component
func (n *NavBar) Render() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(
			float.Left,
		),
		&nav.HorzBar{
			Background: style.Background(
				backgroundImage.LinearGradient(backgroundImage.ToRight, "Blue", "Red"),
			),
			Divs: []vecty.MarkupOrChild{
				elem.Div(
					vecty.Markup(
						marginLeft.Auto,
						float.Left,
						style.Padding(14, 16),
						fontFamily.Arial,
					),
					router.Link("/", "Dashboard", router.LinkOptions{}),
				),
				elem.Div(
					vecty.Markup(
						marginLeft.Auto,
						float.Left,
						style.Padding(14, 16),
						fontFamily.Arial,
					),
					router.Link("/processes", "Voting Processes & Entities", router.LinkOptions{}),
				),
			},
		},
	)
}

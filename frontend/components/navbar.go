package components

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/nathanhack/vectyUI/nav"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/nathanhack/vectyUI/button"
	"github.com/nathanhack/vectyUI/materialdesign"
	"github.com/nathanhack/vectyUI/nav"
	"github.com/nathanhack/vectyUI/style"
	"github.com/nathanhack/vectyUI/style/backgroundImage"
	"github.com/nathanhack/vectyUI/style/backgroundImage/colorStop"
	"github.com/nathanhack/vectyUI/style/display"
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
	return elem.Div{
		&nav.HorzBar{
			Background: style.Background(
				backgroundImage.LinearGradient(backgroundImage.ToRight, "Blue", "Red"),
			),
			Divs: []vecty.MarkupOrChild{
				elem.Div(
					vecty.Markup(
						marginLeft.Auto,
						float.Right,
						style.Padding(14, 16),
						fontFamily.Arial,
					)
					router.Link("/", "dashboard", router.LinkOptions{}),
				),
				elem.Div(
					vecty.Markup(
						marginLeft.Auto,
						float.Right,
						style.Padding(14, 16),
						fontFamily.Arial,
					)
					router.Link("/blocks", "blocks", router.LinkOptions{}),
				),
				elem.Div(
					vecty.Markup(
						marginLeft.Auto,
						float.Right,
						style.Padding(14, 16),
						fontFamily.Arial,
					)
					router.Link("/txs", "transactions", router.LinkOptions{}),
				),
				elem.Div(
					vecty.Markup(
						marginLeft.Auto,
						float.Right,
						style.Padding(14, 16),
						fontFamily.Arial,
					)
					router.Link("/procs", "processes", router.LinkOptions{}),
				),
			},
	}
}

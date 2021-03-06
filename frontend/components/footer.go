package components

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

// Footer renders the page footer
func Footer() vecty.ComponentOrHTML {
	return elem.Footer(
		Container(
			elem.Paragraph(
				elem.Anchor(
					vecty.Markup(vecty.Attribute("href", "https://vocdoni.io")),
					vecty.Text("Powered by"),
					elem.Image(
						vecty.Markup(
							vecty.Attribute("src", "/static/img/logo_labeled_white.png"),
							vecty.Attribute("alt", "Vocdoni"),
						),
					),
				),
			),
		),
	)
}

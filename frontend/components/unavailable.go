package components

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/vocdoni/vocexplorer/frontend/bootstrap"
)

// Unavailable renders an "item unavailable" message bar
func Unavailable(text, secondaryText string) vecty.ComponentOrHTML {
	return Container(
		vecty.Markup(vecty.Attribute("id", "main")),
		renderServerConnectionBanner(),
		elem.Section(
			bootstrap.Card(bootstrap.CardParams{
				Body: vecty.List{
					elem.Heading2(
						vecty.Text(text),
					),
					elem.Heading4(vecty.Text(secondaryText)),
				},
			}),
		),
	)
}

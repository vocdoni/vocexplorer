package bootstrap

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
)

func RenderCard(Header, Body, Footer vecty.ComponentOrHTML) vecty.ComponentOrHTML {
	items := vecty.List{}

	if Header != nil {
		items = append(items, elem.Div(
			vecty.Markup(vecty.Class("card-header")),
			Header,
		))
	}
	if Body != nil {
		items = append(items, elem.Div(
			vecty.Markup(vecty.Class("card-body")),
			Body,
		))
	}
	if Footer != nil {
		items = append(items, elem.Div(
			vecty.Markup(vecty.Class("card-footer")),
			Body,
		))
	}

	return elem.Div(
		items,
	)
}

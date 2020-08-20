package bootstrap

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
)

type CardParams struct {
	Header     vecty.ComponentOrHTML
	Body       vecty.ComponentOrHTML
	Footer     vecty.ComponentOrHTML
	ClassNames []string
}

func Card(p CardParams) vecty.ComponentOrHTML {
	items := vecty.List{}
	classes := p.ClassNames
	classes = append(classes, "card")

	if p.Header != nil {
		items = append(items, elem.Div(
			vecty.Markup(vecty.Class("card-header")),
			p.Header,
		))
	}
	if p.Body != nil {
		items = append(items, elem.Div(
			vecty.Markup(vecty.Class("card-body")),
			p.Body,
		))
	}
	if p.Footer != nil {
		items = append(items, elem.Div(
			vecty.Markup(vecty.Class("card-footer")),
			p.Body,
		))
	}

	return elem.Div(
		vecty.Markup(vecty.Class(classes...)),
		items,
	)
}

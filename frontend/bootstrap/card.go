package bootstrap

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
)

type Card struct {
	vecty.Core
	Header vecty.ComponentOrHTML
	Body   vecty.ComponentOrHTML
	Footer vecty.ComponentOrHTML
}

func (c *Card) Render() vecty.ComponentOrHTML {
	items := vecty.List{}

	if c.Header != nil {
		items = append(items, elem.Div(
			vecty.Markup(vecty.Class("card-header")),
			c.Header,
		))
	}
	if c.Body != nil {
		items = append(items, elem.Div(
			vecty.Markup(vecty.Class("card-body")),
			c.Body,
		))
	}
	if c.Footer != nil {
		items = append(items, elem.Div(
			vecty.Markup(vecty.Class("card-footer")),
			c.Body,
		))
	}

	return elem.Div(
		vecty.Markup(vecty.Class("card")),
		items,
	)
}
